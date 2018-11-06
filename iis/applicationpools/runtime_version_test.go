package applicationpools

import (
	"fmt"
	"log"
	"testing"

	"github.com/tombuildsstuff/golang-iis/iis/cmd"
	"github.com/tombuildsstuff/golang-iis/iis/helpers"
)

func TestRuntimeVersion(t *testing.T) {
	name := fmt.Sprintf("acctestpool-%d", helpers.RandomInt())
	client := AppPoolsClient{
		Client: cmd.Client{},
	}

	err := client.Create(name)
	if err != nil {
		t.Fatalf("Error creating App Pool %q: %+v", name, err)
		return
	}

	defer client.Delete(name)

	appPool, err := client.Get(name)
	if err != nil {
		t.Fatalf("Error retrieving App Pool %q: %+v", name, err)
		return
	}

	if appPool.FrameworkVersion != ManagedFrameworkVersionFour {
		t.Fatalf("Expected the Default Managed Framework Version to be Four but got %q", appPool.FrameworkVersion)
		return
	}

	values := []ManagedFrameworkVersion{
		ManagedFrameworkVersionNone,
		ManagedFrameworkVersionTwo,
		ManagedFrameworkVersionFour,
	}
	for _, v := range values {
		log.Printf("Setting the Managed Runtime Version to %q..", v)
		err = client.SetRuntimeVersion(name, v)
		if err != nil {
			t.Fatalf("Error setting the Managed Runtime Version: %+v", err)
			return
		}

		appPool, err = client.Get(name)
		if err != nil {
			t.Fatalf("Error retrieving App Pool %q: %+v", name, err)
			return
		}

		if appPool.FrameworkVersion != v {
			t.Fatalf("Expected the Managed Framework Version to be %q but got %q", v, appPool.FrameworkVersion)
			return
		}
	}
}
