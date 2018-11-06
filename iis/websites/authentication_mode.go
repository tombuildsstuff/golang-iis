package websites

import (
	"encoding/json"
	"fmt"
	"strings"
)

type getAuthenticationModeResponse struct {
	Value string `json:"Value"`
}

func (c *WebsitesClient) GetAuthenticationMode(websiteName string) (*AuthenticationMode, error) {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
$path = "IIS:\Sites\%s"
Get-WebConfigurationProperty -pspath $path -filter "system.web/authentication" -name "mode" | ConvertTo-Json -Compress
`, websiteName)

	stdout, _, err := c.Run(commands)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Authentication Mode for Website: %+v", err)
	}

	var resp getAuthenticationModeResponse
	if out := stdout; out != nil && *out != "" {
		v := *out
		err := json.Unmarshal([]byte(v), &resp)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling Authentication Mode for Website %q: %+v", websiteName, err)
		}
	}

	if resp.Value == "" {
		authMode := None
		return &authMode, nil
	}

	authMode := AuthenticationMode(resp.Value)
	return &authMode, nil
}

func (c *WebsitesClient) SetAuthenticationMode(websiteName string, mode AuthenticationMode) error {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
$path = "IIS:\Sites\%s"
Set-WebConfigurationProperty -pspath $path -filter "system.web/authentication" -name "mode" -value %q
  `, websiteName, string(mode))

	_, stderr, err := c.Run(commands)
	if err != nil {
		return fmt.Errorf("Error setting Authentication Mode for Website: %+v", err)
	}

	if stde := stderr; stde != nil && *stde != "" {
		e := strings.TrimSpace(*stde)
		return fmt.Errorf("Error setting Authentication Mode for Website: %+v", e)
	}

	return nil
}
