import React from 'react';

export default function DiffView({ diffs }) {
  if (!diffs) {
    return (
      <div className="bg-slate-900 border border-slate-800 rounded-2xl p-6 text-slate-400 shadow-2xl text-center">
        No difference logs to display. Send requests or compare payloads to generate diff analysis.
      </div>
    );
  }

  if (diffs.length === 0) {
    return (
      <div className="bg-slate-900 border border-slate-850 rounded-2xl p-6 text-emerald-400 shadow-2xl text-center font-semibold">
        ✓ Payloads are identical. No diffs found!
      </div>
    );
  }

  const getBadgeStyle = (type) => {
    switch (type) {
      case 'ADDED':
        return 'bg-emerald-950 text-emerald-400 border-emerald-800';
      case 'REMOVED':
        return 'bg-rose-950 text-rose-400 border-rose-800';
      case 'TYPE_MISMATCH':
        return 'bg-amber-950 text-amber-400 border-amber-800';
      default:
        return 'bg-violet-950 text-violet-400 border-violet-855';
    }
  };

  return (
    <div className="bg-slate-900 border border-slate-800 rounded-2xl p-6 text-slate-100 shadow-2xl">
      <h2 className="text-xl font-bold mb-4 bg-gradient-to-r from-teal-400 to-emerald-400 bg-clip-text text-transparent">
        JSON Diff Engine Analysis
      </h2>
      
      <div className="overflow-x-auto">
        <table className="w-full text-left text-xs border-collapse">
          <thead>
            <tr className="border-b border-slate-800 text-slate-450 uppercase tracking-wider font-semibold">
              <th className="py-3 px-4">Field Path</th>
              <th className="py-3 px-4">Change Type</th>
              <th className="py-3 px-4">Original (A)</th>
              <th className="py-3 px-4">New (B)</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-850">
            {diffs.map((d, index) => (
              <tr key={index} className="hover:bg-slate-850/50 transition-colors duration-150">
                <td className="py-3.5 px-4 font-mono font-semibold text-slate-300">
                  {d.path || '(root)'}
                </td>
                <td className="py-3.5 px-4">
                  <span className={`inline-block px-2 py-0.5 rounded-full border text-[10px] font-bold ${getBadgeStyle(d.type)}`}>
                    {d.type}
                  </span>
                </td>
                <td className="py-3.5 px-4 font-mono text-rose-400 max-w-xs truncate">
                  {d.val1 !== undefined ? JSON.stringify(d.val1) : '-'}
                </td>
                <td className="py-3.5 px-4 font-mono text-emerald-400 max-w-xs truncate">
                  {d.val2 !== undefined ? JSON.stringify(d.val2) : '-'}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
