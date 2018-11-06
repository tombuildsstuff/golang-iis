package main

import (
	"fmt"
	"log"

	"github.com/tombuildsstuff/golang-iis/iis/helpers"

	"github.com/tombuildsstuff/golang-iis/iis"
)

func main() {
	log.Printf("Running sample..")
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	log.Printf("Creating CLient..")
	client, err := iis.NewClient()
	if err != nil {
		return fmt.Errorf("Error building client: %s", err)
	}

	rInt := helpers.RandomInt()
	appPoolName := fmt.Sprintf("app-pool-%d", rInt)
	websiteName := fmt.Sprintf("website-%d", rInt)
	physicalPath := "C:\\inetpub\\wwwroot"

	err = createAppPool(client, appPoolName)
	if err != nil {
		return err
	}

	err = createWebsite(client, websiteName, appPoolName, physicalPath)
	if err != nil {
		return err
	}

	log.Printf("Deleting Website %q..", websiteName)
	err = client.Websites.Delete(websiteName)
	if err != nil {
		return fmt.Errorf("Error deleting website %q: %s", websiteName, err)
	}

	log.Printf("Deleting App Pool %q..", appPoolName)
	err = client.AppPools.Delete(appPoolName)
	if err != nil {
		return fmt.Errorf("Error deleting App Pool %q: %s", appPoolName, err)
	}

	return nil
}

func createAppPool(client *iis.Client, name string) error {
	log.Printf("Creating the App Pool (with Name %q)..", name)
	err := client.AppPools.Create(name)
	if err != nil {
		return fmt.Errorf("Error creating App Pool %q: %+v", name, err)
	}

	log.Printf("Checking if the App Pool exists..")
	exists, err := client.AppPools.Exists(name)
	if err != nil {
		return fmt.Errorf("Error checking if App Pool %q exists: %s", name, err)
	}

	if !*exists {
		return fmt.Errorf("App Pool %q didn't exist when it was supposed to: %s", name, err)
	}

	return nil
}

func createWebsite(client *iis.Client, websiteName string, appPoolName string, physicalPath string) error {
	appSettingName := "SomeAppSetting"
	appSettingValue := "SomeValue"
	hostHeader := "example.com"

	log.Printf("Creating Website (with Name %q)..", websiteName)
	err := client.Websites.Create(websiteName, appPoolName, physicalPath)
	if err != nil {
		return fmt.Errorf("Error creating Website %q (App Pool %q): %s", websiteName, appPoolName, err)
	}

	log.Printf("Setting the App Setting %q to %q..", appSettingName, appSettingValue)
	err = client.Websites.SetAppSetting(websiteName, appSettingName, appSettingValue)
	if err != nil {
		return fmt.Errorf("Error setting AppSetting %q (Website %q): %s", appSettingName, websiteName, err)
	}

	log.Printf("Retrieving the AppSetting %q", appSettingName)
	val, err := client.Websites.GetAppSetting(websiteName, appSettingName)
	if err != nil {
		return fmt.Errorf("Error retrieving AppSetting %q (Website %q): %s", appSettingName, websiteName, err)
	}

	if *val != appSettingValue {
		return fmt.Errorf("Expected the AppSetting %q to have value %q but got %q", appSettingName, appSettingValue, *val)
	}

	log.Printf("Adding a binding for %q..", hostHeader)
	err = client.Websites.AddBinding(websiteName, "*", hostHeader, 80)
	if err != nil {
		return fmt.Errorf("Error adding binding: %s", err)
	}

	log.Printf("Retrieving Bindings for Website %q..", websiteName)
	bindings, err := client.Websites.GetBindings(websiteName)
	if err != nil {
		return fmt.Errorf("Error retrieving bindings for Website %q: %s", websiteName, err)
	}

	found := false
	for _, binding := range *bindings {
		if binding.DomainName == hostHeader {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("Expected Website %q to have a binding for %q but didn't find one..", websiteName, hostHeader)
	}

	log.Printf("Removing binding %q from Website %q..", hostHeader, websiteName)
	err = client.Websites.Delete(websiteName)
	if err != nil {
		return fmt.Errorf("Error removing binding %q from website %q: %s", hostHeader, websiteName, err)
	}

	return nil
}
