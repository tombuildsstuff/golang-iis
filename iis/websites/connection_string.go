package websites

import (
	"encoding/json"
	"fmt"
	"strings"
)

type getConnectionStringResponse struct {
	Value string `json:"Value"`
}

func (c *WebsitesClient) GetConnectionString(websiteName string, name string) (*string, error) {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
$path = "IIS:\Sites\%s"
$keyPath = "/connectionStrings/add[@name='%s']"
Get-WebConfigurationProperty -pspath $path -filter $keyPath -name "connectionString" | ConvertTo-Json
`, websiteName, name)

	stdout, _, err := c.Run(commands)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Connection String for Website: %+v", err)
	}

	var resp getConnectionStringResponse
	if out := stdout; out != nil && *out != "" {
		v := *out
		err := json.Unmarshal([]byte(v), &resp)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling Connection String for Website %q: %+v", websiteName, err)
		}
	}

	v := strings.TrimSpace(resp.Value)
	return &v, nil
}

func (c *WebsitesClient) SetConnectionString(websiteName string, name string, value string) error {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
$path = "IIS:\Sites\%s"
$key= %q
$connectionString = %q
$keyPath = "/connectionStrings/add[@name='$key']"
$prop = Get-WebConfigurationProperty -pspath $path -filter $keyPath -name "."
if ($prop -eq $null) {
    Add-WebConfigurationProperty -pspath $path -filter "connectionStrings" -name "." -value @{connectionString=$connectionString;name=$key}
} else {
    Set-WebConfigurationProperty -pspath $path -filter $keyPath -name "connectionString" -value $connectionString
}
  `, websiteName, name, value)

	_, _, err := c.Run(commands)
	if err != nil {
		return fmt.Errorf("Error setting Connection String for Website: %+v", err)
	}

	return nil
}
