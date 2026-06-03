import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL || 'http://localhost:5000';

export default function Login({ onLoginSuccess }) {
  const [usernameOrEmail, setUsernameOrEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const navigate = useNavigate();

  const handleLogin = async (e) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      const res = await fetch(`${BACKEND_URL}/api/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username_or_email: usernameOrEmail, password }),
      });

      let errorMessage = 'Invalid username or password.';
      let data = {};
      try {
        data = await res.json();
        errorMessage = data.message || errorMessage;
      } catch (parseErr) {
        try {
          const text = await res.text();
          if (text) errorMessage = text;
        } catch (_) {}
      }

      if (res.ok) {
        onLoginSuccess(data.token, data.username);
      } else {
        setError(errorMessage);
      }
    } catch (err) {
      setError('Connection to auth server failed.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-slate-950 flex flex-col items-center justify-center px-6 py-12">
      <div className="bg-slate-900 border border-slate-800 rounded-2xl p-8 max-w-md w-full shadow-2xl text-slate-100 relative">
        <div className="absolute -top-10 -left-10 w-40 h-40 bg-violet-650 rounded-full blur-[80px] opacity-10 pointer-events-none" />

        <div className="text-center mb-6">
          <div className="h-12 w-12 rounded-2xl bg-gradient-to-tr from-violet-655 to-indigo-655 flex items-center justify-center font-bold text-white shadow-lg mx-auto mb-3">
            R
          </div>
          <h2 className="text-2xl font-bold">Welcome Back</h2>
          <p className="text-slate-400 text-xs mt-1">Sign in to your Ripple workspace</p>
        </div>

        {error && (
          <div className="bg-rose-955/50 border border-rose-800 text-rose-450 text-xs p-3 rounded-lg mb-4 text-center">
            {error}
          </div>
        )}

        <form onSubmit={handleLogin} className="space-y-4">
          <div>
            <label className="block text-xs font-semibold text-slate-400 mb-1.5">Username or Email</label>
            <input
              type="text"
              required
              value={usernameOrEmail}
              onChange={(e) => setUsernameOrEmail(e.target.value)}
              className="w-full bg-slate-800 border border-slate-700 rounded-xl px-4 py-2.5 text-xs text-slate-200 focus:outline-none focus:ring-2 focus:ring-violet-500"
              placeholder="you@domain.com or username"
            />
          </div>
          <div>
            <label className="block text-xs font-semibold text-slate-400 mb-1.5">Password</label>
            <div className="relative">
              <input
                type={showPassword ? 'text' : 'password'}
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full bg-slate-805 border border-slate-700 rounded-xl pl-4 pr-10 py-2.5 text-xs text-slate-200 focus:outline-none focus:ring-2 focus:ring-violet-500"
                placeholder="••••••••"
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-xs text-slate-400 hover:text-slate-200 font-semibold"
              >
                {showPassword ? 'Hide' : 'Show'}
              </button>
            </div>
          </div>

          <button
            type="submit"
            disabled={isLoading}
            className="w-full bg-gradient-to-r from-violet-600 to-indigo-600 hover:from-violet-500 hover:to-indigo-500 text-white font-semibold text-xs py-3 rounded-xl transition duration-150 disabled:opacity-55"
          >
            {isLoading ? 'Signing In...' : 'Log In'}
          </button>
        </form>

        <div className="text-center mt-6 text-xs text-slate-450">
          New to Ripple?{' '}
          <button onClick={() => navigate('/signup')} className="text-violet-400 hover:text-violet-300 font-semibold">
            Sign Up
          </button>
        </div>
      </div>
    </div>
  );
}
