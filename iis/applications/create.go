package applications

import "fmt"

func (c *ApplicationsClient) Create(name string, applicationPool string, physicalPath string, parentWebsite string) error {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
$path = [IO.Path]::GetFullPath(%q)
New-WebApplication -Name %q -Site %q -PhysicalPath $path -ApplicationPool %q
`, physicalPath, name, parentWebsite, applicationPool)

	_, stderr, err := c.Run(commands)
	if err != nil {
		return fmt.Errorf("Error creating Application: %+v", err)
	}

	if stderr != nil && *stderr != "" {
		return fmt.Errorf("Error creating Application %q: %+v", name, err)
	}

	return nil
}
