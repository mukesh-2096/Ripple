/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package types

import "time"

// NetworkProfile defines network simulation parameters.
type NetworkProfile struct {
	Name        string
	Description string
	Latency     time.Duration
	Jitter      time.Duration
	PacketLoss  float64 // percentage 0-100
	Bandwidth   int64   // bytes per second, 0 = unlimited
}

// ShadowConfig holds the configuration for the shadow client.
type ShadowConfig struct {
	TargetA    string
	TargetB    string
	Profile    NetworkProfile
	OutputFile string
	Verbose    bool
	DiffFormat string // "json" or "text"
}

// PredefinedProfiles returns a map of commonly used network profiles.
func PredefinedProfiles() map[string]NetworkProfile {
	return map[string]NetworkProfile{
		"local": {
			Name:        "local",
			Description: "No network simulation (LAN)",
			Latency:     0,
			Jitter:      0,
			PacketLoss:  0,
			Bandwidth:   0,
		},
		"good-4g": {
			Name:        "good-4g",
			Description: "Good 4G network conditions",
			Latency:     50 * time.Millisecond,
			Jitter:      10 * time.Millisecond,
			PacketLoss:  0.1,
			Bandwidth:   4 * 1024 * 1024, // 4 Mbps
		},
		"poor-4g": {
			Name:        "poor-4g",
			Description: "Poor 4G network conditions",
			Latency:     150 * time.Millisecond,
			Jitter:      50 * time.Millisecond,
			PacketLoss:  1.0,
			Bandwidth:   1 * 1024 * 1024, // 1 Mbps
		},
		"3g": {
			Name:        "3g",
			Description: "3G network conditions",
			Latency:     200 * time.Millisecond,
			Jitter:      100 * time.Millisecond,
			PacketLoss:  2.0,
			Bandwidth:   512 * 1024, // 512 Kbps
		},
		"satellite": {
			Name:        "satellite",
			Description: "High latency satellite link",
			Latency:     500 * time.Millisecond,
			Jitter:      100 * time.Millisecond,
			PacketLoss:  0.5,
			Bandwidth:   2 * 1024 * 1024, // 2 Mbps
		},
		"custom": {
			Name:        "custom",
			Description: "Custom network parameters (specify latency-ms, jitter-ms, packet-loss)",
			Latency:     100 * time.Millisecond,
			Jitter:      10 * time.Millisecond,
			PacketLoss:  0.0,
			Bandwidth:   0,
		},
	}
}

// GetProfile returns a predefined profile or the custom profile if not found.
func GetProfile(name string) NetworkProfile {
	profiles := PredefinedProfiles()
	if profile, ok := profiles[name]; ok {
		return profile
	}
	return profiles["custom"]
}
