# Switch with a virtual machine attached to it

This example demonstrates how to create a switch and two virtual machine. The virtual machines will be attached to the switch.

## How to run

Set environment variables `HYPERV_USER` and `HYPERV_PASSWORD` or configure provider properties `user` and `password`:
```
provider "hyperv" {
	user     = "${var.username}"
	password = "${var.password}"
}
```

then run:
```
terraform init
terraform plan -out=tfplan
terraform apply tfplan
```
