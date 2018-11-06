package applicationpools

import (
	"encoding/json"
	"fmt"
)

type getAppPool struct {
	Name                  string              `json:"name"`
	ManagedPipelineMode   string              `json:"managedPipelineMode"`
	ManagedRuntimeVersion string              `json:"managedRuntimeVersion"`
	Cpu                   getAppPoolCpuLimits `json:"cpu"`
}

type getAppPoolCpuLimits struct {
	Action string `json:"action"`
	Limit  int64  `json:"limit"`
}

func (c *AppPoolsClient) Get(name string) (*ApplicationPool, error) {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
Get-ItemProperty -Path "IIS:\AppPools\%s" | Select-Object | ConvertTo-Json
  `, name)

	stdout, _, err := c.Run(commands)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving App Pool %q: %+v", name, err)
	}

	var appPool getAppPool
	if out := stdout; out != nil && *out != "" {
		v := *out
		err := json.Unmarshal([]byte(v), &appPool)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling App Pool %q: %+v", name, err)
		}
	}

	pool := ApplicationPool{
		Name:              appPool.Name,
		FrameworkVersion:  ManagedFrameworkVersion(appPool.ManagedRuntimeVersion),
		MaxCPUPerInterval: appPool.Cpu.Limit,
	}
	return &pool, nil
}
