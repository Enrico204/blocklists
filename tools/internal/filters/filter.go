package filters

import (
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/url"
)

type Filter interface {
	// Fetch downloads the blocklist and returns the list of CIDRs.
	Fetch(logger *zap.SugaredLogger, client *http.Client, blocklistURL *url.URL) ([]*net.IPNet, error)
}
