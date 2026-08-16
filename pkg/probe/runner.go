package probe

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/txltedxgod/healthgate/pkg/metrics"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Interval    time.Duration `yaml:"interval"`
	HTTPTargets []HTTPTarget  `yaml:"http_targets"`
	TCPTargets  []TCPTarget   `yaml:"tcp_targets"`
	DNSTargets  []DNSTarget   `yaml:"dns_targets"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Interval == 0 {
		cfg.Interval = 10 * time.Second
	}
	return &cfg, nil
}

type Runner struct {
	config *Config
}

func NewRunner(cfg *Config) *Runner {
	return &Runner{config: cfg}
}

func (r *Runner) Start() {
	r.runAll()
	ticker := time.NewTicker(r.config.Interval)
	for range ticker.C {
		r.runAll()
	}
}

func (r *Runner) runAll() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// HTTP
	for _, target := range r.config.HTTPTargets {
		res := ProbeHTTP(ctx, target, 4*time.Second)
		successVal := 0.0
		if res.Success {
			successVal = 1.0
		}
		metrics.ProbeSuccess.WithLabelValues("http", target.URL).Set(successVal)
		metrics.ProbeLatencySeconds.WithLabelValues("http", target.URL).Observe(res.Latency.Seconds())
	}

	// TCP
	for _, target := range r.config.TCPTargets {
		res := ProbeTCP(target, 4*time.Second)
		successVal := 0.0
		if res.Success {
			successVal = 1.0
		}
		metrics.ProbeSuccess.WithLabelValues("tcp", target.Address).Set(successVal)
		metrics.ProbeLatencySeconds.WithLabelValues("tcp", target.Address).Observe(res.Latency.Seconds())
	}

	// DNS
	for _, target := range r.config.DNSTargets {
		res := ProbeDNS(ctx, target, 4*time.Second)
		successVal := 0.0
		if res.Success {
			successVal = 1.0
		}
		metrics.ProbeSuccess.WithLabelValues("dns", target.Host).Set(successVal)
		metrics.ProbeLatencySeconds.WithLabelValues("dns", target.Host).Observe(res.Latency.Seconds())
	}
	log.Println("[healthgate] Probed all targets successfully")
}
