package main

import (
	"fmt"
	"git.netsplit.it/enrico204/blocklists/tools/internal/blocklists"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
)

// loadBlocklistIndex returns the blocklist from its YAML.
func loadBlocklistIndex(logger *zap.SugaredLogger, fname string) (blocklists.BlocklistIndex, error) {
	logger.Info("loading blocklist index")

	fp, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("error opening blocklist index file: %w", err)
	}
	defer func() { _ = fp.Close() }()

	var blocklistIndex blocklists.BlocklistIndex
	err = yaml.NewDecoder(fp).Decode(&blocklistIndex)
	if err != nil {
		return nil, fmt.Errorf("error parsing blocklist index file: %w", err)
	}
	return blocklistIndex, nil
}
