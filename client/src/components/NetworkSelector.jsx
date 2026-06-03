import React from 'react';

const PROFILES = [
  { name: '5G', latency: '20ms', loss: '0%', color: 'from-cyan-500 to-blue-500' },
  { name: '4G', latency: '100ms', loss: '0%', color: 'from-blue-500 to-indigo-500' },
  { name: '3G', latency: '300ms', loss: '1%', color: 'from-amber-500 to-orange-500' },
  { name: '2G', latency: '800ms', loss: '3%', color: 'from-orange-500 to-red-500' },
  { name: 'Slow', latency: '2000ms', loss: '5%', color: 'from-red-600 to-rose-700' },
];

export default function NetworkSelector({ selectedProfile, onSelectProfile }) {
  return (
    <div className="bg-slate-900 border border-slate-800 rounded-2xl p-6 text-slate-100 shadow-2xl">
      <h2 className="text-xl font-bold mb-4 bg-gradient-to-r from-cyan-400 to-blue-400 bg-clip-text text-transparent">
        Network Profile Simulation
      </h2>
      <p className="text-xs text-slate-400 mb-6">
        Simulate real-world network latency and packet loss rates dynamically on the proxy server.
      </p>
      
      <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
        {PROFILES.map((p) => {
          const isSelected = selectedProfile === p.name;
          return (
            <button
              key={p.name}
              onClick={() => onSelectProfile(p.name)}
              className={`relative overflow-hidden rounded-xl p-4 text-left border transition-all duration-200 ${
                isSelected 
                  ? 'border-cyan-400 bg-slate-800/80 shadow-lg shadow-cyan-900/20' 
                  : 'border-slate-800 bg-slate-850 hover:border-slate-700'
              }`}
            >
              {isSelected && (
                <div className={`absolute top-0 left-0 right-0 h-1 bg-gradient-to-r ${p.color}`} />
              )}
              <div className="flex justify-between items-center mb-2">
                <span className="font-bold text-lg">{p.name}</span>
                <span className={`h-2.5 w-2.5 rounded-full ${isSelected ? 'bg-cyan-400 animate-pulse' : 'bg-slate-700'}`} />
              </div>
              <div className="space-y-1">
                <div className="flex justify-between text-xs">
                  <span className="text-slate-400">Latency:</span>
                  <span className="font-semibold text-slate-200">{p.latency}</span>
                </div>
                <div className="flex justify-between text-xs">
                  <span className="text-slate-400">Packet Loss:</span>
                  <span className="font-semibold text-rose-400">{p.loss}</span>
                </div>
              </div>
            </button>
          );
        })}
      </div>
    </div>
  );
}
