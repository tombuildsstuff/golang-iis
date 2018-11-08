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

	err = client.Stop(name)
	if err != nil {
		t.Fatalf("Error stopping Application Pool %q: %+v", name, err)
	}

	err = client.Start(name)
	if err != nil {
		t.Fatalf("Error starting Application Pool %q: %+v", name, err)
	}

	err = client.Delete(name)
	if err != nil {
		t.Fatalf("Error deleting App Pool %q: %+v", name, err)
		return
	}
}
