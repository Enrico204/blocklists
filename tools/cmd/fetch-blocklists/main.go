package main

import (
	"errors"
	"flag"
	"fmt"
	"git.netsplit.it/enrico204/blocklists/tools/internal/blocklists"
	"go.uber.org/zap"
	"os"
	"path"
	"time"
)

const HTTPClientTimeout = 120 * time.Second

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	var blocklistDir = CommandLine.String("lists-directory", ".", "Blocklists YAML files location.")
	var outDir = CommandLine.String("output-directory", "out/", `Output directory.
Cleaned, filtered and aggregated lists will be saved here.
The directory will be created if it does not exist.`)
	var cacheDir = CommandLine.String("cache-directory", "tmp/", `Cache directory.
The program will store downloaded but unfiltered lists in this directory.
It's used to avoid re-downloading the same content again and again.
The directory will be created if it does not exist.`)
	var quietMode = CommandLine.Bool("quiet", false, "Print only warnings and errors to stderr.")
	var verboseMode = CommandLine.Bool("verbose", false, "Verbose log: prints more information about what the program is doing.")
	var debugMode = CommandLine.Bool("debug", false, "Debug mode: verbose log + dumps HTTP headers in the log.")
	CommandLine.Usage = func() {
		_, _ = fmt.Fprint(CommandLine.Output(), "\nFetches all IP lists configured in YAML files and aggregate them.\n\n")
		_, _ = fmt.Fprintf(CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		CommandLine.PrintDefaults()
	}

	err := CommandLine.Parse(os.Args[1:])
	if err != nil && errors.Is(err, flag.ErrHelp) {
		return nil
	}

	if *debugMode {
		*verboseMode = true
	}

	// Setting up the logger.
	var zlogger *zap.Logger
	var zcfg zap.Config
	if *quietMode {
		zcfg = zap.Config{
			Level:            zap.NewAtomicLevelAt(zap.WarnLevel),
			Encoding:         "console",
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}
	} else if *verboseMode {
		zcfg = zap.NewDevelopmentConfig()
	} else {
		zcfg = zap.Config{
			Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
			Development:      true,
			Encoding:         "console",
			EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}
	}
	zlogger, _ = zcfg.Build()
	defer func() { _ = zlogger.Sync() }()
	logger := zlogger.Sugar()

	// Creating cache and output directories, if needed.
	err = os.MkdirAll(*cacheDir, 0750)
	if err != nil {
		logger.Errorw("can't create the cache directory", "err", err)
		return err
	}

	err = os.MkdirAll(*outDir, 0750)
	if err != nil {
		logger.Errorw("can't create the output directory", "err", err)
		return err
	}

	// initialize all blocklists filters.
	var filters = registerFilters(logger, *cacheDir)

	// For each category, open the blocklist file, and for each blocklist in the file, check if the handler exists,
	// whether the blocklist is disabled, and then fetch it if it's time.
	for _, category := range blocklists.BlocklistCategories {
		logger := logger.With("category", category)

		// Open the blocklist index.
		blocklistIndex, err := loadBlocklistIndex(logger, path.Join(*blocklistDir, category+".yaml"))
		if err != nil {
			logger.Errorw("can't load the blocklist index", "err", err)
			return err
		}

		logger.Info("fetching and filtering the lists")
		for tag, bl := range blocklistIndex {
			logger := logger.With("blocklist", tag)

			// Is the blocklist disabled for any reason?
			if bl.DisabledReason != "" {
				logger.Debugw("blocklist disabled, skipping")
				continue
			}

			// Get the filter for the blocklist, if it exists. If not, skip and log.
			hdl, ok := filters[bl.Filter]
			if !ok {
				logger.Warnw("filter not found, skipping list", "filter", bl.Filter)
				continue
			}

			// Work on this blocklist, and refresh it if needed.
			err = handleBlocklist(logger, *outDir, hdl, tag, bl, *debugMode, *cacheDir, bl.CanBeEmpty)
			if err != nil {
				logger.Warnw("error while working on blocklist", "err", err)
			}
		}
	}
	return nil
}
