package applications

import "fmt"

// Checks to see if the given Application exists under the given Website
func (c *ApplicationsClient) Exists(name string, parentWebsite string) (*bool, error) {
	// returns anything if the application exists, nothing if it doesn't
	command := fmt.Sprintf(`
Import-Module WebAdministration
Get-WebApplication -Name %q -Site %q | Select-Object -ExpandProperty Path | ConvertTo-Json -Compress
`, name, parentWebsite)

	stdout, stderr, err := c.Client.Run(command)
	if err != nil {
		return nil, fmt.Errorf("Error determining if Application %q exists: %+v", name, err)
	}

	if stderr != nil && *stderr != "" {
		exists := false
		return &exists, fmt.Errorf("Error retrieving Application: %s", *stderr)
	}

	if out := stdout; out != nil && *out != "" {
		exists := true
		return &exists, nil
	}

	exists := false
	return &exists, nil
}
