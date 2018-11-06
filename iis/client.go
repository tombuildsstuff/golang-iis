package iis

import (
	"fmt"
	"strings"

	"github.com/tombuildsstuff/golang-iis/iis/applicationpools"
	"github.com/tombuildsstuff/golang-iis/iis/cmd"
	"github.com/tombuildsstuff/golang-iis/iis/websites"
)

type Client struct {
	AppPools applicationpools.AppPoolsClient
	Websites websites.WebsitesClient

	client cmd.Client
}

func (c Client) validate() error {
	commands := "Import-Module WebAdministration"
	_, stderr, err := c.client.Run(commands)
	if err != nil {
		return fmt.Errorf("Error running verification command: %+v", err)
	}

	if serr := stderr; serr != nil {
		v := strings.TrimSpace(*serr)
		if v != "" {
			return fmt.Errorf("Error verifying IIS/PowerShell is installed/configured correctly: %s", v)
		}
	}

	return nil
}

// TODO: possible a local user to run as, or something?
func NewClient() (*Client, error) {
	c := cmd.Client{}

	client := Client{
		AppPools: applicationpools.AppPoolsClient{
			Client: c,
		},
		Websites: websites.WebsitesClient{
			Client: c,
		},
	}

	if err := client.validate(); err != nil {
		return nil, fmt.Errorf("Error initializing Client: %+v", err)
	}

	return &client, nil
}
