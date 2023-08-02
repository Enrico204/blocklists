package localfiles

import (
	"git.netsplit.it/enrico204/blocklists/tools/internal/blocklists"
)

type localfiles struct {
}

func New() blocklists.BlocklistFilter {
	return &localfiles{}
}
