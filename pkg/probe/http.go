package probe

import (
	"context"
	"net/http"
	"time"
)

type HTTPTarget struct {
	URL            string `yaml:"url"`
	ExpectedStatus int    `yaml:"expected_status"`
}

type HTTPResult struct {
	URL        string
	StatusCode int
	Latency    time.Duration
	Success    bool
	Error      error
}

func ProbeHTTP(ctx context.Context, target HTTPTarget, timeout time.Duration) *HTTPResult {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target.URL, nil)
	if err != nil {
		return &HTTPResult{URL: target.URL, Success: false, Error: err}
	}

	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // don't follow redirects automatically
		},
	}

	start := time.Now()
	resp, err := client.Do(req)
	latency := time.Since(start)

	if err != nil {
		return &HTTPResult{
			URL:     target.URL,
			Latency: latency,
			Success: false,
			Error:   err,
		}
	}
	defer resp.Body.Close()

	expected := target.ExpectedStatus
	if expected == 0 {
		expected = http.StatusOK
	}

	return &HTTPResult{
		URL:        target.URL,
		StatusCode: resp.StatusCode,
		Latency:    latency,
		Success:    resp.StatusCode == expected,
	}
}
