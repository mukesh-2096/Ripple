package proxy

import (
	"net/url"
	"time"
)

// ProxyTarget represents a single target being proxied through Toxiproxy.
type ProxyTarget struct {
	Name           string
	OriginalURL    string
	ParsedURL      *url.URL
	ProxyURL       string // URL through which to access the proxied target
	ProxyName      string // Toxiproxy proxy name
	UpstreamSocket string // Toxiproxy upstream socket (e.g., "127.0.0.1:8080")
	ListenSocket   string // Toxiproxy listen socket (e.g., "127.0.0.1:8081")
}

// ToxicConfig represents a toxic condition to apply to a proxy.
type ToxicConfig struct {
	Name       string
	Type       string      // "latency", "jitter", "bandwidth_limit", "packet_loss"
	Stream     string      // "upstream" or "downstream"
	Toxicity   float64     // 0.0 to 1.0, percentage of traffic affected
	Attributes map[string]interface{}
}

// LatencyToxic represents a latency toxic condition.
type LatencyToxic struct {
	Latency int `json:"latency"` // in milliseconds
	Jitter  int `json:"jitter"`  // in milliseconds
}

// BandwidthLimitToxic represents a bandwidth limit toxic condition.
type BandwidthLimitToxic struct {
	Rate int64 `json:"rate"` // in bytes per second
}

// PacketLossToxic represents a packet loss toxic condition.
type PacketLossToxic struct {
	Percentage float64 `json:"percentage"` // 0-100
}

// ProxyManager manages Toxiproxy proxies for shadow client operation.
type ProxyManager struct {
	ToxiproxyURL string
	TargetA      *ProxyTarget
	TargetB      *ProxyTarget
	ActiveToxics map[string][]*ToxicConfig // key: proxy name, value: list of applied toxics
}

// ProxyStats represents statistics for a proxy.
type ProxyStats struct {
	ProxyName        string
	Received         int64
	Sent             int64
	ActiveConnections int64
	AppliedToxics    []*ToxicConfig
	Created          time.Time
}
