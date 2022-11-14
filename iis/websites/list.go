package websites

import (
	"encoding/json"
	"fmt"
)

func (c WebsitesClient) List() ([]string, error) {
	commands := `
		Import-Module WebAdministration
		Get-Website | select name | ConvertTo-Json -Compress
  	`

	stdout, _, err := c.Run(commands)
	if err != nil {
		return nil, fmt.Errorf("Error listing all Website: %+v", err)
	}

	var sitesJson = []map[string]string{}

	if out := stdout; out != nil && *out != "" {
		v := *out
		err := json.Unmarshal([]byte(v), &sitesJson)

		if err != nil {
			// Powershell omits the outer array if the return count of sites == 1
			// Hence we unmarshal a single map here

			var singleSite = make(map[string]string)
			err := json.Unmarshal([]byte(v), &singleSite)
			sitesJson = append(sitesJson, singleSite)

			if err != nil {
				return nil, fmt.Errorf("Error unmarshalling Websites: %+v", err)
			}
		}
	}

	sites := []string{}
	for _, i := range sitesJson {
		sites = append(sites, i["name"])
	}

	return sites, nil
}
