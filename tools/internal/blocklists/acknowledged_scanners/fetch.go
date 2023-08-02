package ackscanners

import (
	"encoding/csv"
	"fmt"
	"git.netsplit.it/enrico204/blocklists/tools/internal/gitrepo"
	"git.netsplit.it/enrico204/blocklists/tools/internal/utils"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
)

func (a ackscanners) Fetch(logger *zap.SugaredLogger, _ *http.Client, blocklistURL *url.URL) ([]*net.IPNet, error) {
	var ret []*net.IPNet

	g := gitrepo.New(logger, blocklistURL.String(), a.cachedir)
	repoExists, err := g.Exists()
	if err != nil {
		return ret, fmt.Errorf("checking if repo exists: %w", err)
	}

	if repoExists {
		_, err = g.Pull()
		if err != nil {
			repoExists = false
			err = os.RemoveAll(a.cachedir)
			if err != nil {
				return ret, fmt.Errorf("removing old dir: %w", err)
			}
			err = os.MkdirAll(a.cachedir, 0750)
			if err != nil {
				return ret, fmt.Errorf("creating new dir: %w", err)
			}
		}
	}
	if !repoExists {
		_, err = g.Clone()
		if err != nil {
			return ret, fmt.Errorf("cloning repo: %w", err)
		}
	}

	fp, err := os.Open(path.Join(a.cachedir, "data", "orgs.csv"))
	if err != nil {
		return ret, fmt.Errorf("opening orgs.csv file: %w", err)
	}
	defer func() { _ = fp.Close() }()

	records, err := csv.NewReader(fp).ReadAll()
	if err != nil {
		return ret, fmt.Errorf("reading orgs.csv file: %w", err)
	}

	for idx, row := range records {
		if len(row) == 0 {
			continue
		} else if len(row) != 8 {
			return ret, fmt.Errorf("invalid row format in orgs.csv line %d", idx)
		}

		tag := row[0]
		inactive := row[3]
		if inactive != "" {
			continue
		}

		ipsfp, err := os.Open(path.Join(a.cachedir, "data", tag, "ips.txt"))
		if err != nil {
			logger.Warnw("error opening ips.txt for this tag", "tag", tag, "err", err)
			continue
		}
		networks, err := utils.ExtractNetworks(logger, ipsfp)
		if err != nil {
			logger.Warnw("error extracting IPs", "tag", tag, "err", err)
			continue
		}
		_ = ipsfp.Close()

		ret = append(ret, networks...)
	}

	return utils.Aggregate(ret), nil
}
