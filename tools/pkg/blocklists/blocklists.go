package blocklists

import "time"

var Categories = []string{"abuse", "anonymizers", "attacks", "malware", "organizations", "reputation", "spam", "unroutable", "scanners"}

type Names []string

func (n Names) Contains(name string) bool {
	for _, item := range n {
		if item == name {
			return true
		}
	}
	return false
}

type Index map[string]IndexItem

type IndexItem struct {
	// Filter is the "handler" that implements the logic for extracting data from the blocklist.
	Filter string `yaml:"filter"`

	// UpdateEvery indicates the update interval for the blocklist. It won't be updated in less than this interval.
	UpdateEvery time.Duration `yaml:"update_every"`

	// URL is the blocklist URL.
	URL string `yaml:"url"`

	// Info is a general description of this blocklist.
	Info string `yaml:"info,omitempty"`

	// Maintainer is the name of the maintainer.
	Maintainer string `yaml:"maintainer,omitempty"`

	// MaintainerURL is the URL for the maintainer of this blocklist.
	MaintainerURL string `yaml:"maintainer_url,omitempty"`

	// DisabledReason is the reason why this blocklist has been disabled. If empty, the blocklist is enabled.
	DisabledReason string `yaml:"disabled_reason,omitempty"`

	// CanBeEmpty allow the system to accept empty blocklists.
	CanBeEmpty bool `yaml:"can_be_empty,omitempty"`
}
