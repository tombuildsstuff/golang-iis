package applicationpools

import (
	"fmt"
	"strconv"
	"strings"
)

// SetStartMode will set the Start Mode for an Application Pool within IIS.
func (c *AppPoolsClient) SetStartMode(name string, autoStart bool, mode StartMode) error {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
$AppPool = Get-Item "IIS:\AppPools\%s"
$AppPool.autoStart = "%s"
$AppPool.startMode = "%s"
$AppPool | Set-Item
  `, name, strconv.FormatBool(autoStart), string(mode))

	_, stderr, err := c.Run(commands)
	if err != nil {
		return fmt.Errorf("Error configuring the Start Mode for App Pool %q: %+v", name, err)
	}

	if serr := stderr; serr != nil {
		v := strings.TrimSpace(*serr)
		if v != "" {
			return fmt.Errorf("Error configuring the Start Mode for App Pool %q: %s", name, v)
		}
	}

	return nil
}
