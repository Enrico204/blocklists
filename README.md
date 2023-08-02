# IP (block)lists

This repository contains an index of various IP lists that can be used in firewalls/IDS/IPS to block or detect attacks and other problems. This repo also contains a tool for downloading and processing them.

Lists have been divided in multiple files, one for each category:
* abuse
* anonymizers
* attacks
* geo
* malware
* organizations
* reputation
* scanners
* spam
* unroutable

Each entry in a category have some attributes (like `update_every`) that you can use to optimize the download.

The `tools/` directory contains a Go project with a tool to fetch all lists and merge them, so that you can import them in `ipset` or other tools. [See the README inside the `tools/` directory to discover how to use it.](tools/README.md)

I started this repository from the FireHOL IP lists index, as the FireHOL project development seems to be on hiatus. [Further details on differences between this index and FireHOL is in the FIREHOL.md file.](FIREHOL.md)

## YAML structure

Each entry in the categories section has this structure:

```yaml
tagname:
    filter: filter-name
    update_every: 1h
    url: https://www.example.com/list-url
    info: 'Some description'
    maintainer: Acme corp
    maintainer_url: https://www.example.com
    disabled_reason: "Website is broken."
```

* The `tagname` is the name of the list. It should be something compatible with `ipset` and filesystem names for better compatibility (so, the name should use characters in `[a-zA-Z0-9_-]`)
* The `update_every` indicates the time of validity for a given IP list. Tools updating the blocklist should never refresh the blocklist in less than this period. It is expressed in a form that is compatible with Go `time.Duration` format. Example: `1h30m10s` for 1 hour, 30 minutes and 10 seconds.
* The `url` is the URL for the list.
* `info` contains the description of the IP list.
* `maintainer` and `maintainer_url` are fields populated with the name and the URL of the maintainer.
* `disabled_reason`, **if present and not empty**, denote that the list should be temporarily skipped. The value of the key is the indication of why.
* The `filter` is the name of the filter used to process the IP list. It can be used to discriminate between different formats. See the source of the fetch tool in `tools/`.

## Contributing

Contributions to this repository are welcome! If you want to add/remove/change some list, please follow the usual steps:

1. Fork this repository and create a new branch for your contribution;
2. Add/change the list entry in the relevant YAML file, and commit with descriptive commit messages;
3. Push your changes to your forked repository;
4. Submit a pull request from your branch to the main repository.

## What is/was FireHOL?

FireHOL is a powerful yet easy-to-use `iptables`/netfilter configuration tool. The `update-ipsets` script is a component of FireHOL that facilitates the automatic updating of IP sets, which are used in conjunction with `iptables` to allow or deny traffic based on IP addresses.

Unfortunately, FireHOL development seems to be stopped, and the index of blocklists in `update-ipsets` are not maintained anymore (although the script is still working). See <https://github.com/firehol/blocklist-ipsets/issues/263>.

Please note that **this repository is not the complete FireHOL project**.
