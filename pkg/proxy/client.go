/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"

	"github.com/your-username/shadow-client/types"
)

// NewProxyManager creates and initializes a ProxyManager.
func NewProxyManager(toxiproxyURL string) *ProxyManager {
	return &ProxyManager{
		ToxiproxyURL: toxiproxyURL,
		ActiveToxics: make(map[string][]*ToxicConfig),
	}
}

// InitializeProxies sets up Toxiproxy proxies for Target A and Target B.
// It dynamically allocates free ports and creates proxies that forward traffic.
func (pm *ProxyManager) InitializeProxies(config *types.ShadowConfig) error {
	// Parse target URLs
	urlA, err := url.Parse(config.TargetA)
	if err != nil {
		return fmt.Errorf("invalid target-a URL: %w", err)
	}

	urlB, err := url.Parse(config.TargetB)
	if err != nil {
		return fmt.Errorf("invalid target-b URL: %w", err)
	}

	// Find free ports for proxies
	portA, err := findFreePort()
	if err != nil {
		return fmt.Errorf("could not find free port for target-a proxy: %w", err)
	}

	portB, err := findFreePort()
	if err != nil {
		return fmt.Errorf("could not find free port for target-b proxy: %w", err)
	}

	// Create proxy targets
	pm.TargetA = &ProxyTarget{
		Name:           "target-a",
		OriginalURL:    config.TargetA,
		ParsedURL:      urlA,
		ListenSocket:   fmt.Sprintf("127.0.0.1:%d", portA),
		UpstreamSocket: net.JoinHostPort(urlA.Hostname(), getPort(urlA)),
		ProxyName:      "shadow-client-target-a",
	}

	pm.TargetB = &ProxyTarget{
		Name:           "target-b",
		OriginalURL:    config.TargetB,
		ParsedURL:      urlB,
		ListenSocket:   fmt.Sprintf("127.0.0.1:%d", portB),
		UpstreamSocket: net.JoinHostPort(urlB.Hostname(), getPort(urlB)),
		ProxyName:      "shadow-client-target-b",
	}

	// Create proxies in Toxiproxy
	if err := pm.createProxy(pm.TargetA); err != nil {
		return fmt.Errorf("failed to create proxy for target-a: %w", err)
	}

	if err := pm.createProxy(pm.TargetB); err != nil {
		// Cleanup target-a proxy on failure
		pm.deleteProxy(pm.TargetA)
		return fmt.Errorf("failed to create proxy for target-b: %w", err)
	}

	// Update proxy URLs
	scheme := "http"
	if urlA.Scheme == "https" {
		scheme = "http" // Toxiproxy always exposes as http
	}
	pm.TargetA.ProxyURL = fmt.Sprintf("%s://%s%s", scheme, pm.TargetA.ListenSocket, urlA.Path)
	pm.TargetB.ProxyURL = fmt.Sprintf("%s://%s%s", scheme, pm.TargetB.ListenSocket, urlB.Path)

	slog.Info("Proxies initialized",
		"target_a_listen", pm.TargetA.ListenSocket,
		"target_a_upstream", pm.TargetA.UpstreamSocket,
		"target_b_listen", pm.TargetB.ListenSocket,
		"target_b_upstream", pm.TargetB.UpstreamSocket,
	)

	return nil
}

// ApplyNetworkProfile applies network toxics based on the selected profile.
func (pm *ProxyManager) ApplyNetworkProfile(profile types.NetworkProfile) error {
	slog.Debug("Applying network profile", "profile", profile.Name)

	// Apply latency and jitter as a single toxic
	if profile.Latency > 0 || profile.Jitter > 0 {
		latencyToxic := &ToxicConfig{
			Name:     "latency",
			Type:     "latency",
			Stream:   "upstream",
			Toxicity: 1.0,
			Attributes: map[string]interface{}{
				"latency": int(profile.Latency.Milliseconds()),
				"jitter":  int(profile.Jitter.Milliseconds()),
			},
		}

		if err := pm.addToxic(pm.TargetA.ProxyName, latencyToxic); err != nil {
			slog.Error("Failed to apply latency toxic to target-a", "error", err)
			return err
		}

		if err := pm.addToxic(pm.TargetB.ProxyName, latencyToxic); err != nil {
			slog.Error("Failed to apply latency toxic to target-b", "error", err)
			return err
		}

		slog.Debug("Applied latency toxic",
			"latency_ms", profile.Latency.Milliseconds(),
			"jitter_ms", profile.Jitter.Milliseconds(),
		)
	}

	// Apply bandwidth limit
	if profile.Bandwidth > 0 {
		bandwidthToxic := &ToxicConfig{
			Name:     "bandwidth_limit",
			Type:     "bandwidth_limit",
			Stream:   "upstream",
			Toxicity: 1.0,
			Attributes: map[string]interface{}{
				"rate": profile.Bandwidth,
			},
		}

		if err := pm.addToxic(pm.TargetA.ProxyName, bandwidthToxic); err != nil {
			slog.Error("Failed to apply bandwidth toxic to target-a", "error", err)
			return err
		}

		if err := pm.addToxic(pm.TargetB.ProxyName, bandwidthToxic); err != nil {
			slog.Error("Failed to apply bandwidth toxic to target-b", "error", err)
			return err
		}

		slog.Debug("Applied bandwidth limit toxic", "rate_bps", profile.Bandwidth)
	}

	// Apply packet loss (if non-zero)
	if profile.PacketLoss > 0 {
		packetLossToxic := &ToxicConfig{
			Name:     "packet_loss",
			Type:     "packet_loss",
			Stream:   "upstream",
			Toxicity: 1.0,
			Attributes: map[string]interface{}{
				"percentage": profile.PacketLoss,
			},
		}

		if err := pm.addToxic(pm.TargetA.ProxyName, packetLossToxic); err != nil {
			slog.Error("Failed to apply packet loss toxic to target-a", "error", err)
			return err
		}

		if err := pm.addToxic(pm.TargetB.ProxyName, packetLossToxic); err != nil {
			slog.Error("Failed to apply packet loss toxic to target-b", "error", err)
			return err
		}

		slog.Debug("Applied packet loss toxic", "percentage", profile.PacketLoss)
	}

	return nil
}

