/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package types

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"time"
)

// ValidateURL checks if the provided string is a valid HTTP(S) URL.
func ValidateURL(urlStr string) error {
	if urlStr == "" {
		return errors.New("URL cannot be empty")
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL syntax: %w", err)
	}

	if parsedURL.Scheme == "" {
		return errors.New("URL must include a scheme (http or https)")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("unsupported scheme %q, must be http or https", parsedURL.Scheme)
	}

	if parsedURL.Host == "" {
		return errors.New("URL must include a host")
	}

	return nil
}

// ValidateShadowConfig performs comprehensive validation of the shadow client configuration.
func ValidateShadowConfig(config *ShadowConfig) error {
	if err := ValidateURL(config.TargetA); err != nil {
		return fmt.Errorf("invalid target-a: %w", err)
	}

	if err := ValidateURL(config.TargetB); err != nil {
		return fmt.Errorf("invalid target-b: %w", err)
	}

	if config.TargetA == config.TargetB {
		return errors.New("target-a and target-b cannot be identical")
	}

	if err := ValidateNetworkProfile(config.Profile); err != nil {
		return fmt.Errorf("invalid network profile: %w", err)
	}

	if err := ValidateDiffFormat(config.DiffFormat); err != nil {
		return fmt.Errorf("invalid diff format: %w", err)
	}

	return nil
}

// ValidateNetworkProfile validates network profile parameters.
func ValidateNetworkProfile(profile NetworkProfile) error {
	if profile.Latency < 0 {
		return errors.New("latency cannot be negative")
	}

	if profile.Jitter < 0 {
		return errors.New("jitter cannot be negative")
	}

	if profile.PacketLoss < 0 || profile.PacketLoss > 100 {
		return fmt.Errorf("packet loss must be between 0 and 100, got %f", profile.PacketLoss)
	}

	if profile.Bandwidth < 0 {
		return errors.New("bandwidth cannot be negative")
	}

	if profile.Latency > 10*time.Second {
		return fmt.Errorf("latency exceeds reasonable maximum (10s), got %v", profile.Latency)
	}

	return nil
}

// ValidateDiffFormat validates the diff output format.
func ValidateDiffFormat(format string) error {
	validFormats := map[string]bool{
		"json": true,
		"text": true,
	}

	if !validFormats[format] {
		return fmt.Errorf("invalid diff format %q, must be json or text", format)
	}

	return nil
}

// ValidateOutputFile checks if the output file path is valid.
func ValidateOutputFile(filePath string) error {
	if filePath == "" {
		return nil // stdout is valid
	}

	// Simple check for invalid characters in filename
	invalidCharsPattern := regexp.MustCompile(`[<>"|?*]`)
	if invalidCharsPattern.MatchString(filePath) {
		return fmt.Errorf("output file contains invalid characters: %s", filePath)
	}

	return nil
}
