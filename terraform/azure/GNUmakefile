provision:
	terraform init
	terraform apply --auto-approve

tests:
	terraform taint null_resource.run-tests
	terraform apply --auto-approve
