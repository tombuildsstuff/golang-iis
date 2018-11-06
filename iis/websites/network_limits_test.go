package websites

import (
	"fmt"
	"testing"

	"github.com/tombuildsstuff/golang-iis/iis/applicationpools"
	"github.com/tombuildsstuff/golang-iis/iis/cmd"
	"github.com/tombuildsstuff/golang-iis/iis/helpers"
)

func TestNetworkLimits(t *testing.T) {
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

	site, err := websitesClient.Get(websiteName)
	if err != nil {
		t.Fatalf("Error retrieving Website %q (App Pool %q): %+v", websiteName, appPoolName, err)
		return
	}

	if site.MaxBandwidthPerSecondInBytes != defaultNetworkLimit {
		t.Fatalf("Expected the default Max Bandwidth to be %d but got %d", defaultNetworkLimit, site.MaxBandwidthPerSecondInBytes)
		return
	}

	err = websitesClient.SetNetworkLimits(websiteName, 2048)
	if err != nil {
		t.Fatalf("Error setting Network Limits for %q: %+v", websiteName, err)
		return
	}

	site, err = websitesClient.Get(websiteName)
	if err != nil {
		t.Fatalf("Error re-retrieving Website %q (App Pool %q) after set: %+v", websiteName, appPoolName, err)
		return
	}

	if site.MaxBandwidthPerSecondInBytes != int64(2048) {
		t.Fatalf("Expected the updated Max Bandwidth to be 2048 but got %d", site.MaxBandwidthPerSecondInBytes)
		return
	}

	err = websitesClient.ResetNetworkLimits(websiteName)
	if err != nil {
		t.Fatalf("Error resetting Network Limits for Website %q (App Pool %q): %+v", websiteName, appPoolName, err)
		return
	}

	site, err = websitesClient.Get(websiteName)
	if err != nil {
		t.Fatalf("Error re-retrieving Website %q (App Pool %q) after reset: %+v", websiteName, appPoolName, err)
		return
	}

	if site.MaxBandwidthPerSecondInBytes != defaultNetworkLimit {
		t.Fatalf("Expected the default Max Bandwidth to be %d but got %d", defaultNetworkLimit, site.MaxBandwidthPerSecondInBytes)
		return
	}
}
