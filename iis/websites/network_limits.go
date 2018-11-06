package websites

import "fmt"

// ResetNetworkLimits resets the specified network limits for the Website.
// name is the name of the website
func (c *WebsitesClient) ResetNetworkLimits(name string) error {
	// from IIS6 onwards: https://msdn.microsoft.com/en-us/library/ms524902(v=vs.90).aspx
	return c.SetNetworkLimits(name, 4294967295)
}

// SetNetworkLimits sets the specified network limits for the Website.
// name is the name of the website
// maxBandwidth is the max bandwidth per second in bytes.
func (c *WebsitesClient) SetNetworkLimits(name string, maxBandwidth int64) error {
	// The value must be in the range between 1024 and 4294967295 bytes per second
	if maxBandwidth < 1024 || maxBandwidth > 4294967295 {
		return fmt.Errorf("maxBandwidth must be between 1024 and 4294967295 - got %d", maxBandwidth)
	}

	commands := fmt.Sprintf(`
Import-Module WebAdministration
Set-WebConfigurationProperty '/system.applicationHost/sites/site[@name="%s"]' -Name Limits.MaxBandwidth -Value %d -Force
  `, name, maxBandwidth)

	_, stderr, err := c.Run(commands)
	if err != nil {
		return fmt.Errorf("Error setting Network Limits for Website: %+v", err)
	}

	if stderr != nil && *stderr != "" {
		return fmt.Errorf("Error Network Limits for Website %q: %+v", name, err)
	}

	return nil
}
