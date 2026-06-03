package proxy

import (
	"errors"
	"net/http"
	"time"

	"ripple/simulation"
)

// HandleProxy forwards request to target and injects network profile characteristics
func HandleProxy(w http.ResponseWriter, r *http.Request, profileName string, targetURL string) (*http.Response, time.Duration, error) {
	profile, exists := simulation.Profiles[profileName]
	if !exists {
		profile = simulation.Profiles["5G"] // Default
	}

	// 1. Simulate Latency
	if profile.Latency > 0 {
		time.Sleep(profile.Latency)
	}

	// 2. Simulate Packet Loss
	if profile.ShouldDrop() {
		return nil, 0, errors.New("network packet dropped (simulated packet loss)")
	}

	// 3. Create proxy request
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		return nil, 0, err
	}

	// Copy headers
	for k, vv := range r.Header {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	startTime := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		return nil, duration, err
	}

	return resp, duration, nil
}
