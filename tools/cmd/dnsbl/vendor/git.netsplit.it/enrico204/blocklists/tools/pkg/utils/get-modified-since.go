package utils

import (
	"compress/gzip"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// GetFileIfModifiedSince retrieves the file adding "If-Modified-Since" header (using the file last modification time)
// if the file exists. It also decompresses gzip content on the fly.
func GetFileIfModifiedSince(logger *zap.SugaredLogger, localFileName string, client *http.Client, remoteURL *url.URL) error {
	req, err := http.NewRequest(http.MethodGet, remoteURL.String(), nil)
	if err != nil {
		return err
	}
	return GetFileIfModifiedSinceUsingRequest(logger, localFileName, client, req)
}

// GetFileIfModifiedSinceUsingRequest retrieves the file adding "If-Modified-Since" header (using the file last
// modification time) if the file exists. It also decompresses gzip content on the fly.
// This variant allows the use of custom request. Note that "If-Modified-Since" is overwritten with the value coming
// from the mtime of the localFileName.
func GetFileIfModifiedSinceUsingRequest(logger *zap.SugaredLogger, localFileName string, client *http.Client, req *http.Request) error {
	var lastModified time.Time

	stat, err := os.Stat(localFileName)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	} else if err == nil {
		// File exists, retrieve the last modified
		lastModified = stat.ModTime()
	} else {
		logger.Debugw("local cache file does not exist, forcing re-download", "local-file-name", localFileName)
	}

	// TODO: make this part configurable
	if req.URL.Host == "talosintel.com" {
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
	}

	if !lastModified.IsZero() {
		req.Header.Set("If-Modified-Since", lastModified.UTC().Format(http.TimeFormat))
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotModified {
		logger.Debugw("not modified, skipping re-download")
		// Not modified since last download, end here
		return nil
	} else if resp.StatusCode == 200 {
		var lastModified = time.Now()
		if resp.Header.Get("Last-Modified") != "" {
			lastModified, err = http.ParseTime(resp.Header.Get("Last-Modified"))
			if err != nil {
				return err
			}
		}

		if resp.Header.Get("Content-Type") == "application/x-gzip" {
			resp.Body, err = gzip.NewReader(resp.Body)
			if err != nil {
				return err
			}
		}

		dest, err := os.Create(localFileName)
		if err != nil {
			return err
		}

		_, err = io.Copy(dest, resp.Body)
		if err != nil {
			_ = dest.Close()
			return err
		}
		err = dest.Close()
		if err != nil {
			return err
		}

		return os.Chtimes(localFileName, lastModified, lastModified)
	}
	return fmt.Errorf("unexpected status code %d", resp.StatusCode)
}
