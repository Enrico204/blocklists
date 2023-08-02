package p2p

import (
	"bufio"
	"fmt"
	"git.netsplit.it/enrico204/blocklists/tools/internal/utils"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func (a p2p) Fetch(logger *zap.SugaredLogger, client *http.Client, blocklistURL *url.URL) ([]*net.IPNet, error) {
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

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			// Skip empty lines
			continue
		} else if line[0] == '#' {
			// Skip comments at the beginning of the line
			continue
		}

		fields := strings.Split(line, ":")
		if len(fields) != 2 {
			logger.Warnw("malformed line, skipping")
			continue
		}
		addresses := strings.Split(fields[1], "-")
		if len(addresses) != 2 {
			logger.Warnw("malformed line, skipping")
			continue
		}

		firstIP := net.ParseIP(addresses[0])
		lastIP := net.ParseIP(addresses[1])
		if firstIP == nil || lastIP == nil {
			logger.Warnw("malformed line, skipping")
			continue
		}

		ret = append(ret, utils.IPrangeToCIDR(firstIP, lastIP)...)
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return utils.Aggregate(ret), nil
}
