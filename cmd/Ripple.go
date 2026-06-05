/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/your-username/shadow-client/pkg/proxy"
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

		// Create proxy manager pointing to default Toxiproxy URL
		pm := proxy.NewProxyManager("http://localhost:8474")

		// Initialize proxies (dynamically allocates ports)
		slog.Info("Initializing Toxiproxy proxies...")
		if err := pm.InitializeProxies(config); err != nil {
			slog.Error("Failed to initialize proxies", "error", err)
			return fmt.Errorf("proxy initialization failed: %w", err)
		}
		// Ensure cleanup is deferred
		defer func() {
			slog.Info("Cleaning up proxies...")
			if err := pm.Cleanup(); err != nil {
				slog.Warn("Failed to cleanup proxies", "error", err)
			}
		}()

		// Apply network profile (latency, jitter, bandwidth, packet loss)
		slog.Info("Applying network profile to proxies...", "profile", config.Profile.Name)
		if err := pm.ApplyNetworkProfile(config.Profile); err != nil {
			slog.Error("Failed to apply network profile", "error", err)
			return fmt.Errorf("failed to apply network profile: %w", err)
		}

		// Retrieve proxy URLs to execute the requests
		proxyUrlA, err := pm.GetProxyURL("target-a")
		if err != nil {
			return fmt.Errorf("failed to get proxy URL for target-a: %w", err)
		}
		proxyUrlB, err := pm.GetProxyURL("target-b")
		if err != nil {
			return fmt.Errorf("failed to get proxy URL for target-b: %w", err)
		}

		slog.Info("Routing outbound requests through Toxiproxy...",
			"proxy_a", proxyUrlA,
			"proxy_b", proxyUrlB,
		)

		// Create HTTP Client with a timeout exceeding maximum profile latency (e.g., 30s)
		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		type reqResult struct {
			target   string
			duration time.Duration
			body     []byte
			status   int
			err      error
		}

		resultsChan := make(chan reqResult, 2)

		// Trigger concurrent requests
		slog.Info("Firing concurrent requests through proxies...")
		fireRequest := func(targetName, urlStr string) {
			startTime := time.Now()
			req, err := http.NewRequest("GET", urlStr, nil)
			if err != nil {
				resultsChan <- reqResult{target: targetName, err: err}
				return
			}

			resp, err := client.Do(req)
			duration := time.Since(startTime)
			if err != nil {
				resultsChan <- reqResult{target: targetName, duration: duration, err: err}
				return
			}
			defer resp.Body.Close()

			bodyBytes, err := io.ReadAll(resp.Body)
			resultsChan <- reqResult{
				target:   targetName,
				duration: duration,
				body:     bodyBytes,
				status:   resp.StatusCode,
				err:      err,
			}
		}

		go fireRequest("Target A", proxyUrlA)
		go fireRequest("Target B", proxyUrlB)

		var resA, resB reqResult
		for i := 0; i < 2; i++ {
			res := <-resultsChan
			if res.target == "Target A" {
				resA = res
			} else {
				resB = res
			}
		}

		// Display performance stats
		fmt.Println("\n==================================================")
		fmt.Println("             RIPPLE EXECUTION REPORT              ")
		fmt.Println("==================================================")
		fmt.Printf("Profile Applied: %s (%s)\n", config.Profile.Name, config.Profile.Description)
		fmt.Printf("Latency Profile: %v (Jitter: %v, Loss: %.1f%%, Bandwidth: %d Bps)\n\n",
			config.Profile.Latency, config.Profile.Jitter, config.Profile.PacketLoss, config.Profile.Bandwidth)

		// Print Target A results
		if resA.err != nil {
			fmt.Printf("🔴 Target A [%s]: FAILED\n   Error: %v\n", config.TargetA, resA.err)
		} else {
			fmt.Printf("🟢 Target A [%s]:\n   Status:   %d OK\n   Latency:  %v\n   Bytes:    %d\n",
				config.TargetA, resA.status, resA.duration, len(resA.body))
		}

		fmt.Println()

		// Print Target B results
		if resB.err != nil {
			fmt.Printf("🔴 Target B [%s]: FAILED\n   Error: %v\n", config.TargetB, resB.err)
		} else {
			fmt.Printf("🟢 Target B [%s]:\n   Status:   %d OK\n   Latency:  %v\n   Bytes:    %d\n",
				config.TargetB, resB.status, resB.duration, len(resB.body))
		}

		fmt.Println("==================================================")

		// Payload Diff Comparison
		if resA.err == nil && resB.err == nil {
			fmt.Println("Payload Match Check:")
			if bytes.Equal(resA.body, resB.body) {
				fmt.Println("  ✓ JSON Payloads match exactly!")
			} else {
				fmt.Println("  ✗ JSON Payloads differ.")
				// Fallback to simple string display or unmarshaled comparison
				var valA, valB interface{}
				errA := json.Unmarshal(resA.body, &valA)
				errB := json.Unmarshal(resB.body, &valB)
				if errA == nil && errB == nil {
					if fmt.Sprintf("%v", valA) == fmt.Sprintf("%v", valB) {
						fmt.Println("  ✓ Unmarshaled JSON payloads are structurally identical.")
					} else {
						fmt.Println("  [Warning] Responses have data mismatches.")
						if config.Verbose {
							fmt.Printf("  Target A response: %s\n", string(resA.body))
							fmt.Printf("  Target B response: %s\n", string(resB.body))
						}
					}
				} else {
					fmt.Println("  [Non-JSON or unparseable payload comparison]")
				}
			}
		}
		fmt.Println("==================================================")

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