// RemoveAllToxics removes all toxics from both proxies.
func (pm *ProxyManager) RemoveAllToxics() error {
	var errs []error

	if pm.TargetA != nil {
		if err := pm.removeToxicsFromProxy(pm.TargetA.ProxyName); err != nil {
			slog.Error("Failed to remove toxics from target-a", "error", err)
			errs = append(errs, err)
		}
	}

	if pm.TargetB != nil {
		if err := pm.removeToxicsFromProxy(pm.TargetB.ProxyName); err != nil {
			slog.Error("Failed to remove toxics from target-b", "error", err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to remove some toxics: %v", errs)
	}

	return nil
}

// Cleanup removes both proxies from Toxiproxy.
func (pm *ProxyManager) Cleanup() error {
	var errs []error

	if pm.TargetA != nil {
		if err := pm.deleteProxy(pm.TargetA); err != nil {
			slog.Error("Failed to delete target-a proxy", "error", err)
			errs = append(errs, err)
		}
	}

	if pm.TargetB != nil {
		if err := pm.deleteProxy(pm.TargetB); err != nil {
			slog.Error("Failed to delete target-b proxy", "error", err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to cleanup some proxies: %v", errs)
	}

	return nil
}

// GetProxyURL returns the proxy URL for a target ("target-a" or "target-b").
func (pm *ProxyManager) GetProxyURL(targetName string) (string, error) {
	switch targetName {
	case "target-a":
		if pm.TargetA == nil {
			return "", fmt.Errorf("target-a proxy not initialized")
		}
		return pm.TargetA.ProxyURL, nil
	case "target-b":
		if pm.TargetB == nil {
			return "", fmt.Errorf("target-b proxy not initialized")
		}
		return pm.TargetB.ProxyURL, nil
	default:
		return "", fmt.Errorf("unknown target: %s", targetName)
	}
}

// createProxy creates a new proxy in Toxiproxy via REST API.
func (pm *ProxyManager) createProxy(target *ProxyTarget) error {
	payload := map[string]interface{}{
		"name":     target.ProxyName,
		"listen":   target.ListenSocket,
		"upstream": target.UpstreamSocket,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal proxy payload: %w", err)
	}

	url := fmt.Sprintf("%s/proxies", pm.ToxiproxyURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create proxy: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("toxiproxy returned status %d: %s", resp.StatusCode, string(respBody))
	}

	slog.Debug("Proxy created", "name", target.ProxyName, "listen", target.ListenSocket)
	return nil
}

// deleteProxy removes a proxy from Toxiproxy via REST API.
func (pm *ProxyManager) deleteProxy(target *ProxyTarget) error {
	url := fmt.Sprintf("%s/proxies/%s", pm.ToxiproxyURL, target.ProxyName)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete proxy: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("toxiproxy returned status %d: %s", resp.StatusCode, string(respBody))
	}

	slog.Debug("Proxy deleted", "name", target.ProxyName)
	return nil
}

// addToxic adds a toxic condition to a proxy via REST API.
func (pm *ProxyManager) addToxic(proxyName string, toxic *ToxicConfig) error {
	payload := map[string]interface{}{
		"name":       toxic.Name,
		"type":       toxic.Type,
		"stream":     toxic.Stream,
		"toxicity":   toxic.Toxicity,
		"attributes": toxic.Attributes,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal toxic payload: %w", err)
	}

	url := fmt.Sprintf("%s/proxies/%s/toxics", pm.ToxiproxyURL, proxyName)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to add toxic: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("toxiproxy returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Track applied toxic
	pm.ActiveToxics[proxyName] = append(pm.ActiveToxics[proxyName], toxic)

	slog.Debug("Toxic added",
		"proxy", proxyName,
		"toxic_type", toxic.Type,
		"stream", toxic.Stream,
	)

	return nil
}

// removeToxicsFromProxy removes all toxics from a specific proxy.
func (pm *ProxyManager) removeToxicsFromProxy(proxyName string) error {
	toxics, ok := pm.ActiveToxics[proxyName]
	if !ok || len(toxics) == 0 {
		return nil
	}

	for _, toxic := range toxics {
		url := fmt.Sprintf("%s/proxies/%s/toxics/%s", pm.ToxiproxyURL, proxyName, toxic.Name)
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			return fmt.Errorf("failed to create delete toxic request: %w", err)
		}

		resp, err := (&http.Client{}).Do(req)
		if err != nil {
			return fmt.Errorf("failed to delete toxic: %w", err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
			return fmt.Errorf("toxiproxy returned status %d when deleting toxic", resp.StatusCode)
		}

		slog.Debug("Toxic removed", "proxy", proxyName, "toxic", toxic.Name)
	}

	pm.ActiveToxics[proxyName] = []*ToxicConfig{}
	return nil
}

// Helper functions

// findFreePort finds an available TCP port.
func findFreePort() (int, error) {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: 0,
		IP:   net.ParseIP("127.0.0.1"),
	})
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	return port, nil
}

// getPort extracts the port from a URL, returning the default for the scheme if not specified.
func getPort(u *url.URL) string {
	if u.Port() != "" {
		return u.Port()
	}

	switch u.Scheme {
	case "https":
		return "443"
	case "http":
		fallthrough
	default:
		return "80"
	}
}
