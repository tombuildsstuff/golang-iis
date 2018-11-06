package websites

import "fmt"

func (c WebsitesClient) RemoveBinding(websiteName string, ipAddress string, hostHeader string, port int) error {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
Remove-WebBinding -Name %q -IPAddress %q -HostHeader %q -Port %d
  `, websiteName, ipAddress, hostHeader, port)

	_, stderr, err := c.Run(commands)
	if err != nil {
		return fmt.Errorf("Error creating Binding %q for Website %q: %+v", hostHeader, websiteName, err)
	}

	if stderr != nil && *stderr != "" {
		return fmt.Errorf("Error creating Binding %q for Website %q: %+v", hostHeader, websiteName, err)
	}

	return nil
}
