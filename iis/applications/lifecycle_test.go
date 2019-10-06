package applications

import (
	"fmt"
	"github.com/tombuildsstuff/golang-iis/iis/applicationpools"
	"github.com/tombuildsstuff/golang-iis/iis/cmd"
	"github.com/tombuildsstuff/golang-iis/iis/helpers"
	"github.com/tombuildsstuff/golang-iis/iis/websites"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestApplicationLifecycle(t *testing.T) {
	rInt := helpers.RandomInt()
	appPoolName := fmt.Sprintf("app-pool-%d", rInt)
	websiteName := fmt.Sprintf("website-%d", rInt)
	rootDirectory := fmt.Sprintf("C:\\inetpub")

	t.Logf("[DEBUG] Creating IIS Clients")
	cmdClient := cmd.Client{}
	appPoolsClient := applicationpools.AppPoolsClient{
		Client: cmdClient,
	}
	websitesClient := websites.WebsitesClient{
		Client: cmdClient,
	}
	applicationsClient := ApplicationsClient{
		Client: cmdClient,
	}

	t.Logf("[DEBUG] Creating App Pool %q", appPoolName)
	if err := appPoolsClient.Create(appPoolName); err != nil {
		t.Fatalf("[FATAL] Error creating App Pool %q: %s", appPoolName, err)
	}

	t.Logf("[DEBUG] Creating Website %q to host Applications", websiteName)
	mainWebsitePath := fmt.Sprintf("%s\\wwwroot", rootDirectory)
	if err := websitesClient.Create(websiteName, appPoolName, mainWebsitePath); err != nil {
		t.Fatalf("[FATAL] Error creating Website %q: %s", websiteName, err)
	}

	createdApplications := make([]string, 0, 3)
	for i := 1; i <= 3; i++ {
		randInt := helpers.RandomInt()
		appName := fmt.Sprintf("application-%d", randInt)
		appPath := fmt.Sprintf("%s\\%s", rootDirectory, appName)

		t.Logf("[DEBUG] Creating test Application Content %q", appPath)
		err := createTestDirectory(t, appPath, testAppContent(appName))
		if err != nil {
			t.Fatal(err)
		}

		createdApplications = append(createdApplications, appName)

		t.Logf("[DEBUG] Adding test Application %q to IIS", appName)
		if err := applicationsClient.Create(appName, websiteName, appPoolName, appPath); err != nil {
			t.Fatalf("[FATAL] Error creating Application %q (Website %q, App Pool %q): %s", appName, websiteName, appPoolName, err)
		}
	}

	t.Logf("[DEBUG] Checking that the Applications exist")
	for _, appName := range createdApplications {
		testUrl := fmt.Sprintf("http://localhost/%s", appName)
		appPath := fmt.Sprintf("%s\\%s", rootDirectory, appName)

		if err := assertApplicationExists(t, applicationsClient, appName, appPath, websiteName, testUrl, testAppContent(appName)); err != nil {
			t.Fatal(err)
		}
	}

	// Cleaning things up in the reverse order they were created
	t.Logf("[DEBUG] Cleaning up created resources...")
	for _, appName := range createdApplications {
		t.Logf("[DEBUG] Removing Application %q", appName)
		appPath := fmt.Sprintf("%s\\%s", rootDirectory, appName)

		t.Logf("[DEBUG] Deleting Application %q from IIS", appName)
		if err := applicationsClient.Delete(appName, websiteName); err != nil {
			t.Fatalf("[FATAL] Error deleting Application %q from IIS", appName)
		}

		t.Logf("[DEBUG] Deleting the content directory %q", appPath)
		if err := os.RemoveAll(appPath); err != nil {
			t.Fatalf("[FATAL] Error deleting the directory %q: %s", appPath, err)
		}
	}

	t.Logf("[DEBUG] Removing Website %q from IIS", websiteName)
	if err := websitesClient.Delete(websiteName); err != nil {
		t.Fatalf("[FATAL] Error removing Website %q from IIS: %s", websiteName, err)
	}

	t.Logf("[DEBUG] Removing App Pool %q from IIS", appPoolName)
	if err := appPoolsClient.Delete(appPoolName); err != nil {
		t.Fatalf("[FATAL] Error removing App Pool %q from IIS: %s", appPoolName, err)
	}
}

func assertApplicationExists(t *testing.T, client ApplicationsClient, appName string, appPath string, websiteName string, testUrl, expectedBody string) error {
	t.Logf("[DEBUG] Testing that Application %q exists", appName)
	exists, err := client.Exists(appName, websiteName)
	if err != nil {
		return fmt.Errorf("[FATAL] Error checking that Application %q exists: %s", appName, err)
	} else if !*exists {
		return fmt.Errorf("[FATAL] Expected Application %q to exist but it doesn't!", appName)
	}

	t.Logf("[DEBUG] Getting Application properties for %q", appName)
	props, err := client.Get(appName, websiteName)
	if err != nil {
		return fmt.Errorf("[FATAL] Error getting properties for Application %q", appName)
	}

	if !strings.EqualFold(props.PhysicalPath, appPath) {
		return fmt.Errorf("[FATAL] Expected Application %q to have path %q but it has path %q", appName, appPath, props.PhysicalPath)
	}

	t.Logf("[DEBUG] Getting HTML content for Application %q from %q", appName, testUrl)
	resp, err := http.Get(testUrl)
	if err != nil {
		return fmt.Errorf("[FATAL] Error making HTTP call to %q: %s", testUrl, err)
	}
	respContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("[FATAL] Error getting HTTP response body for Application %q: %s", appName, err)
	}

	respBody := string(respContent)
	if expectedBody != respBody {
		return fmt.Errorf("[FATAL] Expected HTML content %q but got content %q", expectedBody, respBody)
	}

	return nil
}

func createTestDirectory(t *testing.T, path string, content string) error {
	t.Logf("[DEBUG] Creating test application folder %q", path)
	err := os.MkdirAll(path, os.ModeDir)
	if err != nil {
		return fmt.Errorf("[FATAL] Error creating directory %q: %s", path, err)
	}

	t.Logf("[DEBUG] Writing test HTML file")
	testFile, err := os.OpenFile(fmt.Sprintf("%s\\index.html", path), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer testFile.Close()
	if err != nil {
		return fmt.Errorf("[FATAL] Error creating test HTML file %q: %s", fmt.Sprintf("%s\\index.html", path), err)
	}

	_, err = testFile.WriteString(content)
	if err != nil {
		return fmt.Errorf("[FATAL] Error writing test HTML file: %s", err)
	}

	t.Logf("[DEBUG] Wrote test HTML content: %q", content)

	return nil
}

func testAppContent(str string) string {
	return fmt.Sprintf("<html><head><title>IIS Testing</title></head><body><h1>Hello World from Application %s!</h1></body></html>", str)
}
