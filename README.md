# golang-iis

Enables management of Windows IIS from Golang.

### Why?

Good question. IIS is a bit of a unmanageable beast - this library tries to make it friendlier, for if you were doing something crazy like automating it.

## Requirements

- Go 1.10+ (possibly earlier, but untested)
- Access to a machine with IIS

### Example Usage

```go
import (
  "log"

  "github.com/tombuildsstuff/golang-iis/iis"
)

func main() {
    log.Printf("Example app launched..")

    websiteName := "my-website"
    appPoolName := "my-app-pool"
    physicalPath := "C:\\inetpub\\wwwroot"

    err := run(websiteName, appPoolName, physicalPath)
    if err != nil {
        // handle the error better than this
        panic(err)
    }
}

func run(websiteName string, appPoolName string, physicalPath string) error {
    client, err := iis.NewClient()
    if err != nil {
        return err
    }

    log.Printf("Creating App Pool %q..", appPoolName)
    err = client.AppPools.Create(appPoolName)
    if err != nil {
        return fmt.Errorf("Error creating App Pool %q: %+v", appPoolName, err)
    }

    log.Printf("Creating Website %q exists..", websiteName)
    err = client.Websites.Create(websiteName, appPoolName, physicalPath)
    if err != nil {
        return fmt.Errorf("Error creating Website %q: %+v", websiteName, err)
    }

    return nil
}
```

More examples can be found in [the `./examples` directory](https://github.com/tombuildsstuff/golang-iis/tree/master/examples) in the root of this repository.

## Running the Tests

If you're running Windows and have IIS installed, you should be able to run `go test -v ./... -parallel=1` and things should run. This is significantly quicker than provisioning the Virtual Machine since the files don't need to be replicated.

If you're not running Windows - there's a [HashiCorp Terraform](https://terraform.io) script available which will provision a Windows VM (currently only in Azure, but other providers would be cool) which runs the tests remotely [which can be found here](https://github.com/tombuildsstuff/golang-iis/tree/master/terraform).

## Licence

MIT
