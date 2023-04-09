# tf-profile

[![Go Linting, Verification, and Testing](https://github.com/QuintenBruynseraede/tf-profile/actions/workflows/go-fmt-vet-tests.yml/badge.svg?branch=main)](https://github.com/QuintenBruynseraede/tf-profile/actions/workflows/go-fmt-vet-tests.yml)

CLI tool to profile Terraform runs, written in Go.

Main features:
- Modern CLI ([cobra](https://github.com/spf13/cobra)-based), including autocomplete
- Read logs straight from your Terraform process (using pipe) or a log file
- Can generate global stats, resource-level stats or vizualization
- Provides many levels of granularity and aggregation, customizable outputs

## Installation

For now, the only supported way is to build the binary yourself. This requires at least version 1.20 of the `go` cli.

```bash
$ git clone git@github.com:QuintenBruynseraede/tf-profile.git
$ cd tf-profile && go build .
$ sudo ln -s $(pwd)/tf-profile /usr/local/bin  # Optional: only if you want to run tf-profile from other directories
$ tf-profile --help
tf-profile is a CLI tool to profile Terraform runs

Usage:
  tf-profile [command]
```

## Basic usage

`tf-profile` handles input from stdin and from files. These two commands are therefore equivalent:

```bash
$ terraform apply -auto-approve | tf-profile table
$ terraform apply -auto-approve > log.txt && tf-profile table log.txt
```

Three major commands are supported:
- [🔗](#anchor_stats) `tf-profile stats`: provide general statistics about a Terraform run
- [🔗](#anchor_table) `tf-profile table`: provide detailed, resource-level statistics about a Terraform run
- [🔗](#anchor_graph) `tf-profile graph`: generate visual overview of a Terraform run.


## `tf-profile stats`
<a name="anchor_stats"></a>

`tf-profile stats` is the most basic command. Given a Terraform log, it will only provide high-level statistics.

Example:
```bash
$ terraform apply -auto-approve > log.txt
$ tf-profile stats log.txt
> tf-profile stats test/many_modules.log

Key                                Value    
-----------------------------------------------------------------                       
Number of resources created        1510                            
                                                                   
Cumulative duration                36m19s                          
Longest apply time                 7m18s                           
Longest apply resource             time_sleep.foo[*]               
                                                                   
No. resources in state AllCreated  800                             
No. resources in state Created     695                             
No. resources in state Started     15                              
                                                                   
Number of top-level modules        13                              
Largest top-level module           module.core[2]                  
Size of largest top-level module   170                             
Deepest module                     module.core[2].module.role[47]  
Deepest module depth               2                               
Largest leaf module                module.dbt[4]                   
Size of largest leaf module        40  
```

Statistics are divided into four categories:
1. Basic:
    - Number of resources created
2. Related to time:
    - Cumulative duration of all resource creation
    - Which resource took the longest to create (and how long)
3. Related to resource statuses
    - For each possible status (see `tf-profile table`), the number of resources
4. Module statistics:
    - Number of root modules, largest root module
    - Nested module depth, name and size of the deepest module
    - Largest leaf module and its size

## `tf-profile table`
<a name="anchor_table"></a>
`tf-profile table` will parse a log and provide per-resource metrics. By default, resources created with `for_each` and `count` are aggregated into one entry (e.g. `aws_subnet[0]` and `aws_subnet[1]` become `aws_subnet[*]`). The following statistics are shown:

- **resource**: Resource name
- **n**: Number of resources created (usually 1, unless `count` or `for_each` were used)
- **tot_time** (milliseconds): total time spent creating these resources. Note that for resources where `n` > 1, `tot_time` does not equal wall time, as Terraform usually creates resources in parallel. There is currently no way to accurately find the wall time from a Terraform log.
- **idx_creation**: order in which resource creation _started_. This means that Terraform started by creation the resource with `idx_creation = 0`. That does not guarantee the creation of this resource finished first as well (see `idx_created`).
- **idx_created**: order in which resource creation _ended_. this means that the resource with `idx_created = 0` was the first resource to be fully creatd.
- **status**: For single resources, status can be any of: `Started|NotStarted|Created|Failed`. For aggregated resources, status can be any of: `AllCreated|AllFailed|SomeFailed|NoneStarted|AllStarted|SomeStarted`.
   
    With resource aggregation, more informative statuses have precedence over less informative statuses. For example, `AllCreated` will be shown over `AllStarted`.
```bash
$ terraform apply -auto-approve > log.txt
$ tf-profile table log.txt

resource                            n  tot_time  idx_creation  idx_created  status    
-------------------------------------------------------------------------------------- 
time_sleep.count[*]                 5  11000     0             13           AllCreated  
time_sleep.foreach[*]               3  7000      4             11           AllCreated  
module.test[1].time_sleep.count[*]  3  5000      3             9            AllCreated  
module.test[0].time_sleep.count[*]  3  4000      9             7            AllCreated 
```

### Sorting the table with `--sort`

Entries in this table can be sorted by providing a `--sort` (shorthand `-s`) argument. This argument is a comma-separated list of key-value pairs. Valid example include:
- `tot_time=desc` (default): sort by total time descending
- `tot_time=asc,resource=desc,status=asc`: sort by total time, resource name and creation status in that order

```bash
$ terraform apply -auto-approve > log.txt
$ tf-profile table --sort "tot_time=asc,resource=desc" log.txt

resource                            n  tot_time  idx_creation  idx_created  status      
--------------------------------------------------------------------------------------
module.test[0].time_sleep.count[*]  3  4000      9             7            AllCreated  
module.test[1].time_sleep.count[*]  3  5000      3             9            AllCreated  
time_sleep.foreach[*]               3  7000      4             11           AllCreated  
time_sleep.count[*]                 5  11000     0             13           AllCreated 
```

### Mirroring input with `--tee`

When piping stdinput into `tf-profile`, it is convenient to use the `--tee` flag. This flag instructs `tf-profile` to print every line it parses. This way you don't lose your detailed Terraform logs, but still get a table at the end.

Example:
```bash
$ terraform apply -auto-approve | tf-profile table --tee

Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # aws_subnet.test[0] will be created
  ...

resource                            n  tot_time  idx_creation  idx_created  status    
-------------------------------------------------------------------------------------- 
time_sleep.count[*]                 5  11000     0             13           AllCreated  
time_sleep.foreach[*]               3  7000      4             11           AllCreated  
module.test[1].time_sleep.count[*]  3  5000      3             9            AllCreated  
module.test[0].time_sleep.count[*]  3  4000      9             7            AllCreated 
```

### Limit output with `--max_depth`
🚧 Under construction (not implemented) 🚧

When working with a large codebase, viewing statistics for every resource may be too detailed. You can limit the maximum module depth that `tf-profile` parses with `--max_depth` (default: -1, no limit). Any nested modules deeper than `--max_depth` are simply shown as their module name. Statistics of resources within that module are aggregated.

```bash
$ terraform apply -auto-approve | tf-profile table --max_depth 1 --tee

resource                            n  tot_time  idx_creation  idx_created  status    
-------------------------------------------------------------------------------------- 
time_sleep.count[*]                 5  11000     0             13           AllCreated  
time_sleep.foreach[*]               3  7000      4             11           AllCreated  
module.test[1]                      3  5000      3             9            AllCreated  
module.test[0]                      3  4000      9             7            AllCreated 
```

## `tf-profile graph`
<a name="anchor_graph"></a>
🚧 Under construction (not implemented) 🚧
