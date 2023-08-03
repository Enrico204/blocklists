package ip9datacenters

import (
	"encoding/csv"
	"fmt"
	"git.netsplit.it/enrico204/blocklists/tools/internal/utils"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/url"
	"os"
)

func (a ip6datacenters) Fetch(logger *zap.SugaredLogger, client *http.Client, blocklistURL *url.URL) ([]*net.IPNet, error) {
	var ret []*net.IPNet
	tmpfilename := fmt.Sprintf("%s/%s", a.cacheDir, utils.SHA256String(blocklistURL.String()))

	err := utils.GetFileIfModifiedSince(logger, tmpfilename, client, blocklistURL)
	if err != nil {
		return nil, err
	}

	fp, err := os.Open(tmpfilename)
	if err != nil {
		return nil, err
	}
	defer func() { _ = fp.Close() }()

	scanner := csv.NewReader(fp)
	rows, err := scanner.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		if len(row) < 2 {
			logger.Warnw("malformed line, skipping")
			continue
		}

		firstIP := net.ParseIP(row[0])
		lastIP := net.ParseIP(row[1])
		if firstIP == nil || lastIP == nil {
			logger.Warnw("malformed line, skipping")
			continue
		}

		ret = append(ret, utils.IPrangeToCIDR(firstIP, lastIP)...)
	}
	return utils.Aggregate(ret), nil
}
