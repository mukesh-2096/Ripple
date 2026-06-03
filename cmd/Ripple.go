/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/your-username/shadow-client/types"
)

// Flag variables for the Ripple command
var (
	targetA    string
	targetB    string
	profile    string
	latencyMs  int
	jitterMs   int
	packetLoss float64
	bandwidth  int64
	outputFile string
	diffFormat string
	verbose    bool
)

// RippleCmd represents the Ripple command for API shadow client comparison
var RippleCmd = &cobra.Command{
	Use:   "ripple",
	Short: "Compare two API endpoints with network simulation",
	Long: `Ripple compares responses from two API endpoints, optionally applying network 
simulation profiles to measure behavior under different network conditions.

Example:
  shadow-client ripple --target-a https://api.example.com --target-b https://api.staging.example.com --profile good-4g
  shadow-client ripple --target-a http://localhost:8080 --target-b http://localhost:8081 --latency-ms 100 --packet-loss 1.5`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize logger
		handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
		if verbose {
			handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
		}
		slog.SetDefault(slog.New(handler))

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Build network profile
		selectedProfile := types.GetProfile(profile)

		// Override profile parameters if custom flags are provided
		if cmd.Flags().Changed("latency-ms") {
			selectedProfile.Latency = time.Duration(latencyMs) * time.Millisecond
		}
		if cmd.Flags().Changed("jitter-ms") {
			selectedProfile.Jitter = time.Duration(jitterMs) * time.Millisecond
		}
		if cmd.Flags().Changed("packet-loss") {
			selectedProfile.PacketLoss = packetLoss
		}
		if cmd.Flags().Changed("bandwidth") {
			selectedProfile.Bandwidth = bandwidth
		}

		// Create shadow configuration
		config := &types.ShadowConfig{
			TargetA:    targetA,
			TargetB:    targetB,
			Profile:    selectedProfile,
			OutputFile: outputFile,
			Verbose:    verbose,
			DiffFormat: diffFormat,
		}

		// Validate configuration
		if err := types.ValidateShadowConfig(config); err != nil {
			slog.Error("Configuration validation failed", "error", err)
			return fmt.Errorf("invalid configuration: %w", err)
		}

		if err := types.ValidateOutputFile(config.OutputFile); err != nil {
			slog.Error("Output file validation failed", "error", err)
			return fmt.Errorf("invalid output file: %w", err)
		}

		// Log configuration details
		slog.Info("Shadow client initialized",
			"target_a", config.TargetA,
			"target_b", config.TargetB,
			"profile", config.Profile.Name,
			"latency", config.Profile.Latency,
			"jitter", config.Profile.Jitter,
			"packet_loss", config.Profile.PacketLoss,
		)

		// Placeholder for actual comparison logic
		fmt.Printf("Ripple comparison started\n")
		fmt.Printf("Target A: %s\n", config.TargetA)
		fmt.Printf("Target B: %s\n", config.TargetB)
		fmt.Printf("Profile: %s (%s)\n", config.Profile.Name, config.Profile.Description)
		fmt.Printf("Latency: %v, Jitter: %v, Packet Loss: %.2f%%\n",
			config.Profile.Latency, config.Profile.Jitter, config.Profile.PacketLoss)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(RippleCmd)

	// Required flags
	RippleCmd.Flags().StringVar(&targetA, "target-a", "", "First API endpoint URL (required)")
	RippleCmd.Flags().StringVar(&targetB, "target-b", "", "Second API endpoint URL (required)")

	// Mark required flags
	RippleCmd.MarkFlagRequired("target-a")
	RippleCmd.MarkFlagRequired("target-b")

	// Network profile flags
	RippleCmd.Flags().StringVar(&profile, "profile", "local",
		"Network profile: local, good-4g, poor-4g, 3g, satellite, or custom")

	// Custom network parameters (override profile)
	RippleCmd.Flags().IntVar(&latencyMs, "latency-ms", 0,
		"Override profile latency in milliseconds")
	RippleCmd.Flags().IntVar(&jitterMs, "jitter-ms", 0,
		"Override profile jitter in milliseconds")
	RippleCmd.Flags().Float64Var(&packetLoss, "packet-loss", 0,
		"Override profile packet loss percentage (0-100)")
	RippleCmd.Flags().Int64Var(&bandwidth, "bandwidth", 0,
		"Override profile bandwidth in bytes per second (0 = unlimited)")

	// Output options
	RippleCmd.Flags().StringVar(&outputFile, "output", "",
		"Output file for results (default: stdout)")
	RippleCmd.Flags().StringVar(&diffFormat, "diff-format", "json",
		"Diff output format: json or text")

	// Flags
	RippleCmd.Flags().BoolVar(&verbose, "verbose", false,
		"Enable verbose logging")
}
