package probe

import (
	"context"
	"net"
	"time"
)

type DNSTarget struct {
	Host   string `yaml:"host"`   // e.g. "api.example.com"
	Server string `yaml:"server"` // e.g. "8.8.8.8:53"
}

type DNSResult struct {
	Host    string
	Latency time.Duration
	Success bool
	IPs     []string
	Error   error
}

func ProbeDNS(ctx context.Context, target DNSTarget, timeout time.Duration) *DNSResult {
	resolver := net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: timeout}
			server := target.Server
			if server != "" && !hasPort(server) {
				server = net.JoinHostPort(server, "53")
			}
			if server == "" {
				server = "8.8.8.8:53"
			}
			return d.DialContext(ctx, "udp", server)
		},
	}

	start := time.Now()
	ips, err := resolver.LookupHost(ctx, target.Host)
	latency := time.Since(start)

	if err != nil {
		return &DNSResult{Host: target.Host, Latency: latency, Success: false, Error: err}
	}

	return &DNSResult{
		Host:    target.Host,
		Latency: latency,
		Success: len(ips) > 0,
		IPs:     ips,
	}
}

func hasPort(s string) bool {
	_, _, err := net.SplitHostPort(s)
	return err == nil
}
