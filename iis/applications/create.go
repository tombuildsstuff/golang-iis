package applications

import "fmt"

// Create an IIS Application under the given Website
func (c *ApplicationsClient) Create(name string, parentWebsite string, applicationPool string, physicalPath string) error {
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
