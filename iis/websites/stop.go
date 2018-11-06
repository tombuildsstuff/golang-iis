package websites

import "fmt"

func (c *WebsitesClient) Stop(name string) error {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
Stop-Website -Name %q
  `, name)

	_, stderr, err := c.Run(commands)
	if err != nil {
		return fmt.Errorf("Error stopping Website: %+v", err)
	}

	if stderr != nil && *stderr != "" {
		return fmt.Errorf("Error stopping Website %q: %+v", name, err)
	}

	return nil
}
