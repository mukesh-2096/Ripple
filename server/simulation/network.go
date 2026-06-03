package simulation

import (
	"math/rand"
	"time"
)

type NetworkProfile struct {
	Name       string        `json:"name"`
	Latency    time.Duration `json:"latency"`
	PacketLoss float64       `json:"packet_loss"` // percentage (0.0 to 1.0)
}

var Profiles = map[string]NetworkProfile{
	"5G":   {Name: "5G", Latency: 20 * time.Millisecond, PacketLoss: 0.0},
	"4G":   {Name: "4G", Latency: 100 * time.Millisecond, PacketLoss: 0.0},
	"3G":   {Name: "3G", Latency: 300 * time.Millisecond, PacketLoss: 0.01},
	"2G":   {Name: "2G", Latency: 800 * time.Millisecond, PacketLoss: 0.03},
	"Slow": {Name: "Slow", Latency: 2000 * time.Millisecond, PacketLoss: 0.05},
}

// ShouldDrop simulates packet loss based on profile threshold
func (p NetworkProfile) ShouldDrop() bool {
	if p.PacketLoss <= 0 {
		return false
	}
	return rand.Float64() < p.PacketLoss
}
