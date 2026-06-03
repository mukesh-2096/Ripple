import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL || 'http://localhost:5000';

export default function ContactPage() {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setIsLoading(true);

    try {
      const res = await fetch(`${BACKEND_URL}/api/contact`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name, email, message }),
      });

      let errorMessage = 'Submission failed. Please try again.';
      try {
        const data = await res.json();
        errorMessage = data.message || errorMessage;
      } catch (parseErr) {
        try {
          const text = await res.text();
          if (text) errorMessage = text;
        } catch (_) {}
      }

      if (res.ok) {
        setSuccess('Thank you! Your message has been submitted. Check your inbox for confirmation email.');
        setName('');
        setEmail('');
        setMessage('');
      } else {
        setError(errorMessage);
      }
    } catch (err) {
      setError('Connection to server failed.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-slate-950 flex flex-col items-center justify-center px-6 py-12">
      <div className="bg-slate-905 border border-slate-800 rounded-2xl p-8 max-w-lg w-full shadow-2xl text-slate-100 relative">
        <div className="absolute -top-10 -right-10 w-40 h-40 bg-cyan-650 rounded-full blur-[80px] opacity-10 pointer-events-none" />

        <div className="flex justify-between items-center mb-6">
          <button onClick={() => navigate('/')} className="text-xs text-slate-400 hover:text-slate-200">
            ← Back to Home
          </button>
          <span className="text-xs font-semibold text-cyan-400">Mandate 2: Contact Page</span>
        </div>

        <div className="mb-6">
          <h2 className="text-2xl font-bold bg-gradient-to-r from-cyan-400 to-indigo-400 bg-clip-text text-transparent">
            Contact Ripple Support
          </h2>
          <p className="text-slate-400 text-xs mt-1">Submit your message live and verify your inbox confirmation</p>
        </div>

        {error && (
          <div className="bg-rose-955/50 border border-rose-800 text-rose-450 text-xs p-3 rounded-lg mb-4 text-center">
            {error}
          </div>
        )}

        {success && (
          <div className="bg-emerald-955/50 border border-emerald-800 text-emerald-450 text-xs p-3 rounded-lg mb-4 text-center">
            {success}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-xs font-semibold text-slate-400 mb-1.5">Full Name</label>
            <input
              type="text"
              required
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full bg-slate-800 border border-slate-700 rounded-xl px-4 py-2.5 text-xs text-slate-200 focus:outline-none focus:ring-2 focus:ring-cyan-500"
              placeholder="V D S Mukesh"
            />
          </div>
          <div>
            <label className="block text-xs font-semibold text-slate-400 mb-1.5">Email address</label>
            <input
              type="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full bg-slate-800 border border-slate-700 rounded-xl px-4 py-2.5 text-xs text-slate-200 focus:outline-none focus:ring-2 focus:ring-cyan-500"
              placeholder="you@domain.com"
            />
          </div>
          <div>
            <label className="block text-xs font-semibold text-slate-400 mb-1.5">Your Message</label>
            <textarea
              required
              rows={4}
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              className="w-full bg-slate-800 border border-slate-700 rounded-xl p-3 text-xs text-slate-200 focus:outline-none focus:ring-2 focus:ring-cyan-500"
              placeholder="Describe your inquiry..."
            />
          </div>

          <button
            type="submit"
            disabled={isLoading}
            className="w-full bg-gradient-to-r from-cyan-600 to-indigo-650 hover:from-cyan-500 hover:to-indigo-500 text-white font-semibold text-xs py-3 rounded-xl transition duration-150 disabled:opacity-55"
          >
            {isLoading ? 'Sending Inquiry...' : 'Submit Form'}
          </button>
        </form>
      </div>
    </div>
  );
}
