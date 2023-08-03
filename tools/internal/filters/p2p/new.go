package p2p

import (
	"git.netsplit.it/enrico204/blocklists/tools/internal/filters"
	"os"
)

type p2p struct {
	cacheDir string
}

func New(cacheDirectory string) filters.Filter {
	if err := os.MkdirAll(cacheDirectory, 0750); err != nil {
		panic(err)
	}

	return &p2p{
		cacheDir: cacheDirectory,
	}
}
