package apacheconf

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

func (a apacheconf) Fetch(logger *zap.SugaredLogger, client *http.Client, blocklistURL *url.URL) ([]*net.IPNet, error) {
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
		if !strings.HasPrefix(line, "deny from") {
			// Skip comments at the beginning of the line
			continue
		}

		fields := strings.Fields(line)
		if len(fields) != 3 {
			logger.Warnw("malformed line, skipping")
			continue
		}

		addr, err := utils.ParseCIDR(fields[2])
		if err != nil {
			logger.Warnw("skipping malformed line", "err", err)
			continue
		}

		ret = append(ret, addr)
	}
	if scanner.Err() != nil {
		return ret, scanner.Err()
	}

	return utils.Aggregate(ret), err
}
