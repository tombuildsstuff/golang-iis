package applications

import "fmt"

func (c *ApplicationsClient) Delete(name string, parentWebsite string) error {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
Remove-WebApplication -Name %q -Site %q
  `, name, parentWebsite)

	_, stderr, err := c.Run(commands)
	if err != nil {
		return fmt.Errorf("Error deleting Application: %+v", err)
	}

	if stderr != nil && *stderr != "" {
		return fmt.Errorf("Error deleting Application %q: %+v", name, err)
	}

	return nil
}
