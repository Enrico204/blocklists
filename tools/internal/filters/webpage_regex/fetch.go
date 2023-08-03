package webpagerx

import (
	"fmt"
	"git.netsplit.it/enrico204/blocklists/tools/internal/utils"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func (a webpageregex) Fetch(logger *zap.SugaredLogger, client *http.Client, blocklistURL *url.URL) ([]*net.IPNet, error) {
	tmpfilename := fmt.Sprintf("%s/%s", a.cacheDir, utils.SHA256String(blocklistURL.String()))

	err := utils.GetFileIfModifiedSince(logger, tmpfilename, client, blocklistURL)
	if err != nil {
		return nil, err
	}

	fcontentbytes, err := os.ReadFile(tmpfilename)
	if err != nil {
		return nil, err
	}

	submatches := a.matcher.FindAllStringSubmatch(string(fcontentbytes), -1)
	var ret = make([]*net.IPNet, 0, len(submatches))
	for idx, m := range submatches {
		if len(m) < 2 {
			continue
		} else if strings.TrimSpace(m[1]) == "" {
			continue
		}

		addr, err := utils.ParseCIDR(m[1])
		if err != nil {
			logger.Warnw("skipping malformed line", "line_n", idx, "err", err)
			continue
		}

		ret = append(ret, addr)
	}

	return utils.Aggregate(ret), nil
}
