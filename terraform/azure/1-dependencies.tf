locals {
  prefix               = "tomdev"
  virtual_machine_name = "${local.prefix}-vm"
  admin_username       = "testadmin"
  admin_password       = "Password1234!"

  tags = {
    "Acceptance" = "Testing"
  }
}

resource "azurerm_resource_group" "main" {
  name     = "${local.prefix}-resources"
  location = "West Europe"
  tags     = "${local.tags}"
}

resource "azurerm_virtual_network" "main" {
  name                = "${local.prefix}-network"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.main.location}"
  resource_group_name = "${azurerm_resource_group.main.name}"
  tags                = "${local.tags}"
}

resource "azurerm_subnet" "internal" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.main.name}"
  virtual_network_name = "${azurerm_virtual_network.main.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "main" {
  name                         = "${local.prefix}-publicip"
  resource_group_name          = "${azurerm_resource_group.main.name}"
  location                     = "${azurerm_resource_group.main.location}"
  public_ip_address_allocation = "static"
  tags                         = "${local.tags}"
}

resource "azurerm_network_interface" "main" {
  name                = "${local.prefix}-nic"
  location            = "${azurerm_resource_group.main.location}"
  resource_group_name = "${azurerm_resource_group.main.name}"

  ip_configuration {
    name                          = "configuration"
    subnet_id                     = "${azurerm_subnet.internal.id}"
    private_ip_address_allocation = "dynamic"
    public_ip_address_id          = "${azurerm_public_ip.main.id}"
  }
}
