package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/txltedxgod/healthgate/pkg/probe"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to probe config yaml")
	listenAddr := flag.String("listen", ":9115", "Metrics server address")
	flag.Parse()

	log.Printf("[healthgate] Starting network probe exporter on %s\n", *listenAddr)

	cfg, err := probe.LoadConfig(*configPath)
	if err != nil {
		log.Printf("Warning: failed to load config %s, using defaults: %v\n", *configPath, err)
		cfg = &probe.Config{
			Interval: 15 * time.Second,
			HTTPTargets: []probe.HTTPTarget{
				{URL: "https://httpbin.org/status/200", ExpectedStatus: 200},
			},
			TCPTargets: []probe.TCPTarget{
				{Address: "1.1.1.1:53"},
			},
			DNSTargets: []probe.DNSTarget{
				{Host: "google.com", Server: "8.8.8.8"},
			},
		}
	}

	runner := probe.NewRunner(cfg)
	go runner.Start()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK\n"))
	})

	server := &http.Server{
		Addr:         *listenAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("[healthgate] Shutting down probe exporter...")
}
