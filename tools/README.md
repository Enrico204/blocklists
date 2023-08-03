# IP List Downloader and Merger

This Go project provides tools that allows users to manage IP lists based on the YAML IP list index present in the root directory of the repository. These tools streamline the process of fetching, cleaning and combining multiple IP lists into a single, unified file for use with various networking and security applications.

Currently, there are these tools:
* `analyze-traffic`: match IP addresses from a PCAP or live capture against IP lists (to analyze IP lists and matches offline);
* `fetch-blocklists`: fetches all IP lists configured in YAML files (such as those in the repository root), filter and clean the result to have simple IP lists to load in `ipset` or other tools;
* `merge-blocklists`: to merge multiple IP lists together in one file, aggregating prefixes;
* `dnsbl`: a CoreDNS plugin that implements a DNSBL server querying IP lists.

There is a README file in each executable under `cmd/` with information on what it does, how, and how to use it.

The project is a work in progress.

## Usage

Pre-built binaries are not available yet. In the meantime, you need a Go 1.19 compiler for all projects, plus gcc and
libpcap headers and objects for `analyze-traffic`.

To build static binaries:

```sh
CGO_ENABLED=0 go build -o fetch-blocklists -a -ldflags '-extldflags "-static"' ./cmd/fetch-blocklists/ 
```

## Project structure

This project contains multiple executables. They are defined in `cmd/`. Note that `dnsbl` is defined as separate Go
package as it requires a large number of dependencies.

The `internal/` directory contains all internal packages: `internal/filters/` are "filters" used to download and clean
IP lists from their source.

## License

This project is licensed as specified in the COPYING file, at the root of the repository.
