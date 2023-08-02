package main

import (
	"bufio"
	"go.uber.org/zap"
	"net"
	"os"
)

func readBlocklistFile(logger *zap.SugaredLogger, fname string) ([]*net.IPNet, error) {
	fp, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer func() { _ = fp.Close() }()

	var ret []*net.IPNet
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			// Skip empty lines
			continue
		}

		_, ipnet, err := net.ParseCIDR(line)
		if err != nil {
			logger.Warnw("no IP address detected, skipping line", "line", line, "err", err, "file", fname)
			continue
		}

		ret = append(ret, ipnet)
	}

	return ret, scanner.Err()
}
