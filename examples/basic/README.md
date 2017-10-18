# Switch with a virtual machine attached to it

This example demonstrates how to create a switch and a virtual machine. The virtual machine will be attached to the switch.

## How to run

Either `cp terraform.template.tfvars terraform.tfvars` and modify that new file accordingly or provide variables via CLI:

```
terraform init
terraform plan -out=tfplan
terraform apply tfplan\
	-var="prod_access_key=AAAAAAAAAAAAAAAAAAA" \
	-var="prod_secret_key=SuperSecretKeyForAccountA" \
	-var="test_account_id=123456789012" \
	-var="test_access_key=BBBBBBBBBBBBBBBBBBB" \
	-var="test_secret_key=SuperSecretKeyForAccountB" \
	-var="bucket_name=tf-bucket-in-prod" \
```
