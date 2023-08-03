package hostdeny

import (
	"git.netsplit.it/enrico204/blocklists/tools/internal/filters"
	"os"
)

type hostdeny struct {
	cacheDir string
}

func New(cacheDirectory string) filters.Filter {
	if err := os.MkdirAll(cacheDirectory, 0750); err != nil {
		panic(err)
	}

	return &hostdeny{
		cacheDir: cacheDirectory,
	}
}
