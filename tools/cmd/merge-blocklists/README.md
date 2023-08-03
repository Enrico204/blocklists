# IP List merger

This tool read IP lists specified in YAML files and create a merge. It aggregates CIDRs to reduce the size of the final list.

## Usage

The tool requires a YAML configuration with the indication of which lists to use to produce a merge. The YAML contains
the following keys:

* `lists_directory`: indicates the directory where ip lists (downloaded by `fetch-blocklists`) are stored.
* `lists`: definition of merges. Each key under `lists` is a merge definition, which includes one or more IP lists.

The example `cfg/firehol.yaml` contains the definition of multiple merged lists (which roughly resembles the one from
the  FireHOL project).

```sh
$ ./merge-blocklists -h

Merges multiple blocklists, aggregating prefixes.

Usage of ./merge-blocklists:
  -cfg string
        Configuration file. (default "cfg/firehol.yaml")
  -quiet
        Print only warnings and errors to stderr.
  -verbose
        Verbose log: prints more information about what the program is doing.
```
