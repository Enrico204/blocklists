# Traffic analyzer

This Go executable matches source and destination IPs in PCAP files or live capture against IP lists, either in
aggregated form or match by IP list. It allows you to check whether an IP list matches actual traffic before deplyoing
it into production systems.

## Usage

```sh
$ ./analyze-traffic -h

Analyze a PCAP or a live traffic and reports who many IP packets matches which IP list.

Usage of ./analyze-traffic:
  -aggregate
        Aggregate all lists together. In this case, the result will be only match/no match.
  -iface string
        Interface to listen to (required, or specify -pcap).
  -index-directory string
        Blocklists YAML index files location. (default ".")
  -lists-directory string
        Blocklists files location. (default "out/")
  -pcap string
        PCAP file to analyze (required, or specify -iface).
  -promisc
        Whether the interface should be in promiscuous mode (only when live capturing). (default true)
  -quiet
        Print only warnings and errors to stderr.
  -size int
        Capture size (only when live capturing). (default 1600)
  -verbose
        Verbose log: prints more information about what the program is doing.
```

## Build

```sh
$ go build -o analyze-traffic ./cmd/analyze-traffic/
```

## Examples:

```sh
$ go run ./cmd/analyze-traffic/ -pcap ~/smtp_brute_forces.pcap -index-directory ../ --lists-directory out/
Per blocklist results:
                          myip - SRC:         0 - DST:         0 - TOTAL:         0
                         xroxy - SRC:         0 - DST:         0 - TOTAL:         0
...
                    fullbogons - SRC:      2671 - DST:      2671 - TOTAL:      5342
                   urlvir_last - SRC:         0 - DST:         0 - TOTAL:         0
         iblocklist_exclusions - SRC:      2671 - DST:      2671 - TOTAL:      5342
                php_dictionary - SRC:         0 - DST:         0 - TOTAL:         0
```

## Performance notes

When loading all lists and match against IPv4 both source and destination in a PCAP, `analyze-traffic`:
* saturate the CPU when loading all blocklists, reaching slightly more than 1GB of RAM in aggregated mode
* saturate the CPU when loading all blocklists, reaching slightly more than 3GB of RAM in normal mode
* loading all lists requires nearly ~28 seconds on a PC with a NVMe M.2 disk, in aggregated mode
* loading all lists requires nearly ~17 seconds on a PC with a NVMe M.2 disk, in normal mode
* matching on average lasts 670ns for both source and destination combined, for a single packet, in aggregated mode
* matching on average lasts 90µs for both source and destination combined, for a single packet, for all blocklists, in normal mode
