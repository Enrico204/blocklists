// The goal of this program is to analyze a live traffic or a pcap replay to see if the traffic matches an IP address in
// the blocklists.
package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"go.uber.org/zap"
	"os"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

var (
	CommandLine    = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	aggregateLists = CommandLine.Bool("aggregate", false, "Aggregate all lists together. In this case, the result will be only match/no match.")
	indexDir       = CommandLine.String("index-directory", ".", "Blocklists YAML index files location.")
	blocklistsDir  = CommandLine.String("lists-directory", "out/", "Blocklists files location.")
	iface          = CommandLine.String("iface", "", "Interface to listen to (required, or specify -pcap).")
	promisc        = CommandLine.Bool("promisc", true, "Whether the interface should be in promiscuous mode (only when live capturing).")
	captureSize    = CommandLine.Int("size", 1600, "Capture size (only when live capturing).")
	pcapFile       = CommandLine.String("pcap", "", "PCAP file to analyze (required, or specify -iface).")
	quietMode      = CommandLine.Bool("quiet", false, "Print only warnings and errors to stderr.")
	verboseMode    = CommandLine.Bool("verbose", false, "Verbose log: prints more information about what the program is doing.")

	logger *zap.SugaredLogger
)

func run() error {
	// Parse command line.
	CommandLine.Usage = func() {
		_, _ = fmt.Fprint(CommandLine.Output(), "\nAnalyze a PCAP or a live traffic and reports who many IP packets matches which IP list.\n\n")
		_, _ = fmt.Fprintf(CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		CommandLine.PrintDefaults()
	}

	err := CommandLine.Parse(os.Args[1:])
	if err != nil && errors.Is(err, flag.ErrHelp) {
		return nil
	}

	// Validate command line parameters.
	if (*iface == "" && *pcapFile == "") || (*iface != "" && *pcapFile != "") {
		_, _ = fmt.Fprintln(CommandLine.Output(), "-pcap or -iface must be specified (but not both at the same time).")
		os.Exit(2)
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
	logger = zlogger.Sugar()

	// Load IP lists.
	if err = loadBlocklists(); err != nil {
		return err
	}

	// Open a file or a live interface.
	var handle *pcap.Handle
	if *pcapFile != "" {
		handle, err = pcap.OpenOffline(*pcapFile)
		if err != nil {
			logger.Errorw("error opening pcap file", "err", err)
			return err
		}
	} else if *iface != "" {
		handle, err = pcap.OpenLive(*iface, int32(*captureSize), *promisc, pcap.BlockForever)
		if err != nil {
			logger.Errorw("error opening live interface for capture", "err", err)
			return err
		}
	} else {
		panic("this should never happen")
	}
	defer handle.Close()

	// Decode incoming packets.
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		if packet.NetworkLayer() != nil {
			if err = processPacket(packet); err != nil {
				return err
			}
		}
	}

	// Print results on screen.
	if *aggregateLists {
		fmt.Printf("Aggregate matches for SRC: %d - DST: %d - TOTAL: %d\n", aggregatedSrcMatches, aggregatedDstMatches, aggregatedSrcMatches+aggregatedDstMatches)
	} else {
		fmt.Printf("Per blocklist results:\n")
		for tag := range perBlocklistSrcMatches {
			fmt.Printf("%30s - SRC: %9d - DST: %9d - TOTAL: %9d\n", tag, perBlocklistSrcMatches[tag], perBlocklistDstMatches[tag], perBlocklistSrcMatches[tag]+perBlocklistDstMatches[tag])
		}
	}

	return nil
}
