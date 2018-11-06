package websites

import (
	"fmt"
	"testing"

	"github.com/tombuildsstuff/golang-iis/iis/applicationpools"
	"github.com/tombuildsstuff/golang-iis/iis/cmd"
	"github.com/tombuildsstuff/golang-iis/iis/helpers"
)

func TestWebsiteExists(t *testing.T) {
	rInt := helpers.RandomInt()
	applicationPoolName := fmt.Sprintf("acctestpool-%d", rInt)
	websiteName := fmt.Sprintf("acctestwebsites-%d", rInt)

	appPoolsClient := applicationpools.AppPoolsClient{
		Client: cmd.Client{},
	}
	websiteClient := WebsitesClient{
		Client: cmd.Client{},
	}

	err := appPoolsClient.Create(applicationPoolName)
	if err != nil {
		t.Fatalf("Error creating Application Pool %q: %+v", applicationPoolName, err)
		return
	}

	defer appPoolsClient.Delete(applicationPoolName)

	err = websiteClient.Create(websiteName, applicationPoolName, defaultWebsitePath)
	if err != nil {
		t.Fatalf("Error creating Website %q in App Pool %q: %+v", websiteName, applicationPoolName, err)
		return
	}

	defer websiteClient.Delete(websiteName)

	exists, err := websiteClient.Exists(websiteName)
	if err != nil {
		t.Fatalf("Error checking if Website exists: %+v", err)
		return
	}

	if !*exists {
		t.Fatalf("Expected the Website %q to exist, but it didn't..", websiteName)
		return
	}
}

func TestWebsiteDoesNotExist(t *testing.T) {
	name := fmt.Sprintf("doesntexist%d", helpers.RandomInt())
	client := WebsitesClient{
		Client: cmd.Client{},
	}

	exists, err := client.Exists(name)
	if err != nil {
		t.Fatalf("Error checking if Website exists: %+v", err)
		return
	}

	if *exists {
		t.Fatalf("Expected the Website %q to not exist, but it did..", name)
		return
	}
}
