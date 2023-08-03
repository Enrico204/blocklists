package apacheconf

import (
	"git.netsplit.it/enrico204/blocklists/tools/internal/filters"
	"os"
)

type apacheconf struct {
	cacheDir string
}

func New(cacheDirectory string) filters.Filter {
	if err := os.MkdirAll(cacheDirectory, 0750); err != nil {
		panic(err)
	}

	return &apacheconf{
		cacheDir: cacheDirectory,
	}
}
