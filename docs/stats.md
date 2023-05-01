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
- Number of resources modified

Duration:
- Cumulative duration of modifications. This is the sum of the duration of all modifications in the logs. Because Terraform modifies resources in parallel, this will typically be more than the actual wall time.
- 

