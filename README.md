# tf-profile

[![Go Linting, Verification, and Testing](https://github.com/QuintenBruynseraede/tf-profile/actions/workflows/go-fmt-vet-tests.yml/badge.svg?branch=main)](https://github.com/QuintenBruynseraede/tf-profile/actions/workflows/go-fmt-vet-tests.yml)

CLI tool to profile Terraform runs, written in Go.

Main features:
- Modern CLI ([cobra](https://github.com/spf13/cobra)-based), including autocomplete
- Read logs straight from your Terraform process (using pipe) or a log file
- Can generate global stats, resource-level stats or visualizations
- Provides many levels of granularity and aggregation, customizable outputs

## Installation

### Binary download

- Head over to the releases page ([https://github.com/QuintenBruynseraede/tf-profile/releases](https://github.com/QuintenBruynseraede/tf-profile/releases)) 
- Download the correct binary for your operating system
- Copy it to a path that is on your `$PATH`. On a Linux system, `/usr/local/bin` is the most common location.

### Using docker

If you want to try `tf-profile` without installing anything, you can run it using Docker (or similar).

```bash
❱ cat my_log_file.log | docker run -i qbruynseraede/tf-profile:0.0.1 stats

Key                                Value                                     
Number of resources created        1510                                      
                                                                             
Cumulative duration                36m19s                                    
Longest apply time                 7m18s                                     
Longest apply resource             time_sleep.foo[*]                         
...
```

Optionally, define an alias:

```bash
❱ alias tf-profile=docker run -i qbruynseraede/tf-profile:0.0.1
❱ cat my_log_file.log | tf-profile
```

### Build from source

This requires at least version 1.20 of the `go` cli.

```bash
❱ git clone git@github.com:QuintenBruynseraede/tf-profile.git
❱ cd tf-profile && go build .
❱ sudo ln -s $(pwd)/tf-profile /usr/local/bin  # Optional: only if you want to run tf-profile from other directories
❱ tf-profile --help
tf-profile is a CLI tool to profile Terraform runs

Usage:
  tf-profile [command]
```

## Basic usage

`tf-profile` handles input from stdin and from files. These two commands are therefore equivalent:

```bash
❱ terraform apply -auto-approve | tf-profile table
❱ terraform apply -auto-approve > log.txt && tf-profile table log.txt
```

Three major commands are supported:
- [🔗](#anchor_stats) `tf-profile stats`: provide general statistics about a Terraform run
- [🔗](#anchor_table) `tf-profile table`: provide detailed, resource-level statistics about a Terraform run
- [🔗](#anchor_graph) `tf-profile graph`: generate visual overview of a Terraform run.


## `tf-profile stats`
<a name="anchor_stats"></a>

`tf-profile stats` is the most basic command. Given a Terraform log, it will only provide high-level statistics.

```bash
❱ terraform apply -auto-approve > log.txt
❱ tf-profile stats log.txt
❱ tf-profile stats test/many_modules.log

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

For more information, refer to the [reference](./docs/stats.md) for the `stats` command.

## `tf-profile table`
<a name="anchor_table"></a>
`tf-profile table` will parse a log and provide per-resource metrics.

```bash
❱ terraform apply -auto-approve > log.txt
❱ tf-profile table log.txt

resource                            n  tot_time  idx_creation  idx_created  status    
-------------------------------------------------------------------------------------- 
time_sleep.count[*]                 5  11s     0             13           AllCreated  
time_sleep.foreach[*]               3  7s      4             11           AllCreated  
module.test[1].time_sleep.count[*]  3  5s      3             9            AllCreated  
module.test[0].time_sleep.count[*]  3  4s      9             7            AllCreated 
```

For a full description of the options, see the [reference](./docs/table.md) page.


## `tf-profile graph`
<a name="anchor_graph"></a>

`tf-profile graph` is used to visualize your terraform logs. It generates a [Gantt](https://en.wikipedia.org/wiki/Gantt_chart)-like chart that shows in which order resources were created. `tf-profile` does not actually create the final image, but generates a script file that [Gnuplot](https://en.wikipedia.org/wiki/Gnuplot) understands. 

```bash
❱ tf-profile graph my_log.log --out graph.png --size 2000,1000 | gnuplot
```

![graph.png](https://github.com/QuintenBruynseraede/tf-profile/blob/main/.github/graph.png?raw=true)

_Disclaimer:_ Terraform's logs do not contain any absolute timestamps. We can only derive the order in which resources started and finished their modifications. Therefore, the output of `tf-profile graph` gives only a general indication of _how long_ something actually took. In other words: the X axis is meaningless, apart from the fact that it's monotonically increasing.


## Screenshots

![stats.png](https://github.com/QuintenBruynseraede/tf-profile/blob/main/.github/stats.png?raw=true)

![table.png](https://github.com/QuintenBruynseraede/tf-profile/blob/main/.github/table.png?raw=true)

![graph2.png](https://github.com/QuintenBruynseraede/tf-profile/blob/main/.github/graph2.png?raw=true)

## Roadmap

- [x] Release v0.0.1 as binary and as a Docker image
- [ ] Improve parser
  - [x] Detect failed resources (see [#13](https://github.com/QuintenBruynseraede/tf-profile/pull/13))
  - [ ] Use plan and refresh phase to discover more resources
- [x] Implement a basic Gantt chart in `tf-profile graph` (see [#14](https://github.com/QuintenBruynseraede/tf-profile/pull/14))
- [ ] Implement a single-resource view in `tf-profile detail <resource>`
  - This command should filter logs down to 1 single resource (i.e. refresh, plan, changes, and result)
- [ ] Small improvements:
  - [ ] Add `no-agg` option to disable aggregation of for_each and count
  - [ ] Add `max_depth` option to aggregate deep submodules
  - [ ] Find a way to rename the columns in `tf-profile table` without breaking `--sort`
  - [ ] Add go report card: [https://goreportcard.com/report/github.com/QuintenBruynseraede/tf-profile](https://goreportcard.com/report/github.com/QuintenBruynseraede/tf-profile)
