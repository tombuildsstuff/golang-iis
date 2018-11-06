resource "null_resource" "run-tests" {
  provisioner "file" {
    # sync the code
    source      = "../../"
    destination = "${local.checkout_path}"

    connection {
      host     = "${azurerm_public_ip.main.ip_address}"
      user     = "${local.admin_username}"
      password = "${local.admin_password}"
      port     = 5986
      type     = "winrm"
      https    = true
      timeout  = "10m"

      # NOTE: if you're using a real certificate, rather than a self-signed one, you'll want this set to `false`/to remove this.
      insecure = true
    }
  }

  provisioner "remote-exec" {
    inline = [
      "cd ${local.checkout_path}",
      "go version",
      "make acctests",
    ]

    connection {
      host     = "${azurerm_public_ip.main.ip_address}"
      user     = "${local.admin_username}"
      password = "${local.admin_password}"
      port     = 5986
      type     = "winrm"
      https    = true
      timeout  = "10m"

      # NOTE: if you're using a real certificate, rather than a self-signed one, you'll want this set to `false`/to remove this.
      insecure = true
    }
  }

  depends_on = ["azurerm_virtual_machine.main"]
}
