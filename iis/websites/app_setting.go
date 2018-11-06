package websites

import (
	"encoding/json"
	"fmt"
	"strings"
)

type getAppSettingResponse struct {
	Value string `json:"Value"`
}

func (c *WebsitesClient) GetAppSetting(websiteName string, name string) (*string, error) {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
$path = "IIS:\Sites\%s"
$keyPath = "/appSettings/add[@key='%s']"
Get-WebConfigurationProperty -pspath $path -filter $keyPath -name "value" | ConvertTo-Json
`, websiteName, name)

	stdout, _, err := c.Run(commands)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving App Setting for Website: %+v", err)
	}

	var resp getAppSettingResponse
	if out := stdout; out != nil && *out != "" {
		v := *out
		err := json.Unmarshal([]byte(v), &resp)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling App Setting for Website %q: %+v", websiteName, err)
		}
	}

	v := strings.TrimSpace(resp.Value)
	return &v, nil
}

func (c *WebsitesClient) SetAppSetting(websiteName string, name string, value string) error {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
$path = "IIS:\Sites\%s"
$key=%q
$value=%q
$keyPath = "/appSettings/add[@key='$key']"
$prop = Get-WebConfigurationProperty -pspath $path -filter $keyPath -name "value"
if ($prop -eq $null) {
    Add-WebConfigurationProperty -pspath $path -filter "appSettings" -name "." -value @{key=$key;value=$value}
} else {
    Set-WebConfigurationProperty -pspath $path -filter $keyPath -name "value" -value $value
}
  `, websiteName, name, value)

	_, _, err := c.Run(commands)
	if err != nil {
		return fmt.Errorf("Error setting App Setting for Website: %+v", err)
	}

	return nil
}
