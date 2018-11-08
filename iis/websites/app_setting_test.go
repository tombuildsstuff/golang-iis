package websites

import (
	"fmt"
	"testing"

	"github.com/tombuildsstuff/golang-iis/iis/applicationpools"
	"github.com/tombuildsstuff/golang-iis/iis/cmd"
	"github.com/tombuildsstuff/golang-iis/iis/helpers"
)

func TestAppSetting(t *testing.T) {
	rInt := helpers.RandomInt()
	appPoolName := fmt.Sprintf("acctestpool-%d", rInt)
	websiteName := fmt.Sprintf("acctestsite-%d", rInt)
	appSettingName := "example-setting"

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

	vals := []string{"first", "second", "third"}
	for _, v := range vals {
		err = websitesClient.SetAppSetting(websiteName, appSettingName, v)
		if err != nil {
			t.Fatalf("Error setting the App Setting %q for Website %q to %q: %+v", appSettingName, websiteName, v, err)
			return
		}

		setting, err := websitesClient.GetAppSetting(websiteName, appSettingName)
		if err != nil {
			t.Fatalf("Error retrieving App Setting %q for Website %q: %+v", appSettingName, websiteName, err)
			return
		}

		if *setting != v {
			t.Fatalf("Expected App Setting %q to have the value %q but got %q", appSettingName, v, *setting)
			return
		}
	}

	websitesClient.Delete(websiteName)
	appPoolsClient.Delete(appPoolName)
}
