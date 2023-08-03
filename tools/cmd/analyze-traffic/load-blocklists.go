package main

import (
	"git.netsplit.it/enrico204/blocklists/tools/pkg/blocklists"
	"github.com/yl2chen/cidranger"
	"path"
)

// Variables to hold matchers and data
var (
	aggregatedRanger       = cidranger.NewPCTrieRanger()
	perBlocklistRanger     = make(map[string]cidranger.Ranger)
	aggregatedSrcMatches   = uint64(0)
	aggregatedDstMatches   = uint64(0)
	perBlocklistSrcMatches = make(map[string]uint64)
	perBlocklistDstMatches = make(map[string]uint64)
)

// loadBlocklists loads the specified IP lists in memory and prepare the matchers.
func loadBlocklists() error {
	// Load blocklists in RAM
	logger.Debug("starting loading blocklists")
	for _, category := range blocklists.Categories {
		logger := logger.With("category", category)

		// Open the blocklist index.
		logger.Debug("loading blocklist index")
		blocklistIndex, err := blocklists.LoadBlocklistIndex(path.Join(*indexDir, category+".yaml"))
		if err != nil {
			logger.Warnw("can't load the blocklist index", "err", err)
			//return err
			continue
		}

		for tag, details := range blocklistIndex {
			logger := logger.With("blocklist", tag)
			if details.DisabledReason != "" {
				continue
			}

			logger.Debug("loading blocklist")
			blocklistPath := path.Join(*blocklistsDir, tag+".list")
			partial, err := blocklists.ReadBlocklistFile(logger, blocklistPath)
			if err != nil {
				logger.Warnw("error reading blocklist file", "err", err, "blocklist-file", blocklistPath)
				continue
			}

			if *aggregateLists {
				for _, netaddr := range partial {
					err = aggregatedRanger.Insert(cidranger.NewBasicRangerEntry(*netaddr))
					if err != nil {
						logger.Errorw("error adding CIDR to matching tree", "err", err, "blocklist-file", blocklistPath)
						return err
					}
				}
			} else {
				perBlocklistSrcMatches[tag] = 0
				perBlocklistDstMatches[tag] = 0
				perBlocklistRanger[tag] = cidranger.NewPCTrieRanger()
				for _, netaddr := range partial {
					err = perBlocklistRanger[tag].Insert(cidranger.NewBasicRangerEntry(*netaddr))
					if err != nil {
						logger.Errorw("error adding CIDR to matching tree", "err", err, "blocklist-file", blocklistPath)
						return err
					}
				}
			}
		}
	}
	return nil
}
