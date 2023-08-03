package dnsbl

import (
	"github.com/coredns/coredns/plugin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// requestCount exports a prometheus metric that is incremented every time a query is seen by the DNSBL plugin.
var requestCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: plugin.Namespace,
	Subsystem: "dnsbl",
	Name:      "request_count_total",
	Help:      "Counter of requests made.",
}, []string{"server"})

// matchesCount exports a prometheus metric that is incremented every time a match is found by the DNSBL plugin.
var matchesCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: plugin.Namespace,
	Subsystem: "dnsbl",
	Name:      "matches_count_total",
	Help:      "Counter of matches against IP lists.",
}, []string{"server"})
