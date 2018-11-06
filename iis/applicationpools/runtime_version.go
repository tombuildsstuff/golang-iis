package applicationpools

import (
	"fmt"
)

// Create will create an Application Pool within IIS.
func (c *AppPoolsClient) SetRuntimeVersion(name string, frameworkVersion ManagedFrameworkVersion) error {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
Set-ItemProperty "IIS:\AppPools\%s" managedRuntimeVersion %q
  `, name, frameworkVersion)

	_, _, err := c.Run(commands)
	if err != nil {
		return fmt.Errorf("Error setting Managed Runtime Version for App Pool: %+v", err)
	}

	// TODO: error handling

	return nil
}
