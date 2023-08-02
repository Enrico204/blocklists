package etpix

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

func (a emergingthreatsPix) Fetch(logger *zap.SugaredLogger, client *http.Client, blocklistURL *url.URL) ([]*net.IPNet, error) {
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
		fields := strings.Fields(line)
		if len(fields) == 2 {
			// access list creation rule, skipping
			continue
		} else if len(fields) != 7 {
			logger.Warnw("malformed line, skipping")
			continue
		}

		ipaddr := net.ParseIP(fields[4])
		subnetaddr := net.ParseIP(fields[5])
		if ipaddr == nil || subnetaddr == nil {
			logger.Warnw("malformed line, skipping")
			continue
		}
		addr := &net.IPNet{
			IP:   ipaddr,
			Mask: net.IPMask(subnetaddr),
		}
		ret = append(ret, addr)
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return utils.Aggregate(ret), err
}
