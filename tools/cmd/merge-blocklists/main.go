package main

import (
	"errors"
	"flag"
	"fmt"
	"git.netsplit.it/enrico204/blocklists/tools/internal/blocklists"
	"git.netsplit.it/enrico204/blocklists/tools/internal/utils"
	"go.uber.org/zap"
	"net"
	"os"
	"path"
	"strings"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	var cfgfile = CommandLine.String("cfg", "cfg/firehol.yaml", "Configuration file.")
	var quietMode = CommandLine.Bool("quiet", false, "Print only warnings and errors to stderr.")
	var verboseMode = CommandLine.Bool("verbose", false, "Verbose log: prints more information about what the program is doing.")
	CommandLine.Usage = func() {
		_, _ = fmt.Fprint(CommandLine.Output(), "\nMerges multiple blocklists, aggregating prefixes.\n\n")
		_, _ = fmt.Fprintf(CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		CommandLine.PrintDefaults()
	}

	err := CommandLine.Parse(os.Args[1:])
	if err != nil && errors.Is(err, flag.ErrHelp) {
		return nil
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

	// Loading config.
	cfg, err := loadConfig(*cfgfile)
	if err != nil {
		logger.Errorw("error loading config", "err", err)
		return err
	}

	// Merge lists.
	for blocklist, details := range cfg.Lists {
		logger := logger.With("blocklist", blocklist)

		var aggregate []*net.IPNet
		logger.Info("loading and merging lists")
		for _, tag := range details.Include {
			blocklistPath := path.Join(cfg.ListsDirectory, tag+".list")
			partial, err := blocklists.ReadBlocklistFile(logger, blocklistPath)
			if err != nil {
				logger.Errorw("error reading blocklist file", "err", err, "blocklist-file", blocklistPath)
				return err
			}

			aggregate = append(aggregate, partial...)
		}

		aggregate = utils.Aggregate(aggregate)

		var out strings.Builder
		for _, ip := range aggregate {
			_, _ = out.WriteString(ip.String())
			_, _ = out.WriteRune('\n')
		}

		destinationPath := path.Join(cfg.ListsDirectory, blocklist+".list")
		err = os.WriteFile(destinationPath, []byte(out.String()), 0600)
		if err != nil {
			logger.Errorw("error saving blocklist", "err", err)
			return err
		}
	}
	return nil
}
