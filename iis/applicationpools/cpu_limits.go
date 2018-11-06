package applicationpools

import (
	"fmt"
)

// ResetCPULimits resets the specified CPU limits for the Application Pool.
// name is the name of the Application Pool
func (c *AppPoolsClient) ResetCPULimits(name string) error {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
Set-ItemProperty -Path "IIS:\AppPools\%s" -Name cpu.action -Value "NoAction"
Set-ItemProperty -Path "IIS:\AppPools\%s" -Name cpu.limit -Value 0
  `, name, name)

	_, stderr, err := c.Run(commands)
	if err != nil {
		return fmt.Errorf("Error resetting CPU Limits for Application Pool: %+v", err)
	}

	if stderr != nil && *stderr != "" {
		return fmt.Errorf("Error resetting CPU Limits for Application Pool %q: %+v", name, err)
	}

	return nil
}

// SetCPULimits sets the specified network limits for the Application Pool.
// name is the name of the website
// maxCPU is the max CPU (1/1000th/interval).
func (c *AppPoolsClient) SetCPULimits(name string, maxCpuPerInterval int64) error {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
Set-ItemProperty -Path "IIS:\AppPools\%s" -Name cpu.action -Value "Throttle"
Set-ItemProperty -Path "IIS:\AppPools\%s" -Name cpu.limit -Value %d
  `, name, name, maxCpuPerInterval)

	_, stderr, err := c.Run(commands)
	if err != nil {
		return fmt.Errorf("Error setting CPU Limits for Application Pool: %+v", err)
	}

	if stderr != nil && *stderr != "" {
		return fmt.Errorf("Error setting CPU Limits for Application Pool %q: %+v", name, err)
	}

	return nil
}
