package csv

import (
	"git.netsplit.it/enrico204/blocklists/tools/internal/blocklists"
	"os"
)

type Options struct {
	// SkipHeader indicates whether we need to skip the first row
	SkipHeader bool

	// Column indicates in which column the IP address is
	Column int

	// Comma is the separator character for the CSV. Default: ','
	Comma rune

	// Comment, if not 0, is the comment character. Lines beginning with the
	// Comment character without preceding whitespace are ignored.
	// With leading whitespace the Comment character becomes part of the
	// field, even if TrimLeadingSpace is true.
	// Comment must be a valid rune and must not be \r, \n,
	// or the Unicode replacement character (0xFFFD).
	// It must also not be equal to Comma.
	Comment rune
}

type csv struct {
	cacheDir string
	opts     Options
}

func New(cacheDirectory string, opts Options) blocklists.BlocklistFilter {
	if err := os.MkdirAll(cacheDirectory, 0750); err != nil {
		panic(err)
	}

	if opts.Comma == 0 {
		opts.Comma = ','
	}

	return &csv{
		cacheDir: cacheDirectory,
		opts:     opts,
	}
}
