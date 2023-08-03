package localfiles

import (
	"git.netsplit.it/enrico204/blocklists/tools/internal/filters"
)

type localfiles struct {
}

func New() filters.Filter {
	return &localfiles{}
}
