# Filter

**Syntax:** `tf-profile filter [resource_filter] [log_file]`

**Description:** filter a Terraform log to only selected resources

**Arguments:**

- resource_filter: Output will only contain resources whose name matches this filter. Wildcards are supported. See below for examples on how to filter resources.
- log_file: _Optional_. Instruct `tf-profile` to read input from a text file instead of stdin. 

## Description

This command will filter a log line-by-line, and print only the lines that are related to certain resources (as specified by the resource filter). More specifically, the following lines will remain in the output:
1. Part of the Terraform plan that describes the changes to this resource.
2. The full error message if an error occurred while modifying the resource.
3. Any line that matches the resource filter, but is not part of a plan or error message.

Note that lines are matched to patterns using standard Go `regexp` functions. Out of the box, this would lead to cumbersome resource filtering (filters like `module\\.foo\\.resource\[.*\]` instead of the more natural `module.foo.resource[*]`). To support this, resource filters are transformed from the latter into the former. This entails:
- `*` is replaced by `.*`
- The following characters are escaped: `.`, `[`, `]`

## Examples

Basic usage:
```sh
❱ terraform apply -auto-approve | tf-profile filter "aws_ssm_parameter.test"
  # aws_ssm_parameter.test will be created
  + resource "aws_ssm_parameter" "test" {
      + arn            = (known after apply)
      + name           = "my_ssm_parameter"
      + type           = "String"
      + value          = (sensitive value)
      + version        = (known after apply)
    }

aws_ssm_parameter.test: Creating...
aws_ssm_parameter.test: Creation complete after 1s [id=my-param]
```

Reading from a log file:
```sh
❱ tf-profile filter "aws_ssm_parameter.test" log.txt
... # Output identical to above
```

Using wildcards to filter on multiple resources:
```sh
❱ tf-profile filter "aws_ssm_parameter.*" log.txt

  # aws_ssm_parameter.param1 will be created
  + resource "aws_ssm_parameter" "param1" {
    ...
    }

  # aws_ssm_parameter.param2 will be created
  + resource "aws_ssm_parameter" "param2" {
    ...
    }

aws_ssm_parameter.param1: Creating...
aws_ssm_parameter.param2: Creating...
aws_ssm_parameter.param1: Creation complete after 1s [id=param1]

Error: creating SSM Parameter (param2): ValidationException: Parameter name must not end with slash.
        status code: 400, request id: f7abc744-2fff-4d92-824b-b73g40ab256a

  with aws_ssm_parameter.param2,
  on provider.tf line 27, in resource "aws_ssm_parameter" "param2":
  27: resource "aws_ssm_parameter" "param2" {
```

Multiple wildcards can be combined:
```sh
❱ tf-profile filter "module.*.null_resource.*" log.txt

  # module.mod1.null_resource.foo will be created
  + resource "null_resource" "foo" {
    ...
    }

  # module.mod2.null_resource.bar will be created
  + resource "null_resource" "bar" {
    ...
    }

module.mod1.null_resource.foo: Creating...
module.mod2.null_resource.bar: Creating...
module.mod1.null_resource.foo: Creation complete after 1s [id=foo]
module.mod2.null_resource.bar: Creation complete after 1s [id=bar]
```