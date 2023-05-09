# Stats

**Syntax:** `tf-profile stats [options] [log_file]`

**Description:** reads a Terraform log file and print high-level information about the run.

**Options:**
- -t, --tee: print logs while parsing them. Shorthand for `terraform apply | tee >(tf-profile stats)`

**Arguments:**

Log_file: Optional. Instruct `tf-profile` to read input from a text file instead of stdin. 

## Description

The following statistics will be printed:

General:
- **Number of resources created**: Number of resources detected in your log. 

Duration:
- **Cumulative duration**: Cumulative duration of modifications. This is the sum of the duration of all modifications in the logs. Because Terraform modifies resources in parallel, this will typically be more than the actual wall time.
- **Longest apply time**: Lngest time it took to modify a single resource. The next metric shows which resource that was.
- **Longest apply resource**: the name of the resource that took the longest to modify.

Resource status:
- **No. resources in state \<STATE\>**: This statistic shows per state how many resources are in that state after the modifications.

Modules:
- **Number of top-level modules**: Number of modules called in the root module.
- **Largest top-level module**: Name of the largest top-level module.
- **Size of largest top-level module**: Number of resources in this largest top-level module. Note that this number includes all resources in submodules as well.
- **Deepest module**: Name of the deepest nested module. For example. `module.a.module.b` is two levels deep, but `module.a.module.b.module.c` is three levels deep. If multiple modules are equally as deep, the first one detected in the log will be printed.
- **Deepest module depth**: The depth of the module in the previous statistic. 
- **Largest leaf module**: A module is considered a "leaf module", if it does not make any recursive module calls. This metric prints the name of the largest leaf module.
- **Size of largest leaf module**: Number of resources in the largest leaf module. As a leaf module has no submodules, these are only the resources created directly inside this leaf module.

