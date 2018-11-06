## Provisions a VM to run the Acceptance Tests in Azure

This directory provisions a Windows Virtual Machine running on Azure which installs IIS and Golang - required to run the acceptance tests.

### Notes

- The initial provisioning takes about ~25m to install IIS

-> **NOTE:** If you're running Windows as a Host OS you may prefer to run the Acceptance Tests on your machine, since this'll be faster.

- This script will SCP this repository onto the Virtual Machine; you may wish to compile the AzureRM and Null Providers for [HashiCorp's Terraform](https://terraform.io) so they're not uploaded to the remote machine.

### Commands

Assuming your Azure credentials are configured [as per the documentation](https://www.terraform.io/docs/providers/azurerm/index.html) and Terraform is available on your path:

You can provision the Virtual Machine by running:

```
make provision
```

You can then run (and re-run) the tests by running:

```
make test
```

-> **NOTE:** This'll re-sync the files from your local directory to the Virtual Machine, which can take a couple of minutes.