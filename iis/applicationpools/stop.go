package applicationpools

import (
	"fmt"
	"strings"
)

// Stop will start an Application Pool within IIS.
func (c *AppPoolsClient) Stop(name string) error {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
Stop-WebAppPool -Name %q
  `, name)

	_, stderr, err := c.Run(commands)
	if err != nil {
		return fmt.Errorf("Error stopping App Pool %q: %+v", name, err)
	}

	if serr := stderr; serr != nil {
		v := strings.TrimSpace(*serr)
		if v != "" {
			return fmt.Errorf("Error stopping App Pool %q: %s", name, v)
		}
	}

	return nil
}
