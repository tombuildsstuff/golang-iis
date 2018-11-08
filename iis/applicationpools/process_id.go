package applicationpools

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// GetWorkerProcessID returns the ID of the Worker Process used for this App Pool
// this will only return the ProcessID for an App Pool with an associated Website
// in the running state
func (c *AppPoolsClient) GetWorkerProcessID(name string) (*[]int, error) {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
$pids = dir "IIS:\AppPools\%s\WorkerProcesses" | Select-Object -expand processId | ConvertTo-Json

if ($pids.Count -eq 0) {
    Write-Host "[]"
} else {
    if ($pids.Count -gt 1) {
      $v = $pids | ConvertTo-Json
	    Write-Host $v
    } else {
       $v = "[""{0}""]" -f $pids[0].ToString()
        Write-Host $v
    }
}
  `, name)

	stdout, _, err := c.Run(commands)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Worker Process ID's for App Pool %q: %+v", name, err)
	}

	var processIds []string
	if out := stdout; out != nil && *out != "" {
		v := *out
		err := json.Unmarshal([]byte(v), &processIds)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling Worker Process ID's for App Pool %q: %+v", name, err)
		}
	}

	var actualProcessIds []int
	for _, v := range processIds {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing %q as an int: %+v", v, err)
		}

		actualProcessIds = append(actualProcessIds, i)
	}

	return &actualProcessIds, nil
}
