package gpfcomics

import (
	"fmt"
	"git.netsplit.it/enrico204/blocklists/tools/pkg/utils"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func (a gpfcomics) Fetch(logger *zap.SugaredLogger, client *http.Client, blocklistURL *url.URL) ([]*net.IPNet, error) {
	tmpfilename := fmt.Sprintf("%s/%s", a.cacheDir, utils.SHA256String(blocklistURL.String()))

	// **** IPV4 ****
	req, err := http.NewRequest(http.MethodPost, blocklistURL.String(), strings.NewReader("ipv6=0&export_type=text&submit=Export"))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	err = utils.GetFileIfModifiedSinceUsingRequest(logger, tmpfilename+".ipv4", client, req)
	if err != nil {
		return nil, err
	}

	fp, err := os.Open(tmpfilename + ".ipv4")
	if err != nil {
		return nil, err
	}
	defer func() { _ = fp.Close() }()

	v4networks, err := utils.ExtractNetworks(logger, fp)
	if err != nil {
		return nil, err
	}

	// **** IPV6 ****
	req, err = http.NewRequest(http.MethodPost, blocklistURL.String(), strings.NewReader("ipv6=1&export_type=text&submit=Export"))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	err = utils.GetFileIfModifiedSinceUsingRequest(logger, tmpfilename+".ipv6", client, req)
	if err != nil {
		return nil, err
	}

	fp6, err := os.Open(tmpfilename + ".ipv6")
	if err != nil {
		return nil, err
	}
	defer func() { _ = fp6.Close() }()

	v6networks, err := utils.ExtractNetworks(logger, fp6)
	if err != nil {
		return nil, err
	}

	return utils.Aggregate(append(v4networks, v6networks...)), nil
}
