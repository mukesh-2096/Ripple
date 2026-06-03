import React, { useState } from 'react';

export default function LoadTester({ onRunLoadTest, isLoading, report }) {
  const [concurrency, setConcurrency] = useState(10);
  const [totalRequests, setTotalRequests] = useState(100);
  const [url, setUrl] = useState('https://jsonplaceholder.typicode.com/todos/1');

  const handleRun = () => {
    onRunLoadTest({ url, concurrency, requests: totalRequests });
  };

  return (
    <div className="bg-slate-900 border border-slate-800 rounded-2xl p-6 text-slate-100 shadow-2xl">
      <h2 className="text-xl font-bold mb-4 bg-gradient-to-r from-amber-400 to-orange-400 bg-clip-text text-transparent">
        Goroutine Load Tester
      </h2>
      <p className="text-xs text-slate-400 mb-6">
        Fire high-concurrency HTTP requests using Go's lightweight thread scheduling mechanisms (goroutines).
      </p>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
        <div>
          <label className="block text-xs font-semibold text-slate-400 mb-2">Target URL</label>
          <input 
            type="text" 
            value={url} 
            onChange={(e) => setUrl(e.target.value)} 
            className="w-full bg-slate-800 border border-slate-700 rounded-xl px-4 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-orange-500"
          />
        </div>
        <div>
          <label className="block text-xs font-semibold text-slate-400 mb-2">
            Concurrency (Goroutines): {concurrency}
          </label>
          <input 
            type="range" 
            min="1" 
            max="100" 
            value={concurrency} 
            onChange={(e) => setConcurrency(parseInt(e.target.value))}
            className="w-full h-1.5 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-orange-500"
          />
        </div>
        <div>
          <label className="block text-xs font-semibold text-slate-400 mb-2">
            Total Requests: {totalRequests}
          </label>
          <input 
            type="range" 
            min="10" 
            max="1000" 
            step="10"
            value={totalRequests} 
            onChange={(e) => setTotalRequests(parseInt(e.target.value))}
            className="w-full h-1.5 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-orange-500"
          />
        </div>
      </div>

      <div className="flex justify-end mb-6">
        <button
          onClick={handleRun}
          disabled={isLoading}
          className="bg-gradient-to-r from-orange-600 to-amber-600 hover:from-orange-500 hover:to-amber-500 text-white font-semibold rounded-xl px-6 py-2 text-sm transition-all duration-200 shadow-lg shadow-orange-900/30 disabled:opacity-50"
        >
          {isLoading ? 'Running Load Test...' : 'Start Load Test'}
        </button>
      </div>

      {report && (
        <div className="border-t border-slate-800 pt-6">
          <h3 className="font-bold text-sm text-slate-300 mb-4">Results Dashboard</h3>
          <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
            <div className="bg-slate-950 border border-slate-850 p-4 rounded-xl text-center">
              <div className="text-xs text-slate-500 font-semibold mb-1">Total Requests</div>
              <div className="text-xl font-bold text-slate-200">{report.total_requests}</div>
            </div>
            <div className="bg-slate-950 border border-slate-850 p-4 rounded-xl text-center">
              <div className="text-xs text-slate-500 font-semibold mb-1">Success Rate</div>
              <div className="text-xl font-bold text-emerald-400">{report.success_rate.toFixed(1)}%</div>
            </div>
            <div className="bg-slate-950 border border-slate-850 p-4 rounded-xl text-center">
              <div className="text-xs text-slate-500 font-semibold mb-1">Avg Latency</div>
              <div className="text-xl font-bold text-cyan-400">
                {(report.avg_response_time / 1000000).toFixed(1)} ms
              </div>
            </div>
            <div className="bg-slate-950 border border-slate-850 p-4 rounded-xl text-center">
              <div className="text-xs text-slate-500 font-semibold mb-1">p95 Latency</div>
              <div className="text-xl font-bold text-violet-400">
                {(report.p95_latency / 1000000).toFixed(1)} ms
              </div>
            </div>
            <div className="bg-slate-950 border border-slate-850 p-4 rounded-xl text-center col-span-2 md:col-span-1">
              <div className="text-xs text-slate-500 font-semibold mb-1">Failed Requests</div>
              <div className="text-xl font-bold text-rose-450">{report.failed_requests}</div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
