package ip9datacenters

import (
	"git.netsplit.it/enrico204/blocklists/tools/internal/blocklists"
	"os"
)

type ip6datacenters struct {
	cacheDir string
}

func New(cacheDirectory string) blocklists.BlocklistFilter {
	if err := os.MkdirAll(cacheDirectory, 0750); err != nil {
		panic(err)
	}

	return &ip6datacenters{
		cacheDir: cacheDirectory,
	}
}
