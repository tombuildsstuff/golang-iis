package websites

import (
	"fmt"
	"testing"

	"github.com/tombuildsstuff/golang-iis/iis/applicationpools"
	"github.com/tombuildsstuff/golang-iis/iis/cmd"
	"github.com/tombuildsstuff/golang-iis/iis/helpers"
)

func TestConnectionString(t *testing.T) {
	rInt := helpers.RandomInt()
	appPoolName := fmt.Sprintf("acctestpool-%d", rInt)
	websiteName := fmt.Sprintf("acctestsite-%d", rInt)
	connectionStringName := "example-setting"

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

	defer appPoolsClient.Delete(appPoolName)

	err = websitesClient.Create(websiteName, appPoolName, defaultWebsitePath)
	if err != nil {
		t.Fatalf("Error creating Website %q in App Pool %q: %+v", websiteName, appPoolName, err)
		return
	}

	defer websitesClient.Delete(websiteName)

	exists, err := websitesClient.Exists(websiteName)
	if err != nil {
		t.Fatalf("Error checking if Website %q exists (App Pool %q): %+v", websiteName, appPoolName, err)
		return
	}

	if !*exists {
		t.Fatalf("Expected Website %q to exist, but it didn't!", websiteName)
		return
	}

	vals := []string{"first", "second", "third"}
	for _, v := range vals {
		err = websitesClient.SetConnectionString(websiteName, connectionStringName, v)
		if err != nil {
			t.Fatalf("Error setting the Connection String %q for Website %q to %q: %+v", connectionStringName, websiteName, v, err)
			return
		}

		connectionString, err := websitesClient.GetConnectionString(websiteName, connectionStringName)
		if err != nil {
			t.Fatalf("Error retrieving Connection String %q for Website %q: %+v", connectionStringName, websiteName, err)
			return
		}

		if *connectionString != v {
			t.Fatalf("Expected Connection String %q to have the value %q but got %q", connectionStringName, v, *connectionString)
			return
		}
	}
}
