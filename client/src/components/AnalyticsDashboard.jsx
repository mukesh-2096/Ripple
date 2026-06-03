import React from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

export default function AnalyticsDashboard({ data, onRefresh }) {
  return (
    <div className="bg-slate-900 border border-slate-800 rounded-2xl p-6 text-slate-100 shadow-2xl">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h2 className="text-xl font-bold bg-gradient-to-r from-emerald-400 to-cyan-400 bg-clip-text text-transparent">
            Performance Analytics
          </h2>
          <p className="text-xs text-slate-400">
            Real-time metrics aggregated directly from transaction logs saved inside PostgreSQL database.
          </p>
        </div>
        <button
          onClick={onRefresh}
          className="text-xs font-semibold text-cyan-400 hover:text-cyan-300 bg-slate-850 hover:bg-slate-800 border border-slate-800 px-3.5 py-2 rounded-xl transition-all"
        >
          Refresh Data
        </button>
      </div>

      {!data ? (
        <div className="text-center py-10 text-slate-500 text-xs">
          Loading performance telemetry or database connection offline.
        </div>
      ) : (
        <div className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="bg-slate-950 border border-slate-850 rounded-xl p-4">
              <span className="text-xs text-slate-500 font-semibold block mb-1">Total Transactions Logged</span>
              <span className="text-2xl font-bold text-slate-200">{data.total_requests || 0}</span>
            </div>
            <div className="bg-slate-950 border border-slate-850 rounded-xl p-4">
              <span className="text-xs text-slate-500 font-semibold block mb-1">Error Rate Average</span>
              <span className={`text-2xl font-bold ${data.error_rate > 5 ? 'text-rose-400' : 'text-emerald-450'}`}>
                {(data.error_rate || 0).toFixed(2)}%
              </span>
            </div>
            <div className="bg-slate-950 border border-slate-850 rounded-xl p-4">
              <span className="text-xs text-slate-500 font-semibold block mb-1">p95 Latency Average</span>
              <span className="text-2xl font-bold text-violet-400">{(data.p95_latency || 0).toFixed(1)} ms</span>
            </div>
          </div>

          <div>
            <h3 className="text-sm font-semibold text-slate-300 mb-3">Slowest Endpoint Responses</h3>
            {data.slowest_routes && data.slowest_routes.length > 0 ? (
              <div className="h-64">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={data.slowest_routes} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" />
                    <XAxis dataKey="url" stroke="#64748b" fontSize={10} tickLine={false} />
                    <YAxis stroke="#64748b" fontSize={10} unit="ms" tickLine={false} />
                    <Tooltip 
                      contentStyle={{ backgroundColor: '#0f172a', borderColor: '#1e293b', borderRadius: 8 }}
                      labelStyle={{ color: '#94a3b8', fontSize: 11, fontWeight: 'bold' }}
                      itemStyle={{ color: '#38bdf8', fontSize: 11 }}
                    />
                    <Bar dataKey="avg_duration" fill="#0284c7" radius={[4, 4, 0, 0]} name="Avg Latency" />
                  </BarChart>
                </ResponsiveContainer>
              </div>
            ) : (
              <div className="text-center py-6 text-slate-600 text-xs border border-dashed border-slate-800 rounded-xl">
                No latency records available yet. Make requests via the proxy to populate records.
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
