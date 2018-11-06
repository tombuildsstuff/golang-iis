package applicationpools

import (
	"fmt"
	"testing"

	"github.com/tombuildsstuff/golang-iis/iis/cmd"
	"github.com/tombuildsstuff/golang-iis/iis/helpers"
)

func TestAppPoolExists(t *testing.T) {
	name := fmt.Sprintf("acctestapppool-%d", helpers.RandomInt())
	client := AppPoolsClient{
		Client: cmd.Client{},
	}

	err := client.Create(name)
	if err != nil {
		t.Fatalf("Error creating App Pool %q: %+v", name, err)
		return
	}

	defer client.Delete(name)

	exists, err := client.Exists(name)
	if err != nil {
		t.Fatalf("Error checking if App Pool exists: %+v", err)
		return
	}

	if !*exists {
		t.Fatalf("Expected the App Pool %q to exist, but it didn't..", name)
		return
	}
}

func TestAppPoolDoesNotExist(t *testing.T) {
	name := fmt.Sprintf("doesntexist%d", helpers.RandomInt())
	client := AppPoolsClient{
		Client: cmd.Client{},
	}

	exists, err := client.Exists(name)
	if err != nil {
		t.Fatalf("Error checking if App Pool exists: %+v", err)
		return
	}

	if *exists {
		t.Fatalf("Expected the App Pool %q to not exist, but it did..", name)
		return
	}
}
