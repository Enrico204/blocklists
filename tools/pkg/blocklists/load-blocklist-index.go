package blocklists

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// LoadBlocklistIndex returns the blocklist from its YAML.
func LoadBlocklistIndex(fname string) (Index, error) {
	fp, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("error opening blocklist index file: %w", err)
	}
	defer func() { _ = fp.Close() }()

	var blocklistIndex Index
	err = yaml.NewDecoder(fp).Decode(&blocklistIndex)
	if err != nil {
		return nil, fmt.Errorf("error parsing blocklist index file: %w", err)
	}
	return blocklistIndex, nil
}
