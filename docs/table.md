# Stats

**Syntax:** `tf-profile table [options] [log_file]`

**Description:** reads a Terraform log file and print high-level information about the run.

**Options:**
- -t, --tee: print logs while parsing them. Shorthand for `terraform apply | tee >(tf-profile stats)`
- -d, --max_depth: aggregate resources nested deeper than `-d` levels into a resource that represents the module at depth `-d`.
- -s, --sort: comma-separated key-value pairs that instruct how to sort the output table. Valid values follow the format `column1:(asc|desc),column2:(asc|desc):...`. By default, `tot_time=desc,resource=asc` is used: sort first by descending modification time, second by resource name in alphabetical order.

**Arguments:**

Log_file: Optional. Instruct `tf-profile` to read input from a text file instead of stdin. 

## Description

A table generated based on the log file or input, sorted according to `-s / --sort` and printed to the terminal. 

```
resource              n  tot_time  modify_started  modify_ended  desired_state  operation  final_state  
aws_ssm_parameter.p6  1  0s        6               7             Created        Replace    Created      
aws_ssm_parameter.p1  1  0s        7               5             Created        Replace    Created      
aws_ssm_parameter.p3  1  0s        5               6             Created        Replace    Created      
aws_ssm_parameter.p4  1  0s        /               1             NotCreated     Destroy    NotCreated   
aws_ssm_parameter.p5  1  0s        4               4             Created        Modify     Created      
aws_ssm_parameter.p2  1  0s        /               /             Created        None       Created      
```

The column names are lowercase and separated by underscores to allow for easy referencing in the `--sort` option. The meaning of each column is:

- **resource**: Name of the resource. In case a resource is created by a `for_each` or `count` statement, resources are aggregated and individual resource identifiers are replaced by an asterisk (*). See 
- **n**: Number of resources represented by this resource name. For regular resources, this will be 1. For resourced created with `for_each` or `count`, this number represents the number of resources created in that loop.
- **tot_time**: Total cumulative time of all resources identified by this resource name. This is typically higher than the actual wall time, as Terraform can modify multiple resources at the same time.
- **modify_started**: order in which resource modification _started_. This means that Terraform started by modifying the resource with `modify_started = 0`. It does not guarantee the changes to this resource finished first as well (see `modify_ended`).
- **modify_ended**: order in which resource modifications _ended_. This means that the resource with `modify_ended = 0` was the first resource to finish its modifications (either a creation, deletion or replacement). Resources that were already consistent with the desired state do not have this property.
- **desired_state**: state (Created, NotCreated) that Terraform will try to achieve with this run.
- **operation**: the name of the operation the Terraform will use to reconcile the current and desired situation. Operations can be: Create, Destroy, Replace, Modify, None.
- **final_state**: Final state of the resource after this run. In addition to Created and NotCreated, Failed is used to indicate the operation failed.

## Sorting

Any of the columns above can be used to sort the output table, by means of the `--sort` (shorthand `-s`) option. This option follows the format `column1:(asc|desc),column2:(asc|desc):...`. For example:
- `tot_time=desc,resource=asc`: sort first by total modification time (showing the highest first). For entries with the same modification time, sort alphabetically.
- `idx_creation=asc`: sort in order of creation, showing the resources that Terraform finished modifying first.

When sorting on resource status (`desired_state` or `final_state`), statuses are mapped onto integers before sorting.

- Unknown: 0
- NotCreated: 1
- Created: 2
- Failed: 3
- Tainted: 4
- Multiple (for aggregated resources): 5

When sorting on resource operations (`operation`), these are mapped onto integers:

- None: 0
- Create: 1
- Modify: 2
- Replace: 3
- Destroy: 4
- Multiple (for aggregated resources): 5