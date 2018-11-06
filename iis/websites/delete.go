package websites

import (
	"fmt"
)

func (c *WebsitesClient) Delete(name string) error {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
Remove-Website -Name %q
  `, name)

	_, stderr, err := c.Run(commands)
	if err != nil {
		return fmt.Errorf("Error deleting Website: %+v", err)
	}

	if stderr != nil && *stderr != "" {
		return fmt.Errorf("Error deleting Website %q: %+v", name, err)
	}

	return nil
}
