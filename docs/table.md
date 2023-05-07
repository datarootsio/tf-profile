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
resource                            n  tot_time  idx_creation  idx_created  status    
-------------------------------------------------------------------------------------- 
time_sleep.count[*]                 5  15s     0             13           AllCreated  
time_sleep.foreach[*]               3  1m30s      4             11           AllCreated  
module.test[1].time_sleep.count[*]  3  6      3             9            AllCreated  
module.test[0].time_sleep.count[*]  3  3s      9             7            AllCreated 
```

The column names are lowercase and separated by underscores to allow for easy referencing in the `--sort` option. The meaning of each column is:

- **resource**: Name of the resource. In case a resource is created by a `for_each` or `count` statement, resources are aggregated and individual resource identifiers are replaced by an asterisk (*). See 
- **n**: Number of resources represented by this resource name. For regular resources, this will be 1. For resourced created with `for_each` or `count`, this number represents the number of resources created in that loop.
- **tot_time**: Total cumulative time of all resources identified by this resource name. This is typically higher than the actual wall time, as Terraform can modify multiple resources at the same time.
- **idx_creation**: order in which resource creation _started_. This means that Terraform started by creation the resource with `idx_creation = 0`. That does not guarantee the creation of this resource finished first as well (see `idx_created`).
- **idx_created**: order in which resource creation _ended_. this means that the resource with `idx_created = 0` was the first resource to be fully creatd.
- **status**: For single resources, status can be any of: `Started|NotStarted|Created|Failed`. For aggregated resources, status can be any of: `AllCreated|AllFailed|SomeFailed|NoneStarted|AllStarted|SomeStarted`.
   
    With resource aggregation, more informative statuses have precedence over less informative statuses. For example, `AllCreated` will be shown over `AllStarted`.

## Sorting

Any of the columns above can be used to sort the output table, by means of the `--sort` (shorthand `-s`) option. This option follows the format `column1:(asc|desc),column2:(asc|desc):...`. For example:
- `tot_time=desc,resource=asc`: sort first by total modification time (showing the highest first). For entries with the same modification time, sort alphabetically.
- `idx_creation=asc`: sort in order of creation, showing the resources that Terraform finished modifying first.

When sorting on resource status (`status`), statuses are mapped onto integers before sorting.

- NotStarted: 0
- Started: 1
- Created: 2
- Failed: 3
- SomeStarted: 4
- AllStarted: 5
- NoneStarted: 6
- SomeFailed: 7
- AllFailed: 8
- AllCreated: 9