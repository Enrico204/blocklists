package dshield

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

func (a dshield) Fetch(logger *zap.SugaredLogger, client *http.Client, blocklistURL *url.URL) ([]*net.IPNet, error) {
	var ret []*net.IPNet
	tmpfilename := fmt.Sprintf("%s/%s", a.cacheDir, utils.SHA256String(blocklistURL.String()))

	err := utils.GetFileIfModifiedSince(logger, tmpfilename, client, blocklistURL)
	if err != nil {
		return ret, err
	}

	fp, err := os.Open(tmpfilename)
	if err != nil {
		return ret, err
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
		fields := strings.Fields(line)
		firstIP := net.ParseIP(fields[0])
		lastIP := net.ParseIP(fields[1])
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
