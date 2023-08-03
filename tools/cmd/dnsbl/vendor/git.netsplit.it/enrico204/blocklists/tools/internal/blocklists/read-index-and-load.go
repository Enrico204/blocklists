package blocklists

import (
	"git.netsplit.it/enrico204/blocklists/tools/internal/utils"
	"go.uber.org/zap"
	"net"
	"path"
)

// ReadIndexAndLoad returns the aggregated list of all IPs in lists specified in listsToLoad. If listsToLoad is empty,
// then all lists are loaded.
func ReadIndexAndLoad(logger *zap.SugaredLogger, yamlDirectory string, listsDirectory string, listsToLoad Names) []*net.IPNet {
	logger.Debug("Test")

	var aggregate []*net.IPNet
	for _, category := range Categories {
		logger := logger.With("category", category)

		// Open the blocklist index.
		logger.Debug("loading blocklist index")
		blocklistIndex, err := LoadBlocklistIndex(path.Join(yamlDirectory, category+".yaml"))
		if err != nil {
			logger.Warnw("can't load the blocklist index", "err", err)
			continue
		}

		for tag, details := range blocklistIndex {
			logger := logger.With("blocklist", tag)
			if details.DisabledReason != "" {
				continue
			} else if len(listsToLoad) > 0 && !listsToLoad.Contains(tag) {
				continue
			}

			logger.Debug("loading blocklist")
			blocklistPath := path.Join(listsDirectory, tag+".list")
			partial, err := ReadBlocklistFile(logger, blocklistPath)
			if err != nil {
				logger.Warnw("error reading blocklist file", "err", err, "blocklist-file", blocklistPath)
				continue
			}
			aggregate = append(aggregate, partial...)
		}
	}
	return utils.Aggregate(aggregate)
}
