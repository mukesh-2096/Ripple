/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package proxy

// Example usage of the ProxyManager:
//
// package main
//
// import (
//     "log"
//     "github.com/your-username/shadow-client/pkg/proxy"
//     "github.com/your-username/shadow-client/types"
// )
//
// func main() {
//     // Initialize proxy manager pointing to running Toxiproxy instance
//     pm := proxy.NewProxyManager("http://localhost:8474")
//
//     // Prepare shadow client configuration
//     config := &types.ShadowConfig{
//         TargetA: "http://api.example.com:8080",
//         TargetB: "http://api.example.com:8081",
//         Profile: types.GetProfile("good-4g"),
//     }
//
//     // Initialize proxies (dynamically allocates ports)
//     if err := pm.InitializeProxies(config); err != nil {
//         log.Fatalf("Failed to initialize proxies: %v", err)
//     }
//
//     defer pm.Cleanup()
//
//     // Apply network profile (latency, jitter, bandwidth, packet loss)
//     if err := pm.ApplyNetworkProfile(config.Profile); err != nil {
//         log.Fatalf("Failed to apply network profile: %v", err)
//     }
//
//     // Get proxy URLs to use in actual API calls
//     proxyUrlA, _ := pm.GetProxyURL("target-a")
//     proxyUrlB, _ := pm.GetProxyURL("target-b")
//
//     // Now use proxyUrlA and proxyUrlB instead of original URLs
//     // Requests through these URLs will have network conditions applied
//
//     // Later, remove all toxics while keeping proxies alive
//     if err := pm.RemoveAllToxics(); err != nil {
//         log.Fatalf("Failed to remove toxics: %v", err)
//     }
// }

// ProxyManager provides an interface to Toxiproxy for simulating network conditions.
// It automatically handles:
// - Dynamic port allocation for proxy listeners
// - Proxy creation and deletion
// - Toxic condition application (latency, jitter, bandwidth, packet loss)
// - Cleanup of proxies
//
// Requirements:
// - Toxiproxy server running (typically on localhost:8474)
// - Both target URLs must be reachable from Toxiproxy
// - Sufficient privileges to bind to allocated ports
