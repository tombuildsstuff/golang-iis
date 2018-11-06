package websites

import (
	"fmt"
	"log"
	"testing"

	"github.com/tombuildsstuff/golang-iis/iis/applicationpools"
	"github.com/tombuildsstuff/golang-iis/iis/cmd"
	"github.com/tombuildsstuff/golang-iis/iis/helpers"
)

func TestAuthenticationMethod(t *testing.T) {
	rInt := helpers.RandomInt()
	appPoolName := fmt.Sprintf("acctestpool-%d", rInt)
	websiteName := fmt.Sprintf("acctestsite-%d", rInt)

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

	// TODO: switch over to Exists when it exists
	_, err = websitesClient.Get(websiteName)
	if err != nil {
		t.Fatalf("Error retrieving Website %q (App Pool %q): %+v", websiteName, appPoolName, err)
		return
	}

	modes := []AuthenticationMode{
		None,
		Federated,
		Forms,
		Passport,
		Windows,
	}

	for _, v := range modes {
		log.Printf("Setting Authentication Mode to %q..", v)
		err = websitesClient.SetAuthenticationMode(websiteName, v)
		if err != nil {
			t.Fatalf("Error setting the Authentication Mode for %q to %q: %+v", websiteName, string(v), err)
			return
		}

		mode, err := websitesClient.GetAuthenticationMode(websiteName)
		if err != nil {
			t.Fatalf("Error retrieving the Authentication Mode for %q: %+v", websiteName, err)
			return
		}

		if *mode != v {
			t.Fatalf("Expected the Authentication Mode to be %q but got %q", string(v), string(*mode))
			return
		}
	}
}
