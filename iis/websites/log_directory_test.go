package websites

import (
	"fmt"
	"testing"

	"github.com/tombuildsstuff/golang-iis/iis/applicationpools"
	"github.com/tombuildsstuff/golang-iis/iis/cmd"
	"github.com/tombuildsstuff/golang-iis/iis/helpers"
)

func TestLogDirectory(t *testing.T) {
	rInt := helpers.RandomInt()
	appPoolName := fmt.Sprintf("acctestpool-%d", rInt)
	websiteName := fmt.Sprintf("acctestsite-%d", rInt)
	updatedPath := "C:\\inetpub"

	appPoolsClient := applicationpools.AppPoolsClient{
		Client: cmd.Client{},
	}
	websitesClient := WebsitesClient{
		Client: cmd.Client{},
	}

	err := appPoolsClient.Create(appPoolName)
	if err != nil {
		t.Fatalf("Error creating App Pool %q: %+v", appPoolName, err)
		return
	}

	err = websitesClient.Create(websiteName, appPoolName, defaultWebsitePath)
	if err != nil {
		t.Fatalf("Error creating Website %q in App Pool %q: %+v", websiteName, appPoolName, err)
		return
	}

	exists, err := websitesClient.Exists(websiteName)
	if err != nil {
		t.Fatalf("Error checking if Website %q exists (App Pool %q): %+v", websiteName, appPoolName, err)
		return
	}

	if !*exists {
		t.Fatalf("Expected Website %q to exist, but it didn't!", websiteName)
		return
	}

	path, err := websitesClient.GetLogDirectory(websiteName)
	if err != nil {
		t.Fatalf("Error checking the Log Directory for Website %q (App Pool %q): %+v", websiteName, appPoolName, err)
		return
	}

	if *path != defaultLogPath {
		t.Fatalf("Expected the Log Directory to be %q but was %q", defaultLogPath, *path)
		return
	}

	err = websitesClient.SetLogDirectory(websiteName, updatedPath)
	if err != nil {
		t.Fatalf("Error updating the Log Directory for Website %q (App Pool %q): %+v", websiteName, appPoolName, err)
		return
	}

	path, err = websitesClient.GetLogDirectory(websiteName)
	if err != nil {
		t.Fatalf("Error checking the updated Log Directory for Website %q (App Pool %q): %+v", websiteName, appPoolName, err)
		return
	}

	if *path != updatedPath {
		t.Fatalf("Expected the updated Log Directory to be %q but was %q", updatedPath, *path)
		return
	}

	websitesClient.Delete(websiteName)
	appPoolsClient.Delete(appPoolName)
}
