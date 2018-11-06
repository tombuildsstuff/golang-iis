package websites

import "fmt"

func (c *WebsitesClient) Start(name string) error {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
Start-Website -Name %q
  `, name)

	_, stderr, err := c.Run(commands)
	if err != nil {
		return fmt.Errorf("Error starting Website: %+v", err)
	}

	if stderr != nil && *stderr != "" {
		return fmt.Errorf("Error starting Website %q: %+v", name, err)
	}

	return nil
}
