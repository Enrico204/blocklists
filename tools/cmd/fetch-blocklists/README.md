# IP List downloader

This tools downloads the blocklists specified in YAML files, clean up all files, aggregate CIDR prefixes and save the
result in a directory.


When downloading the IP lists, it uses the Last-Modified header and supports GZIP transfer compression to avoid wasting
traffic. Also, it tries to download a list no more frequent than the period specified in the configuration file.

The tool is a work in progress. IPv4 compatibility is complete; IPv6 "should work", but more testing is needed.

## Usage

The tool logs to `stderr` if there is any problem fetching or filtering a list.

```sh
$ ./fetch-blocklists -h

Fetches all IP lists configured in YAML files and aggregate them.

Usage of ./fetch-blocklists:
  -cache-directory string
        Cache directory.
        The program will store downloaded but unfiltered lists in this directory.
        It's used to avoid re-downloading the same content again and again.
        The directory will be created if it does not exist. (default "tmp/")
  -debug
        Debug mode: verbose log + dumps HTTP headers in the log.
  -lists-directory string
        Blocklists YAML files location. (default ".")
  -output-directory string
        Output directory.
        Cleaned, filtered and aggregated lists will be saved here.
        The directory will be created if it does not exist. (default "out/")
  -quiet
        Print only warnings and errors to stderr.
  -verbose
        Verbose log: prints more information about what the program is doing.
```

## Example

```sh
$ go run ./cmd/fetch-blocklists/ -lists-directory ../ > fetch.log
2023-08-03T11:24:00.294+0200	INFO	fetch-blocklists/register-filters.go:25	initializing filters
```
