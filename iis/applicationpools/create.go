package applicationpools

import (
	"fmt"
	"strings"
)

// Create will create an Application Pool within IIS.
func (c *AppPoolsClient) Create(name string) error {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
New-WebAppPool -Name %q
  `, name)

	// TODO: support for setting which user this should be run as -as a separate script
	/*
	   $appPool.processModel.userName = $userAccountName
	   $appPool.processModel.password = $userAccountPassword
	*/

	_, stderr, err := c.Run(commands)
	if err != nil {
		return fmt.Errorf("Error creating App Pool %q: %+v", name, err)
	}

	if serr := stderr; serr != nil {
		v := strings.TrimSpace(*serr)
		if v != "" {
			return fmt.Errorf("Error creating App Pool %q: %s", name, v)
		}
	}

	return nil
}
