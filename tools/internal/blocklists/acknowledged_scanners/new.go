package ackscanners

import (
	"git.netsplit.it/enrico204/blocklists/tools/internal/blocklists"
	"os"
	"path"
)

type ackscanners struct {
	cachedir string
}

func New(cacheDirectory string) blocklists.BlocklistFilter {
	cacheDirectory = path.Join(cacheDirectory, "acknowledge_scanners")

	if err := os.MkdirAll(cacheDirectory, 0750); err != nil {
		panic(err)
	}

	return &ackscanners{
		cachedir: cacheDirectory,
	}
}
