package utils

import (
	"bufio"
	"go.uber.org/zap"
	"io"
	"net"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	hashCommentRx      = regexp.MustCompile(`#.*$`)
	semicolonCommentRx = regexp.MustCompile(`;.*$`)
)

// ExtractNetworks extracts all CIDRs/IPs from `src` in aggregated form (see Aggregate).
// The `src` should be an IP or CIDR list. Empty lines are ignored. Comments can start with `#` or `;`, and they are
// ignored as well. The newline should be \n and/or \r.
func ExtractNetworks(logger *zap.SugaredLogger, src io.Reader) ([]*net.IPNet, error) {
	return extractNetworksCustomFunction(logger, src, bufio.ScanLines)
}

// ExtractNetworksCustomNewLine extracts all CIDRs/IPs from `src` in aggregated form (see Aggregate).
// The `src` should be an IP or CIDR list. Empty lines are ignored. Comments can start with `#` or `;`, and they are
// ignored as well. It uses the `newlineSeparator` as character for line separator.
func ExtractNetworksCustomNewLine(logger *zap.SugaredLogger, src io.Reader, newlineSeparator rune) ([]*net.IPNet, error) {
	return extractNetworksCustomFunction(logger, src, generateScanForNewline(newlineSeparator))
}

func extractNetworksCustomFunction(logger *zap.SugaredLogger, src io.Reader, splitFn bufio.SplitFunc) ([]*net.IPNet, error) {
	scanner := bufio.NewScanner(src)

	var ret []*net.IPNet
	scanner.Split(splitFn)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			// Skip empty lines
			continue
		}

		if line[0] == '#' || line[0] == ';' {
			// Skip comments at the beginning of the line
			continue
		}

		if strings.ContainsRune(line, '#') {
			line = hashCommentRx.ReplaceAllString(line, "")
		} else if strings.ContainsRune(line, ';') {
			line = semicolonCommentRx.ReplaceAllString(line, "")
		}

		line = strings.ToLower(strings.TrimSpace(line))
		if len(line) == 0 {
			// Skip empty lines
			continue
		}

		addr, err := ParseCIDR(line)
		if err != nil {
			logger.Warnw("malformed line, skipping", "err", err)
			continue
		}

		ret = append(ret, addr)
	}

	err := scanner.Err()
	if err != nil {
		return nil, err
	}

	return Aggregate(ret), err
}

func generateScanForNewline(newline rune) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		start := 0
		var width int
		// Scan until space, marking end of word.
		for i := 0; i < len(data); i += width {
			var r rune
			r, width = utf8.DecodeRune(data[i:])
			if r == newline {
				return i + width, data[start:i], nil
			}
		}
		// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
		if atEOF && len(data) > start {
			return len(data), data[start:], nil
		}
		// Request more data.
		return start, nil, nil
	}
}
