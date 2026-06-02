#  Ripple

> **Go-based API testing & network simulation platform — HTTP proxying, latency & packet-loss simulation, goroutine load testing, JSON diff analysis, and performance analytics dashboard.**

---

## 🧩 Problem Statement

Testing backend APIs for reliability and backward compatibility under real-world network conditions is currently fragmented. Developers must stitch together separate tools to fire requests, throttle network speeds, run load tests, and manually copy-paste responses into diff tools — creating a slow, error-prone workflow that leaves performance issues undiscovered until production.

## 💡 Solution

A lightweight Go-based developer tool that unifies request firing, network simulation, load testing, and response comparison into a single workflow. It leverages Go's native concurrency to simultaneously stress-test endpoints with configurable traffic loads through an artificially throttled network layer — simulating real-world lag and packet loss — and outputs instant side-by-side payload diffs with a performance analytics dashboard tracking response times, error rates, and p95 latency.

---

## ✨ Features

| Feature | Description |
|---|---|
| 🔀 **HTTP Request Proxying** | Route requests through Ripple's Go proxy to any local or remote backend endpoint |
| 📶 **Network Simulation** | Simulate 5G, 4G, 3G, 2G, and Slow network profiles with configurable latency and packet-loss injection |
| ⚡ **Goroutine Load Testing** | Fire N concurrent requests using Go goroutines — see how your backend handles real traffic spikes |
| 🔍 **JSON Diff Analysis** | Recursive field-level comparison of two API responses — detects added, removed, changed, and type-mismatched fields |
| ⚖️ **Response Comparison** | Side-by-side comparison of two responses across different network conditions or endpoint versions |
| 📊 **Performance Analytics** | Dashboard tracking response time trends, error rates, p95 latency, and slowest endpoints over time |

---

## 🛠 Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go (`net/http`, `goroutines`, `sync.WaitGroup`, `GORM`) |
| Frontend | React + Vite |
| Database | PostgreSQL |
| Styling | Tailwind CSS |

---

## 🏗 Project Structure

```
ripple/
├── server/                  # Go backend
│   ├── main.go              # Entry point
│   ├── proxy/               # HTTP proxy + latency injection
│   │   └── proxy.go
│   ├── simulation/          # Network profile definitions
│   │   └── network.go
│   ├── loadtest/            # Goroutine-based load tester
│   │   └── loadtest.go
│   ├── diff/                # JSON diff engine
│   │   └── diff.go
│   ├── analytics/           # Analytics API handlers
│   │   └── analytics.go
│   └── db/                  # PostgreSQL models + GORM setup
│       └── db.go
└── client/                  # React frontend
    ├── src/
    │   ├── components/
    │   │   ├── RequestBuilder.jsx     # Method, URL, headers, body
    │   │   ├── NetworkSelector.jsx    # 5G / 4G / 3G / 2G / Slow
    │   │   ├── ResponsePanel.jsx      # Side-by-side response view
    │   │   ├── DiffView.jsx           # JSON diff with colour coding
    │   │   ├── LoadTester.jsx         # Concurrency slider + results
    │   │   └── AnalyticsDashboard.jsx # Charts + metrics
    │   └── App.jsx
    └── vite.config.js
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 15+

### 1. Clone the repo

```bash
git clone https://github.com/mukesh-2096/ripple.git
cd ripple
```

### 2. Set up the database

```bash
# Create a PostgreSQL database named ripple
createdb ripple
```

### 3. Start the Go backend

```bash
cd server
go mod tidy
go run main.go
# Runs on http://localhost:5000
```

### 4. Start the React frontend

```bash
cd client
npm install
npm run dev
# Opens on http://localhost:3000
```

### 5. Open Ripple

Go to `http://localhost:3000` in your browser, point it at your own backend (e.g. `http://localhost:8000/api/login`), and start testing.

---

## 🌐 Network Profiles

| Profile | Simulated Latency | Packet Loss |
|---|---|---|
| 5G | 20ms | 0% |
| 4G | 100ms | 0% |
| 3G | 300ms | 1% |
| 2G | 800ms | 3% |
| Slow | 2000ms | 5% |

---

## ⚡ Load Testing

Ripple uses Go goroutines to fire N concurrent requests at an endpoint and captures:

- **Total requests** sent
- **Success rate** (2xx responses)
- **Average response time**
- **p95 latency** (95th percentile)
- **Failed requests** (timeouts, errors)

```go
// Under the hood — simplified
var wg sync.WaitGroup
results := make(chan Result, concurrency)

for i := 0; i < concurrency; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        results <- fireRequest(url, method, headers, body)
    }()
}

wg.Wait()
close(results)
```

---

## 📊 Analytics Dashboard

Every request made through Ripple is logged to PostgreSQL. The dashboard shows:

- Response time over time (line chart)
- Error rate trends (bar chart)
- Slowest endpoints ranked
- Network profile performance comparison

---

## 👥 Team

**Team Name — The Team**

| Name | GitHub |
|---|---|
| V D S Mukesh | [@mukesh-2096](https://github.com/mukesh-2096) |
| Abir Panda | [@AbirpandaA](https://github.com/AbirpandaA) |

---

## 📄 License

MIT License — feel free to use, modify, and distribute.

---

<p align="center">Built with ❤️ at a Hackathon Sprint</p>
