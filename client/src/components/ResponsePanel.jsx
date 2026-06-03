import React from 'react';

export default function ResponsePanel({ response1, response2, duration1, duration2 }) {
  const formatJSON = (jsonStr) => {
    try {
      if (!jsonStr) return '';
      if (typeof jsonStr === 'object') return JSON.stringify(jsonStr, null, 2);
      const parsed = JSON.parse(jsonStr);
      return JSON.stringify(parsed, null, 2);
    } catch (e) {
      return jsonStr;
    }
  };

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
      {/* Response 1 Panel */}
      <div className="bg-slate-900 border border-slate-800 rounded-2xl p-6 text-slate-100 shadow-2xl">
        <div className="flex justify-between items-center mb-4">
          <h3 className="font-bold text-lg text-slate-200">Response A</h3>
          {duration1 !== undefined && (
            <span className="text-xs bg-slate-800 text-cyan-400 font-semibold px-2.5 py-1 rounded-lg">
              {duration1} ms
            </span>
          )}
        </div>
        <pre className="bg-slate-950 border border-slate-850 rounded-xl p-4 text-xs font-mono text-emerald-400 overflow-x-auto max-h-96 min-h-48 text-left whitespace-pre-wrap">
          {response1 ? formatJSON(response1) : '// No response data yet'}
        </pre>
      </div>

      {/* Response 2 Panel */}
      <div className="bg-slate-900 border border-slate-800 rounded-2xl p-6 text-slate-100 shadow-2xl">
        <div className="flex justify-between items-center mb-4">
          <h3 className="font-bold text-lg text-slate-200">Response B</h3>
          {duration2 !== undefined && (
            <span className="text-xs bg-slate-800 text-violet-400 font-semibold px-2.5 py-1 rounded-lg">
              {duration2} ms
            </span>
          )}
        </div>
        <pre className="bg-slate-950 border border-slate-850 rounded-xl p-4 text-xs font-mono text-indigo-400 overflow-x-auto max-h-96 min-h-48 text-left whitespace-pre-wrap">
          {response2 ? formatJSON(response2) : '// No response data yet'}
        </pre>
      </div>
    </div>
  );
}
