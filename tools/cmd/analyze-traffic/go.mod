module git.netsplit.it/enrico204/blocklists/tools/cmd/analyze-traffic

go 1.20

replace git.netsplit.it/enrico204/blocklists/tools => ../../

require (
	git.netsplit.it/enrico204/blocklists/tools v0.0.0-00010101000000-000000000000
	github.com/google/gopacket v1.1.19
	github.com/yl2chen/cidranger v1.0.2
	go.uber.org/zap v1.25.0
)

require (
	github.com/ldkingvivi/go-aggregate v0.0.0-20200406164845-67d85734711c // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
