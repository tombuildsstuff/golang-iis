package applicationpools

import (
	"fmt"
	"testing"

	"github.com/tombuildsstuff/golang-iis/iis/cmd"
	"github.com/tombuildsstuff/golang-iis/iis/helpers"
)

func TestCPULimits(t *testing.T) {
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

	if appPool.MaxCPUPerInterval != 0 {
		t.Fatalf("Expected the initial Max CPU per interval to be 0 but got %d", appPool.MaxCPUPerInterval)
		return
	}

	err = client.SetCPULimits(name, 500)
	if err != nil {
		t.Fatalf("Error setting CPU Limits for App Pool: %+v", err)
	}

	appPool, err = client.Get(name)
	if err != nil {
		t.Fatalf("Error retrieving App Pool %q: %+v", name, err)
		return
	}

	if appPool.MaxCPUPerInterval != 500 {
		t.Fatalf("Expected the updated Max CPU per interval to be 500 but got %d", appPool.MaxCPUPerInterval)
		return
	}

	err = client.ResetCPULimits(name)
	if err != nil {
		t.Fatalf("Error resetting CPU Limits for App Pool: %+v", err)
	}

	appPool, err = client.Get(name)
	if err != nil {
		t.Fatalf("Error retrieving App Pool %q: %+v", name, err)
		return
	}

	if appPool.MaxCPUPerInterval != 0 {
		t.Fatalf("Expected the reset Max CPU per interval to be 0 but got %d", appPool.MaxCPUPerInterval)
		return
	}
}
