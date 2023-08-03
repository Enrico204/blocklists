package main

import (
	"fmt"
	"git.netsplit.it/enrico204/blocklists/tools/internal/filters"
	"git.netsplit.it/enrico204/blocklists/tools/pkg/blocklists"
	"git.netsplit.it/enrico204/blocklists/tools/pkg/utils"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"os"
	"path"
)

// handleBlocklist checks if it's time to update the blocklist. If yes, the update is executed using the handler.
// If no error occurs, the final IP address list is saved in "outDir" using the "tag" as filename
func handleBlocklist(logger *zap.SugaredLogger, outDir string, hdl filters.Filter, tag string, bl blocklists.IndexItem, debugMode bool, tmpdir string, canBeEmpty bool) error {
	tempFileName := path.Join(tmpdir, tag+".list")
	destinationFileName := path.Join(outDir, tag+".list")

	// Check if we need to refresh the blocklist.
	timeToUpdate, noUpdateUntil, err := utils.IsTimeToUpdateYet(destinationFileName, bl.UpdateEvery)
	if err != nil {
		return err
	} else if !timeToUpdate {
		logger.Debugw("No need to update the blocklist for another " + noUpdateUntil.String())
		return nil
	}

	// Set some default options in the HTTP client.
	var httpc = &http.Client{
		Timeout: HTTPClientTimeout,
	}
	if debugMode {
		httpc.Transport = &utils.HTTPLoggingTransport{Logger: logger}
	}

	parsedURL, err := url.Parse(bl.URL)
	if err != nil {
		return fmt.Errorf("parsing blocklist url: %w", err)
	}

	// Create a temporary file where we will put our data. If the fetch is successful, we will move this file to the
	// final destination path.
	fp, err := os.Create(tempFileName)
	if err != nil {
		return fmt.Errorf("creating blocklist file: %w", err)
	}

	logger.Debug("fetching blocklist")
	addresses, err := hdl.Fetch(logger, httpc, parsedURL)
	if err != nil {
		_ = fp.Close()
		return fmt.Errorf("fetching blocklist: %w", err)
	}
	addresses = utils.Aggregate(addresses)

	for _, addr := range addresses {
		_, err = fp.WriteString(addr.String() + "\n")
		if err != nil {
			_ = fp.Close()
			return fmt.Errorf("saving blocklist: %w", err)
		}
	}

	// File has been flushed successfully on disk?
	err = fp.Close()
	if err != nil {
		return fmt.Errorf("saving blocklist: %w", err)
	}

	if !canBeEmpty {
		// Is file empty?
		stat, err := os.Stat(tempFileName)
		if err != nil {
			return fmt.Errorf("error checking file size: %w", err)
		} else if stat.Size() == 0 {
			return fmt.Errorf("file empty")
		}
	}

	// Move temporary file into the real destination.
	return utils.MoveFileKeepingLastModifiedTime(tempFileName, destinationFileName)
}
