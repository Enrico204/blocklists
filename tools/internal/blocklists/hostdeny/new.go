package hostdeny

import (
	"git.netsplit.it/enrico204/blocklists/tools/internal/blocklists"
	"os"
)

type hostdeny struct {
	cacheDir string
}

func New(cacheDirectory string) blocklists.BlocklistFilter {
	if err := os.MkdirAll(cacheDirectory, 0750); err != nil {
		panic(err)
	}

	return &hostdeny{
		cacheDir: cacheDirectory,
	}
}
