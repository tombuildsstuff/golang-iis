package main

import (
	"fmt"
	"github.com/tombuildsstuff/golang-iis/iis"
	"github.com/tombuildsstuff/golang-iis/iis/helpers"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Printf("Starting example...")

	if err := run(); err != nil {
		log.Fatalf(err.Error())
	}
}

func run() error {
	randInt := helpers.RandomInt()
	testAppPath := fmt.Sprintf("C:\\inetpub\\ExampleWebApp-%d", randInt)
	testAppContent := "<html><head><title>IIS Testing</title></head><body><h1>Hello World!</h1></body></html>"
	applicationName := fmt.Sprintf("ExampleWebApp-%d", randInt)
	applicationPoolName := fmt.Sprintf("ExampleWebApp-%d", randInt)
	websiteName := fmt.Sprintf("ExampleWebsite-%d", randInt)
	websitePort := 8888
	websitePoolName := fmt.Sprintf("ExampleWebsite-%d", randInt)
	testUrl := fmt.Sprintf("http://localhost:%d/ExampleWebApp-%d", websitePort, randInt)

	err := createTestDirectory(testAppPath, testAppContent)
	if err != nil {
		return fmt.Errorf("Error creating the example directory: %s", err)
	}

	log.Printf("Creating IIS Client")
	client, err := iis.NewClient()
	if err != nil {
		return fmt.Errorf("Error creating IIS client: %s", err)
	}

	log.Printf("Creating IIS App Pool %q", websitePoolName)
	err = createAppPool(client, websitePoolName)
	if err != nil {
		return fmt.Errorf("Error creating App Pool: %s", err)
	}

	log.Printf("Creating IIS Website %q for the example", websiteName)
	err = createExampleWebsite(client, websiteName, websitePoolName, testAppPath, websitePort)
	if err != nil {
		return err
	}

	log.Printf("Creating IIS App Pool %q", applicationPoolName)
	err = createAppPool(client, applicationPoolName)
	if err != nil {
		return fmt.Errorf("Error creating App Pool: %s", err)
	}

	log.Printf("Adding IIS Application %q to the Website %q", applicationName, websiteName)
	err = createApplication(client, applicationName, applicationPoolName, testAppPath, websiteName)
	if err != nil {
		_ = client.AppPools.Delete(applicationPoolName)
		_ = os.RemoveAll(testAppPath)
		return fmt.Errorf("Error creating IIS Application: %s", err)
	}

	log.Printf("Trying to get the test HTML from IIS %s", testUrl)
	err = testApplication(testUrl, testAppContent)
	if err != nil {
		return err
	}

	log.Printf("Cleaning up...")

	cleanup(client, applicationName, applicationPoolName, websiteName, websitePoolName, testAppPath)

	log.Printf("Done!")

	return nil
}

func createExampleWebsite(client *iis.Client, websiteName string, websitePoolName string, testAppPath string, websitePort int) error {
	err := client.Websites.Create(websiteName, websitePoolName, testAppPath)
	if err != nil {
		return fmt.Errorf("Error creating example Website: %s", err)
	}
	// Website has a binding for port *:80 by default that needs to be removed
	err = client.Websites.RemoveBinding(websiteName, "*", "", 80)
	if err != nil {
		return fmt.Errorf("Error removing binding from Website %q: %s", websiteName, err)
	}
	err = client.Websites.AddBinding(websiteName, "*", "", websitePort)
	if err != nil {
		return fmt.Errorf("Error adding a binding for port %d to website %q: %s", websitePort, websiteName, err)
	}
	return nil
}

func cleanup(client *iis.Client, applicationName string, applicationPoolName string, websiteName string, websitePoolName string, testAppPath string) {
	// Logging errors but ignoring them because we want to clean up as much as possible
	err := client.Applications.Delete(applicationName, websiteName)
	if err != nil {
		log.Printf("Error removing Application: %s", err)
	}

	err = client.AppPools.Delete(applicationPoolName)
	if err != nil {
		log.Printf("Error removing App Pool %q: %s", applicationPoolName, err)
	}

	err = client.Websites.Delete(websiteName)
	if err != nil {
		log.Printf("Error removing Website %q: %s", websiteName, err)
	}

	err = client.AppPools.Delete(websitePoolName)
	if err != nil {
		log.Printf("Error removing App Pool %q: %s", websitePoolName, err)
	}

	err = os.RemoveAll(testAppPath)
	if err != nil {
		log.Printf("Error removing test directory: %s", err)
	}
}

func testApplication(testUrl string, testAppContent string) error {
	resp, err := http.Get(testUrl)
	if err != nil {
		return fmt.Errorf("Error trying to get the test HTML from IIS: %s", err)
	}

	respContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error getting HTTP response body: %s", err)
	}

	log.Printf("Got test HTML content: %q", respContent)

	if testAppContent != string(respContent) {
		return fmt.Errorf("Content received from HTTP request doesn't match the expected content")
	}

	return nil
}

func createApplication(client *iis.Client, applicationName string, appPoolName string, testAppPath string, websiteName string) error {
	err := client.Applications.Create(applicationName, websiteName, appPoolName, testAppPath)
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
	defer testFile.Close()
	_, err = testFile.WriteString(testAppContent)
	if err != nil {
		_ = os.RemoveAll(testAppPath)
		return fmt.Errorf("Error writing test HTML file: %s", err)
	}
	log.Printf("Wrote test HTML content: %q", testAppContent)

	return nil
}
