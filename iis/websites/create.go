package websites

import (
	"fmt"
)

func (c *WebsitesClient) Create(name string, applicationPool string, physicalPath string) error {
	// we normalize the path via PS as otherwise AppSettings/AuthMode fails
	// when combining `C:\\inetpub\\site\web.config` (also fails on C:/inetpub/site\web.config`)
	commands := fmt.Sprintf(`
Import-Module WebAdministration
$path = [IO.Path]::GetFullPath(%q)
New-Website -Name %q -ApplicationPool %q -PhysicalPath $path
  `, physicalPath, name, applicationPool)

	_, stderr, err := c.Run(commands)
	if err != nil {
		return fmt.Errorf("Error creating Website: %+v", err)
	}

	if stderr != nil && *stderr != "" {
		return fmt.Errorf("Error creating Website %q: %+v", name, err)
	}

	return nil
}
