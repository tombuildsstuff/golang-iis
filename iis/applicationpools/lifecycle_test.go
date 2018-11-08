package applicationpools

import (
	"fmt"
	"testing"

	"github.com/tombuildsstuff/golang-iis/iis/helpers"
)

func TestApplicationPoolLifecycle(t *testing.T) {
	name := fmt.Sprintf("app-pool-%d", helpers.RandomInt())

	client := AppPoolsClient{}
	err := client.Create(name)
	if err != nil {
		t.Fatalf("Error creating App Pool %q: %+v", name, err)
		return
	}

	exists, err := client.Exists(name)
	if err != nil {
		t.Fatalf("Error checking if App Pool %q exists: %+v", name, err)
		return
	}

	if !*exists {
		t.Fatalf("Application Pool %q does not exist!", name)
		return
	}

	pool, err := client.Get(name)
	if err != nil {
		t.Fatalf("Error retrieving Application Pool %q: %+v", name, err)
	}

	if !pool.AutoStart {
		t.Fatalf("Expected the App Pool %q to be enabled by default but it wasn't", name)
	}

	if pool.StartMode != StartModeOnDemand {
		t.Fatalf("Expected the App Pool %q to be OnDemand but it wasn't", name)
	}

	if pool.State != StateStarted {
		t.Fatalf("Expected the App Pool %q State to be Started but it wasn't", name)
	}

	err = client.SetStartMode(name, false, StartModeAlwaysRunning)
	if err != nil {
		t.Fatalf("Error setting StartMode for Application Pool %q: %+v", name, err)
	}

	pool, err = client.Get(name)
	if err != nil {
		t.Fatalf("Error retrieving Application Pool %q: %+v", name, err)
	}

	if pool.AutoStart {
		t.Fatalf("Expected the App Pool %q to be disabled but it wasn't", name)
	}

	if pool.StartMode != StartModeAlwaysRunning {
		t.Fatalf("Expected the App Pool %q to be AlwaysRunning but it wasn't", name)
	}

	err = client.Stop(name)
	if err != nil {
		t.Fatalf("Error stopping Application Pool %q: %+v", name, err)
	}

	pool, err = client.Get(name)
	if err != nil {
		t.Fatalf("Error retrieving Application Pool %q: %+v", name, err)
	}

	if pool.State != StateStopped {
		t.Fatalf("Expected the App Pool %q State to be Stopped but it wasn't", name)
	}

	err = client.Start(name)
	if err != nil {
		t.Fatalf("Error starting Application Pool %q: %+v", name, err)
	}

	pool, err = client.Get(name)
	if err != nil {
		t.Fatalf("Error retrieving Application Pool %q: %+v", name, err)
	}

	if pool.State != StateStarted {
		t.Fatalf("Expected the App Pool %q State to be Started but it wasn't", name)
	}

	err = client.Delete(name)
	if err != nil {
		t.Fatalf("Error deleting App Pool %q: %+v", name, err)
		return
	}
}
