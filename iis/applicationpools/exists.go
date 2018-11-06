package applicationpools

import (
	"encoding/json"
	"fmt"
)

// Exists returns whether a given Application Pool exists within IIS.
func (c *AppPoolsClient) Exists(name string) (*bool, error) {
	// returns `["first", "second"]`
	command := `
Import-Module WebAdministration
$appPools = Get-Item IIS:\AppPools
$appPoolNames = $appPools.Children.Keys
if ($appPoolNames.Count -eq 0) {
    Write-Host "[]"
} else {
    if ($appPoolNames.Count -gt 1) {
        $v = $appPoolNames | ConvertTo-Json
	    Write-Host $v
    } else {
       $v = "[""{0}""]" -f $appPoolNames[0].ToString()
        Write-Host $v
    }
}
`
	stdout, stderr, err := c.Client.Run(command)
	if err != nil {
		return nil, fmt.Errorf("Error determining if AppPool %q exists: %+v", name, err)
	}

	if stderr != nil && *stderr != "" {
		exists := false
		return &exists, fmt.Errorf("Error retrieving App Pool: %s", *stderr)
	}

	if out := stdout; out != nil {
		var names []string
		err := json.Unmarshal([]byte(*out), &names)
		if err != nil {
			return nil, fmt.Errorf("Error parsing App Pools: %+v", err)
		}

		for _, v := range names {
			if v == name {
				exists := true
				return &exists, nil
			}
		}
	}

	exists := false
	return &exists, nil
}
