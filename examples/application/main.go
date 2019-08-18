package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"../../iis"
)

func main() {
	log.Printf("Starting test...")

	testAppPath := "C:\\inetpub\\testing"
	testAppContent := "<html><head><title>IIS Testing</title></head><body><h1>Hello World!</h1></body></html>"
	appPoolName := "testing"
	applicationName := "testing"
	websiteName := "Default Web Site"
	testUrl := "http://localhost/testing"

	err := createTestDirectory(testAppPath, testAppContent)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Creating IIS Client")
	client, err := iis.NewClient()
	if err != nil {
		log.Fatalf("Error creating IIS client: %s", err)
	}

	log.Printf("Creating IIS App Pool %q", appPoolName)
	err = createAppPool(client, appPoolName)
	if err != nil {
		_ = os.RemoveAll(testAppPath)
		log.Fatal(err)
	}

	log.Printf("Adding IIS Application %s to the Default Web Site", "testing")
	err = createApplication(client, applicationName, appPoolName, testAppPath, websiteName)
	if err != nil {
		_ = client.AppPools.Delete(appPoolName)
		_ = os.RemoveAll(testAppPath)
		log.Fatal(err)
	}

	log.Printf("Trying to get the test HTML from IIS %s", testUrl)
	testApplication(testUrl, testAppContent)

	log.Printf("Cleaning up...")

	cleanup(client, applicationName, websiteName, appPoolName, testAppPath)

	log.Printf("Done!")
}

func cleanup(client *iis.Client, applicationName string, websiteName string, appPoolName string, testAppPath string) {
	err := client.Applications.Delete(applicationName, websiteName)
	if err != nil {
		log.Printf("Error removing Application: %s", err)
		err = nil
	}
	err = client.AppPools.Delete(appPoolName)
	if err != nil {
		log.Printf("Error removing App Pool: %s", err)
		err = nil
	}
	err = os.RemoveAll(testAppPath)
	if err != nil {
		log.Printf("Error removing test directory: %s", err)
		err = nil
	}
}

func testApplication(testUrl string, testAppContent string) {
	resp, err := http.Get(testUrl)
	if err != nil {
		log.Printf("Error trying to get the test HTML from IIS: %s", err)
		err = nil
	} else {
		respContent, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error getting HTTP response body: %s", err)
			err = nil
		} else {
			log.Printf("Got test HTML content: %q", respContent)

			diff := strings.Compare(testAppContent, string(respContent))
			log.Printf("Difference between what was written and read: %d", diff)
		}
	}
}

func createApplication(client *iis.Client, applicationName string, appPoolName string, testAppPath string, websiteName string) error {
	err := client.Applications.Create(applicationName, appPoolName, testAppPath, websiteName)
	if err != nil {
		_ = client.AppPools.Delete(appPoolName)
		_ = os.RemoveAll(testAppPath)
		return fmt.Errorf("Error creating IIS Application: %s", err)
	}

	exists, err := client.Applications.Exists(applicationName, websiteName)
	if err != nil {
		_ = client.AppPools.Delete(appPoolName)
		_ = os.RemoveAll(testAppPath)
		return fmt.Errorf("Error checking for IIS Application existence: %s", err)
	}
	if !*exists {
		return fmt.Errorf("Tried to create Application %q under Website %q but it doesn't exist! %s", applicationName, websiteName, err)
	}
	log.Printf("Application created successfully!")

	return nil
}

func createAppPool(client *iis.Client, appPoolName string) error {
	err := client.AppPools.Create(appPoolName)
	if err != nil {
		return fmt.Errorf("Error creating IIS App Pool: %s", err)
	}

	exists, err := client.AppPools.Exists(appPoolName)
	if err != nil {
		return fmt.Errorf("Error checking for IIS App Pool existence: %s", err)
	}
	if !*exists {
		return fmt.Errorf("Tried to create an app pool %q but it doesn't exist! %s", appPoolName, err)
	}
	log.Printf("App Pool created successfully!")

	return nil
}

func createTestDirectory(testAppPath string, testAppContent string) error {
	log.Printf("Creating test application folder %q", testAppPath)
	err := os.MkdirAll(testAppPath, os.ModeDir)
	if err != nil {
		return fmt.Errorf("Error creating directory %q: %s", testAppPath, err)
	}

	log.Printf("Writing test HTML file")
	testFile, err := os.OpenFile(fmt.Sprintf("%s\\index.html", testAppPath), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	_, err = testFile.WriteString(testAppContent)
	if err != nil {
		_ = testFile.Close()
		_ = os.RemoveAll(testAppPath)
		return fmt.Errorf("Error writing test HTML file: %s", err)
	}
	_ = testFile.Close()
	log.Printf("Wrote test HTML content: %q", testAppContent)

	return nil
}
