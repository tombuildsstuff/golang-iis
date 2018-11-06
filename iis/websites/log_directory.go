package websites

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tombuildsstuff/golang-iis/iis/helpers"
)

type getLogDirectoryResponse struct {
	Directory string `json="directory"`
}

func (c *WebsitesClient) GetLogDirectory(name string) (*string, error) {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
Get-ItemProperty "IIS:\Sites\%s" -name logFile | ConvertTo-Json -Compress
`, name)

	stdout, _, err := c.Run(commands)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Log Directory for Website: %+v", err)
	}

	var resp getLogDirectoryResponse
	if out := stdout; out != nil && *out != "" {
		v := *out
		err := json.Unmarshal([]byte(v), &resp)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling Log Directory for Website %q: %+v", name, err)
		}
	}

	d := helpers.FixPowerShellPath(resp.Directory)
	v := strings.TrimSpace(d)
	return &v, nil
}

func (c *WebsitesClient) SetLogDirectory(name string, physicalPath string) error {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
Set-ItemProperty "IIS:\Sites\%s" -name logFile -value @{directory=%q}
  `, name, physicalPath)

	_, stderr, err := c.Run(commands)
	if err != nil {
		return fmt.Errorf("Error setting Log Directory: %+v", err)
	}

	if stderr != nil && *stderr != "" {
		return fmt.Errorf("Error setting Log Directory %q: %+v", name, err)
	}

	return nil
}
