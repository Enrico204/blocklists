package plaintext

import (
	"git.netsplit.it/enrico204/blocklists/tools/internal/blocklists"
	"os"
)

type Options struct {
	// Unzip will unzip all files in the archive and squash them as one single continuous stream
	Unzip bool

	// UseCustomSeparator enables the use of CustomSeparator (as CustomSeparator is zero by default, and zero may be a
	// valid separator)
	UseCustomSeparator bool

	// CustomSeparator is a custom separator for new lines
	CustomSeparator rune
}

type plaintext struct {
	cacheDir string
	opts     Options
}

func New(cacheDirectory string, opts Options) blocklists.BlocklistFilter {
	if err := os.MkdirAll(cacheDirectory, 0750); err != nil {
		panic(err)
	}

	return &plaintext{
		cacheDir: cacheDirectory,
		opts:     opts,
	}
}
