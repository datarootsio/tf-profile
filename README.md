# tf-profile

[![Go Linting, Verification, and Testing](https://github.com/QuintenBruynseraede/tf-profile/actions/workflows/go-fmt-vet-tests.yml/badge.svg?branch=main)](https://github.com/QuintenBruynseraede/tf-profile/actions/workflows/go-fmt-vet-tests.yml)

CLI tool to profile Terraform runs, written in Go


## Usage

Basic usage

```
terraform plan | tf-profile 
terraform plan -no-color > log.txt && tf-profile log.txt  # Read from file 
terraform plan | tf-profile --stats  # Print global stats only
```

tf-profile can handle different types of logs 

```
terraform apply -auto-approve | tf-profile  # Include plan and apply phases when gathering stats
TF_LOG=TRACE && terraform apply -auto-approve | tf-profile  # Use detailed trace logs
```

Limit which resources are included

```
terraform plan | tf-profile --target=module.mymodule  # Only include logs for certain resources
terraform plan | tf-profile --max_depth=1  # Only profile the root module, aggregating metrics from submodules
```


