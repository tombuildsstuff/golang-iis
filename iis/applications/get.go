package applications

import (
	"encoding/json"
	"fmt"
)

// PowerShell response shape
type getWebApplicationResponse struct {
	Path            string `json:"path"`
	ApplicationPool string `json:"applicationPool"`
	PhysicalPath    string `json:"physicalPath"`
}

// Get configuration data for an IIS Application
func (c *ApplicationsClient) Get(name string, parentWebsite string) (*Application, error) {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
Get-WebApplication -Name %q -Site %q | ConvertTo-Json -Compress
`, name, parentWebsite)

	stdout, _, err := c.Run(commands)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Application: %+v", err)
	}

	var app getWebApplicationResponse
	if out := stdout; out != nil && *out != "" {
		v := *out
		err := json.Unmarshal([]byte(v), &app)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling Application %q: %+v", name, err)
		}
	}

	if app.Path == "" {
		return nil, fmt.Errorf("Application %q was not found", name)
	}

	application := Application{
		Name:            name,
		Path:            app.Path,
		ApplicationPool: app.ApplicationPool,
		PhysicalPath:    app.PhysicalPath,
	}

	return &application, nil
}
