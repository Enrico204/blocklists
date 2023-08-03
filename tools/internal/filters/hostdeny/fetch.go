package hostdeny

import (
	"bufio"
	"fmt"
	"git.netsplit.it/enrico204/blocklists/tools/pkg/utils"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func (a hostdeny) Fetch(logger *zap.SugaredLogger, client *http.Client, blocklistURL *url.URL) ([]*net.IPNet, error) {
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
		fields := strings.Split(line, " : ")
		if len(fields) != 2 {
			logger.Warnw("malformed line, skipping")
			continue
		}

		addr, err := utils.ParseCIDR(fields[1])
		if err != nil {
			logger.Warnw("skipping malformed line", "err", err)
			continue
		}

		ret = append(ret, addr)
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return utils.Aggregate(ret), nil
}
