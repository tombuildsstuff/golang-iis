package applicationpools

import (
	"fmt"
	"strings"
)

// Delete will create an Application Pool within IIS.
func (c *AppPoolsClient) Delete(name string) error {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
Remove-WebAppPool -Name %q
  `, name)

	_, stderr, err := c.Run(commands)
	if err != nil {
		return fmt.Errorf("Error deleting App Pool %q: %+v", name, err)
	}

	if serr := stderr; serr != nil {
		v := strings.TrimSpace(*serr)
		if v != "" {
			return fmt.Errorf("Error deleting App Pool %q: %s", name, v)
		}
	}

	return nil
}
