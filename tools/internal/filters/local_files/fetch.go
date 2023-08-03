package localfiles

import (
	"git.netsplit.it/enrico204/blocklists/tools/pkg/utils"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/url"
	"os"
)

func (a localfiles) Fetch(logger *zap.SugaredLogger, _ *http.Client, blocklistURL *url.URL) ([]*net.IPNet, error) {
	if blocklistURL.Host == "." {
		// Local path
		blocklistURL.Path = "." + blocklistURL.Path
	}

	fp, err := os.Open(blocklistURL.Path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = fp.Close() }()

	return utils.ExtractNetworks(logger, fp)
}
