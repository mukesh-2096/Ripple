package analytics

import (
	"encoding/json"
	"net/http"
	"time"

	"ripple/db"
)

type AnalyticsSummary struct {
	TotalRequests int                    `json:"total_requests"`
	ErrorRate     float64                `json:"error_rate"`
	P95Latency    float64                `json:"p95_latency"`
	SlowestRoutes []RoutePerformance     `json:"slowest_routes"`
	History       []TimePerformancePoint `json:"history"`
}

type RoutePerformance struct {
	URL          string  `json:"url"`
	AvgDuration  float64 `json:"avg_duration"`
	RequestCount int64   `json:"request_count"`
}

type TimePerformancePoint struct {
	Time        time.Time `json:"time"`
	AvgDuration float64   `json:"avg_duration"`
	ErrorCount  int64     `json:"error_count"`
}

func HandleGetAnalytics(w http.ResponseWriter, r *http.Request) {
	if db.DB == nil {
		http.Error(w, "database connection not initialized", http.StatusInternalServerError)
		return
	}

	var totalCount int64
	db.DB.Model(&db.RequestLog{}).Count(&totalCount)

	var errorCount int64
	db.DB.Model(&db.RequestLog{}).Where("status_code >= 400 OR status_code = 0").Count(&errorCount)

	errorRate := 0.0
	if totalCount > 0 {
		errorRate = float64(errorCount) / float64(totalCount) * 100
	}

	// Calculate p95 latency
	var p95 float64
	var responseTimes []float64
	db.DB.Model(&db.RequestLog{}).Order("response_time asc").Pluck("response_time", &responseTimes)
	if len(responseTimes) > 0 {
		idx := int(float64(len(responseTimes)) * 0.95)
		if idx >= len(responseTimes) {
			idx = len(responseTimes) - 1
		}
		p95 = responseTimes[idx]
	}

	// Get slowest routes
	var slowest []RoutePerformance
	db.DB.Model(&db.RequestLog{}).
		Select("url, avg(response_time) as avg_duration, count(*) as request_count").
		Group("url").
		Order("avg_duration desc").
		Limit(5).
		Scan(&slowest)

	summary := AnalyticsSummary{
		TotalRequests: int(totalCount),
		ErrorRate:     errorRate,
		P95Latency:    p95,
		SlowestRoutes: slowest,
		History:       []TimePerformancePoint{}, // Optional list for charts
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
