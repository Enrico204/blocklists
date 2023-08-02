package gpfcomics

import (
	"git.netsplit.it/enrico204/blocklists/tools/internal/blocklists"
	"os"
)

type gpfcomics struct {
	cacheDir string
}

func New(cacheDirectory string) blocklists.BlocklistFilter {
	if err := os.MkdirAll(cacheDirectory, 0750); err != nil {
		panic(err)
	}

	return &gpfcomics{
		cacheDir: cacheDirectory,
	}
}
