package probe

import (
	"net"
	"time"
)

type TCPTarget struct {
	Address string `yaml:"address"` // e.g. "127.0.0.1:5432"
}

type TCPResult struct {
	Address string
	Latency time.Duration
	Success bool
	Error   error
}

func ProbeTCP(target TCPTarget, timeout time.Duration) *TCPResult {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", target.Address, timeout)
	latency := time.Since(start)

	if err != nil {
		return &TCPResult{
			Address: target.Address,
			Latency: latency,
			Success: false,
			Error:   err,
		}
	}
	defer conn.Close()

	return &TCPResult{
		Address: target.Address,
		Latency: latency,
		Success: true,
	}
}
