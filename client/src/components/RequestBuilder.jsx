import React, { useState } from 'react';

export default function RequestBuilder({ onSendRequest, isLoading }) {
  const [method, setMethod] = useState('GET');
  const [url, setUrl] = useState('https://jsonplaceholder.typicode.com/todos/1');
  const [headers, setHeaders] = useState([{ key: '', value: '' }]);
  const [body, setBody] = useState('{\n  "name": "RippleTest"\n}');

  const addHeader = () => setHeaders([...headers, { key: '', value: '' }]);
  const removeHeader = (index) => setHeaders(headers.filter((_, i) => i !== index));
  const updateHeader = (index, field, val) => {
    const updated = [...headers];
    updated[index][field] = val;
    setHeaders(updated);
  };

  const handleSend = () => {
    const headersObj = {};
    headers.forEach(h => {
      if (h.key && h.value) headersObj[h.key] = h.value;
    });
    onSendRequest({ method, url, headers: headersObj, body });
  };

  return (
    <div className="bg-slate-900 border border-slate-800 rounded-2xl p-6 text-slate-100 shadow-2xl">
      <h2 className="text-xl font-bold mb-4 bg-gradient-to-r from-violet-400 to-indigo-400 bg-clip-text text-transparent">
        Request Builder
      </h2>
      
      {/* Method & URL */}
      <div className="flex gap-3 mb-6">
        <select 
          value={method} 
          onChange={(e) => setMethod(e.target.value)}
          className="bg-slate-800 border border-slate-700 rounded-xl px-4 py-2.5 text-sm font-semibold focus:outline-none focus:ring-2 focus:ring-violet-500"
        >
          {['GET', 'POST', 'PUT', 'DELETE', 'PATCH'].map(m => (
            <option key={m} value={m}>{m}</option>
          ))}
        </select>
        <input 
          type="text" 
          value={url} 
          onChange={(e) => setUrl(e.target.value)} 
          placeholder="Enter API URL" 
          className="flex-1 bg-slate-800 border border-slate-700 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-violet-500"
        />
        <button 
          onClick={handleSend}
          disabled={isLoading}
          className="bg-gradient-to-r from-violet-600 to-indigo-600 hover:from-violet-500 hover:to-indigo-500 text-white font-semibold rounded-xl px-6 py-2.5 text-sm transition-all duration-200 shadow-lg shadow-violet-900/30 disabled:opacity-50"
        >
          {isLoading ? 'Sending...' : 'Send'}
        </button>
      </div>

      {/* Headers Section */}
      <div className="mb-6">
        <div className="flex justify-between items-center mb-2">
          <label className="text-sm font-medium text-slate-400">Headers</label>
          <button 
            type="button" 
            onClick={addHeader} 
            className="text-xs text-violet-400 hover:text-violet-300 font-semibold"
          >
            + Add Header
          </button>
        </div>
        <div className="space-y-2">
          {headers.map((h, i) => (
            <div key={i} className="flex gap-2 items-center">
              <input 
                type="text" 
                placeholder="Key" 
                value={h.key} 
                onChange={(e) => updateHeader(i, 'key', e.target.value)} 
                className="flex-1 bg-slate-800 border border-slate-700 rounded-lg px-3 py-1.5 text-xs focus:outline-none focus:ring-1 focus:ring-violet-500"
              />
              <input 
                type="text" 
                placeholder="Value" 
                value={h.value} 
                onChange={(e) => updateHeader(i, 'value', e.target.value)} 
                className="flex-1 bg-slate-800 border border-slate-700 rounded-lg px-3 py-1.5 text-xs focus:outline-none focus:ring-1 focus:ring-violet-500"
              />
              <button 
                type="button" 
                onClick={() => removeHeader(i)} 
                className="text-slate-500 hover:text-rose-400 text-xs px-2"
              >
                ✕
              </button>
            </div>
          ))}
        </div>
      </div>

      {/* Body Section */}
      <div>
        <label className="block text-sm font-medium text-slate-400 mb-2">JSON Body</label>
        <textarea 
          rows={5} 
          value={body} 
          onChange={(e) => setBody(e.target.value)} 
          className="w-full bg-slate-800 border border-slate-700 rounded-xl p-3 text-xs font-mono focus:outline-none focus:ring-2 focus:ring-violet-500 text-emerald-400"
        />
      </div>
    </div>
  );
}
