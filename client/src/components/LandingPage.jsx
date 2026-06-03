import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL || 'http://localhost:5000';

export default function LandingPage() {
  const [isConnected, setIsConnected] = useState(null); // null = checking, true = connected, false = disconnected
  const navigate = useNavigate();

  useEffect(() => {
    const checkConnection = async () => {
      try {
        const res = await fetch(`${BACKEND_URL}/health`);
        if (res.ok) {
          setIsConnected(true);
        } else {
          setIsConnected(false);
        }
      } catch (e) {
        setIsConnected(false);
      }
    };
    checkConnection();
  }, []);

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans overflow-hidden">
      {/* Hero Section */}
      <div className="relative flex-1 flex flex-col items-center justify-center text-center px-6 py-20 max-w-5xl mx-auto z-10">
        <div className="absolute top-1/4 left-1/2 -translate-x-1/2 -translate-y-1/2 w-80 h-80 bg-violet-650 rounded-full blur-[120px] opacity-20 pointer-events-none" />
        <div className="absolute top-1/3 left-1/4 w-72 h-72 bg-cyan-650 rounded-full blur-[100px] opacity-15 pointer-events-none" />

        {/* Backend Status indicator */}
        <div className="mb-6 flex flex-col items-center gap-2">
          {isConnected === null && (
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-slate-900 border border-slate-800 text-xs font-semibold text-slate-400">
              <span className="h-2 w-2 rounded-full bg-slate-500 animate-pulse" />
              Checking backend connection...
            </div>
          )}
          {isConnected === true && (
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-950/40 border border-emerald-800/80 text-xs font-semibold text-emerald-450 shadow-lg shadow-emerald-900/10">
              <span className="h-2 w-2 rounded-full bg-emerald-400 animate-pulse" />
              Backend Connected
            </div>
          )}
          {isConnected === false && (
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-rose-955/40 border border-rose-800/80 text-xs font-semibold text-rose-455 shadow-lg shadow-rose-900/10">
              <span className="h-2 w-2 rounded-full bg-rose-500" />
              Backend Disconnected
            </div>
          )}
        </div>

        <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-slate-900 border border-slate-800 text-xs font-semibold text-violet-400 mb-6">
          ⚡ Network Simulation & Concurrency testing
        </div>

        <h1 className="text-5xl md:text-7xl font-extrabold tracking-tight mb-6 leading-none">
          Simulate Network Lag <br />
          <span className="bg-gradient-to-r from-violet-400 via-indigo-400 to-cyan-400 bg-clip-text text-transparent">
            Stress Test API Endpoints
          </span>
        </h1>

        <p className="text-slate-400 text-base md:text-lg max-w-2xl mb-10 leading-relaxed">
          Ripple is a developer platform that aggregates network profile simulations, concurrent load testing, payload diff comparison, and real-time GORM analytics into one dashboard.
        </p>

        <div className="flex flex-wrap gap-4 justify-center">
          <button
            onClick={() => navigate('/signup')}
            className="px-8 py-3.5 rounded-xl bg-gradient-to-r from-violet-600 to-indigo-600 hover:from-violet-500 hover:to-indigo-500 text-white font-semibold text-sm transition-all shadow-lg shadow-violet-900/30 hover:scale-105 duration-150"
          >
            Get Started (Free Account)
          </button>
          <button
            onClick={() => navigate('/login')}
            className="px-8 py-3.5 rounded-xl bg-slate-900 hover:bg-slate-805 border border-slate-800 text-slate-200 font-semibold text-sm transition-all duration-150"
          >
            Log In
          </button>
          <button
            onClick={() => navigate('/contact')}
            className="px-8 py-3.5 rounded-xl bg-slate-950 hover:bg-slate-900 border border-slate-850 text-slate-400 font-semibold text-sm transition-all duration-150"
          >
            Contact Sales
          </button>
        </div>
      </div>

      {/* Feature Grid */}
      <section className="border-t border-slate-900 bg-slate-950/50 py-20 px-6">
        <div className="max-w-7xl mx-auto">
          <h2 className="text-center text-3xl font-extrabold mb-12 bg-gradient-to-r from-white to-slate-400 bg-clip-text text-transparent">
            Engineered For API Reliability
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div className="bg-slate-900/50 border border-slate-850 rounded-2xl p-6 hover:border-slate-800 transition-colors">
              <div className="h-10 w-10 rounded-lg bg-violet-900/50 text-violet-400 flex items-center justify-center font-bold text-lg mb-4">
                📶
              </div>
              <h3 className="text-lg font-bold mb-2">Network Emulation</h3>
              <p className="text-slate-400 text-sm">
                Simulate 5G, 4G, 3G, 2G, or Slow connections to see how packet drops and latencies impact client integrations.
              </p>
            </div>
            <div className="bg-slate-900/50 border border-slate-850 rounded-2xl p-6 hover:border-slate-800 transition-colors">
              <div className="h-10 w-10 rounded-lg bg-cyan-900/50 text-cyan-400 flex items-center justify-center font-bold text-lg mb-4">
                ⚡
              </div>
              <h3 className="text-lg font-bold mb-2">Load Simulation</h3>
              <p className="text-slate-400 text-sm">
                Fire concurrent requests from Go goroutines to stress test APIs and measure response distribution.
              </p>
            </div>
            <div className="bg-slate-900/50 border border-slate-850 rounded-2xl p-6 hover:border-slate-800 transition-colors">
              <div className="h-10 w-10 rounded-lg bg-emerald-900/50 text-emerald-450 flex items-center justify-center font-bold text-lg mb-4">
                🔍
              </div>
              <h3 className="text-lg font-bold mb-2">JSON Schema Diff</h3>
              <p className="text-slate-400 text-sm">
                Detect schema breaking changes instantly with a recursive JSON validation and comparison diff engine.
              </p>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}
