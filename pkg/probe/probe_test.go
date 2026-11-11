package probe

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfig_DefaultsAndMissing(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	yamlData := "http_targets:\n  - url: http://localhost:8080/health\n    expected_status: 200\n"
	if err := os.WriteFile(cfgPath, []byte(yamlData), 0o600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("unexpected LoadConfig error: %v", err)
	}
	if cfg.Interval != 10*time.Second {
		t.Fatalf("expected default interval 10s, got %v", cfg.Interval)
	}
	if len(cfg.HTTPTargets) != 1 || cfg.HTTPTargets[0].ExpectedStatus != 200 {
		t.Fatalf("unexpected http targets: %+v", cfg.HTTPTargets)
	}

	if _, err := LoadConfig(filepath.Join(dir, "missing.yaml")); err == nil {
		t.Fatal("expected error for missing config file")
	}
}

func TestProbeHTTP_LocalServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/created" {
			w.WriteHeader(http.StatusCreated)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ctx := context.Background()

	okRes := ProbeHTTP(ctx, HTTPTarget{URL: srv.URL + "/ok"}, 2*time.Second)
	if !okRes.Success || okRes.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 success, got %+v", okRes)
	}

	createdRes := ProbeHTTP(ctx, HTTPTarget{URL: srv.URL + "/created", ExpectedStatus: http.StatusCreated}, 2*time.Second)
	if !createdRes.Success || createdRes.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 success, got %+v", createdRes)
	}

	mismatchRes := ProbeHTTP(ctx, HTTPTarget{URL: srv.URL + "/created", ExpectedStatus: http.StatusOK}, 2*time.Second)
	if mismatchRes.Success {
		t.Fatalf("expected status mismatch failure, got %+v", mismatchRes)
	}
}

func TestProbeTCP_LocalListenerAndRunner(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to open local listener: %v", err)
	}
	addr := ln.Addr().String()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			conn, acceptErr := ln.Accept()
			if acceptErr != nil {
				return
			}
			_ = conn.Close()
		}
	}()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	tcpRes := ProbeTCP(TCPTarget{Address: addr}, time.Second)
	if !tcpRes.Success || tcpRes.Error != nil {
		t.Fatalf("expected local TCP probe success, got %+v", tcpRes)
	}

	runner := NewRunner(&Config{
		Interval:    time.Second,
		HTTPTargets: []HTTPTarget{{URL: srv.URL, ExpectedStatus: http.StatusOK}},
		TCPTargets:  []TCPTarget{{Address: addr}},
	})
	runner.runAll()

	_ = ln.Close()
	<-done

	closedRes := ProbeTCP(TCPTarget{Address: addr}, 200*time.Millisecond)
	if closedRes.Success || closedRes.Error == nil {
		t.Fatalf("expected closed TCP target to fail, got %+v", closedRes)
	}
}
