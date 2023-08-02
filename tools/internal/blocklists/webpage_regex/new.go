package webpagerx

import (
	"git.netsplit.it/enrico204/blocklists/tools/internal/blocklists"
	"os"
	"regexp"
)

type webpageregex struct {
	cacheDir string
	matcher  *regexp.Regexp
}

func New(cacheDirectory string, matcher *regexp.Regexp) blocklists.BlocklistFilter {
	if err := os.MkdirAll(cacheDirectory, 0750); err != nil {
		panic(err)
	}

	return &webpageregex{
		cacheDir: cacheDirectory,
		matcher:  matcher,
	}
}
