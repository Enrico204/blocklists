package plaintext

import (
	"archive/zip"
	"fmt"
	"git.netsplit.it/enrico204/blocklists/tools/internal/utils"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/url"
	"os"
)

func (a plaintext) Fetch(logger *zap.SugaredLogger, client *http.Client, blocklistURL *url.URL) ([]*net.IPNet, error) {
	tmpfilename := fmt.Sprintf("%s/%s", a.cacheDir, utils.SHA256String(blocklistURL.String()))

	err := utils.GetFileIfModifiedSince(logger, tmpfilename, client, blocklistURL)
	if err != nil {
		return nil, err
	}

	// We have a ZIP file, so we open the ZIP and pass all files in the cleaning+copy function.
	if a.opts.Unzip {
		zipfp, err := zip.OpenReader(tmpfilename)
		if err != nil {
			return nil, err
		}
		defer func() { _ = zipfp.Close() }()

		var ret []*net.IPNet
		// For each file in the ZIP file, open the file and read the content.
		for _, fp := range zipfp.File {
			fp, err := fp.Open()
			if err != nil {
				return nil, err
			}

			var networks []*net.IPNet
			if a.opts.UseCustomSeparator {
				networks, err = utils.ExtractNetworksCustomNewLine(logger, fp, a.opts.CustomSeparator)
			} else {
				networks, err = utils.ExtractNetworks(logger, fp)
			}
			if err != nil {
				_ = fp.Close()
				return nil, err
			}

			_ = fp.Close()

			ret = append(ret, networks...)
		}
		// End of ZIP file.
		return utils.Aggregate(ret), nil
	}

	// No zip, open the file normally and copy the content to "out" after cleaning
	fp, err := os.Open(tmpfilename)
	if err != nil {
		return nil, err
	}
	defer func() { _ = fp.Close() }()

	if a.opts.UseCustomSeparator {
		return utils.ExtractNetworksCustomNewLine(logger, fp, a.opts.CustomSeparator)
	}
	return utils.ExtractNetworks(logger, fp)
}
