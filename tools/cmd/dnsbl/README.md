# DNSBL

This Go executable implements a DNSBL server using CoreDNS and a plugin. A DNSBL server is a special DNS server used to distributed IP white lists
or block lists to SMTP servers. This DNSBL server is fully compatible with IPv4 and IPv6.

The `plugin/` directory contains the CoreDNS plugin.

## Usage

```sh
$ ./dnsbl-coredns -h
Usage of dnsbl-coredns:
  -alsologtostderr
        log to standard error as well as files
  -conf string
        Corefile to load (default "Corefile")
  -dns.port string
        Default port (default "53")
  -log_backtrace_at value
        when logging hits line file:N, emit a stack trace
  -log_dir string
        If non-empty, write log files in this directory
  -logtostderr
        log to standard error instead of files
  -p string
        Default port (default "53")
  -pidfile string
        Path to write pid file
  -plugins
        List installed plugins
  -quiet
        Quiet mode (no initialization output)
  -stderrthreshold value
        logs at or above this threshold go to stderr
  -v value
        log level for V logs
  -version
        Show version
  -vmodule value
        comma-separated list of pattern=N settings for file-filtered logging
```

## Build

```sh
$ cd tools/cmd/dnsbl/
$ go build -o dnsbl-coredns .
```

## Examples:

The `Corefile` included in this directory contains a configuration example. It can be used directly:

```sh
$ ./dnsbl-coredns
.:1053
CoreDNS-1.10.1
linux/amd64, go1.20.6, 
```

Now, if we query the DNSBL with an IP that is in a list, we will see the DNSBL reply of "match":

```sh
dig +short -p 1053 0.0.0.0.bl.example.com @127.0.0.1
127.0.0.2
```

By default, `127.0.0.2` is the DNSBL reply for match.

For further information about DNSBL: https://en.wikipedia.org/wiki/Domain_Name_System_blocklist

## Caveats - known bugs

To refresh the IP list you need to restart the server. Starting a server is slow as it loads the IP lists during the
configuration phrase.
