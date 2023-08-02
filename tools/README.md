# IP List Downloader and Merger

This Go project provides a tool that allows users to download and merge IP lists based on the YAML IP list index present in the root directory of the repository. The tool streamlines the process of fetching, cleaning and combining multiple IP lists into a single, unified file for use with various networking and security applications.

The project is a work in progress.

## Features

* Download multiple IP lists from a YAML index, **respecting the Last-Modified header and the update interval specified in the YAML**.
* Merge together multiple IP lists specified in a YAML configuration file, aggregating prefixes.
* IPv4 and IPv6 compatibility.
* Lightweight and written in Go for efficient performance.

## Usage

Pre-built binaries are not available yet. In the meantime, you need a Go 1.19 compiler. You can download the source code, and execute these commands inside the `tools/` directory:

```sh
$ go build -o fetch-blocklists ./cmd/fetch-blocklists
$ go build -o merge-blocklists ./cmd/merge-blocklists
```

Then, you can launch them as command line tools:

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

An example of the merger configuration is in `cfg/firehol.yaml`. The example produces a set of IP lists similar to the FireHOL ones.

## License

This project is licensed as specified in the COPYING file, at the root of the repository.
