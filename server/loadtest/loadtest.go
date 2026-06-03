package loadtest

import (
	"bytes"
	"io"
	"net/http"
	"sort"
	"sync"
	"time"
)

type TestResult struct {
	StatusCode   int
	ResponseTime time.Duration
	Error        bool
}

type LoadTestReport struct {
	TotalRequests   int           `json:"total_requests"`
	SuccessRate     float64       `json:"success_rate"`
	AvgResponseTime time.Duration `json:"avg_response_time"`
	P95Latency      time.Duration `json:"p95_latency"`
	FailedRequests  int           `json:"failed_requests"`
}

func RunLoadTest(url string, method string, headers map[string]string, body []byte, concurrency int, totalRequests int) LoadTestReport {
	results := make(chan TestResult, totalRequests)
	var wg sync.WaitGroup

	reqChan := make(chan struct{}, totalRequests)
	for i := 0; i < totalRequests; i++ {
		reqChan <- struct{}{}
	}
	close(reqChan)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range reqChan {
				req, err := http.NewRequest(method, url, bytes.NewReader(body))
				if err != nil {
					results <- TestResult{Error: true}
					continue
				}

				for k, v := range headers {
					req.Header.Set(k, v)
				}

				start := time.Now()
				resp, err := client.Do(req)
				duration := time.Since(start)

				if err != nil {
					results <- TestResult{Error: true, ResponseTime: duration}
					continue
				}

				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()

				results <- TestResult{
					StatusCode:   resp.StatusCode,
					ResponseTime: duration,
					Error:        resp.StatusCode >= 400,
				}
			}
		}()
	}

	wg.Wait()
	close(results)

	var durations []time.Duration
	successCount := 0
	failedCount := 0
	var totalDuration time.Duration

	for res := range results {
		if res.Error {
			failedCount++
		} else {
			successCount++
		}
		if res.ResponseTime > 0 {
			durations = append(durations, res.ResponseTime)
			totalDuration += res.ResponseTime
		}
	}

	avgResponseTime := time.Duration(0)
	if len(durations) > 0 {
		avgResponseTime = totalDuration / time.Duration(len(durations))
	}

	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})

	p95 := time.Duration(0)
	if len(durations) > 0 {
		idx := int(float64(len(durations)) * 0.95)
		if idx >= len(durations) {
			idx = len(durations) - 1
		}
		p95 = durations[idx]
	}

	successRate := 0.0
	if totalRequests > 0 {
		successRate = float64(successCount) / float64(totalRequests) * 100
	}

	return LoadTestReport{
		TotalRequests:   totalRequests,
		SuccessRate:     successRate,
		AvgResponseTime: avgResponseTime,
		P95Latency:      p95,
		FailedRequests:  failedCount,
	}
}
