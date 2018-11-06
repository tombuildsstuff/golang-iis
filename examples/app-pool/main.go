package main

import (
	"fmt"
	"log"

	"github.com/tombuildsstuff/golang-iis/iis/applicationpools"
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
	name := fmt.Sprintf("app-pool-%d", rInt)

	log.Printf("Creating the App Pool (with Name %q)..", name)
	err = client.AppPools.Create(name)
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

	log.Printf("Retrieving App Pool %q..", name)
	pool, err := client.AppPools.Get(name)
	if err != nil {
		return fmt.Errorf("Error retrieving App Pool %q: %s", name, err)
	}

	log.Printf("The current framework version is: %s", string(pool.FrameworkVersion))

	versions := []applicationpools.ManagedFrameworkVersion{
		applicationpools.ManagedFrameworkVersionNone,
		applicationpools.ManagedFrameworkVersionTwo,
		applicationpools.ManagedFrameworkVersionFour,
	}
	for _, version := range versions {
		log.Printf("Setting the Managed Runtime Version to %q..", string(version))
		err = client.AppPools.SetRuntimeVersion(name, version)
		if err != nil {
			return fmt.Errorf("Error setting the Managed Runtime Version to %q: %+v", version, err)
		}
	}

	log.Printf("Deleting App Pool..")
	err = client.AppPools.Delete(name)
	if err != nil {
		return fmt.Errorf("Error deleting App Pool %q: %s", name, err)
	}

	return nil
}
