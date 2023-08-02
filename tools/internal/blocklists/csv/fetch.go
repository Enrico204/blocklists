package csv

import (
	csvpkg "encoding/csv"
	"fmt"
	"git.netsplit.it/enrico204/blocklists/tools/internal/utils"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/url"
	"os"
)

func (a csv) Fetch(logger *zap.SugaredLogger, client *http.Client, blocklistURL *url.URL) ([]*net.IPNet, error) {
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

	scanner := csvpkg.NewReader(fp)
	scanner.Comment = a.opts.Comment
	scanner.Comma = a.opts.Comma
	rows, err := scanner.ReadAll()
	if err != nil {
		return nil, err
	}

	var ret = make([]*net.IPNet, 0, len(rows))
	for idx, row := range rows {
		if idx == 0 && a.opts.SkipHeader {
			continue
		}

		if len(row) < a.opts.Column {
			logger.Warnw("malformed line, skipping")
			continue
		}

		addr, err := utils.ParseCIDR(row[a.opts.Column])
		if err != nil {
			logger.Warnw("skipping malformed line", "err", err)
			continue
		}

		ret = append(ret, addr)
	}
	return utils.Aggregate(ret), nil
}
