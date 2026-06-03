import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL || 'http://localhost:5000';

export default function Signup({ onSignupSuccess }) {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  
  // Username availability state
  const [isCheckingUsername, setIsCheckingUsername] = useState(false);
  const [usernameStatus, setUsernameStatus] = useState(''); // 'available', 'taken', ''

  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const navigate = useNavigate();

  // Debounced username availability check
  useEffect(() => {
    if (!username.trim()) {
      setUsernameStatus('');
      return;
    }

    const delayDebounce = setTimeout(async () => {
      setIsCheckingUsername(true);
      setUsernameStatus('');
      try {
        const res = await fetch(`${BACKEND_URL}/api/auth/check-username?username=${encodeURIComponent(username)}`);
        if (res.ok) {
          const data = await res.json();
          if (data.available) {
            setUsernameStatus('available');
          } else {
            setUsernameStatus('taken');
          }
        }
      } catch (err) {
        console.error('Failed to check username availability:', err);
      } finally {
        setIsCheckingUsername(false);
      }
    }, 500); // 500ms debounce delay

    return () => clearTimeout(delayDebounce);
  }, [username]);

  const handleSignup = async (e) => {
    e.preventDefault();
    setError('');
    setSuccess('');

    if (usernameStatus === 'taken') {
      setError('Username is already taken. Please try another.');
      return;
    }

    if (password !== confirmPassword) {
      setError('Passwords do not match.');
      return;
    }

    setIsLoading(true);

    try {
      const res = await fetch(`${BACKEND_URL}/api/auth/signup`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, email, password }),
      });

      let errorMessage = 'Signup failed. Please try again.';
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
        setSuccess('Account created successfully! Redirecting to login...');
        setTimeout(() => {
          onSignupSuccess();
        }, 1500);
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
          <div className="h-12 w-12 rounded-2xl bg-gradient-to-tr from-violet-650 to-indigo-650 flex items-center justify-center font-bold text-white shadow-lg mx-auto mb-3">
            R
          </div>
          <h2 className="text-2xl font-bold">Create your Account</h2>
          <p className="text-slate-400 text-xs mt-1">Start running concurrency tests in seconds</p>
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

        <form onSubmit={handleSignup} className="space-y-4">
          <div>
            <div className="flex justify-between items-center mb-1.5">
              <label className="block text-xs font-semibold text-slate-400">Username</label>
              {isCheckingUsername && (
                <span className="text-[10px] text-slate-500 animate-pulse">Checking availability...</span>
              )}
              {!isCheckingUsername && usernameStatus === 'available' && (
                <span className="text-[10px] text-emerald-400 font-semibold">✓ Username available</span>
              )}
              {!isCheckingUsername && usernameStatus === 'taken' && (
                <span className="text-[10px] text-rose-400 font-semibold">✕ Username already taken, try another</span>
              )}
            </div>
            <input
              type="text"
              required
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className={`w-full bg-slate-800 border rounded-xl px-4 py-2.5 text-xs text-slate-200 focus:outline-none focus:ring-2 focus:ring-violet-500 ${
                usernameStatus === 'available' ? 'border-emerald-600' : usernameStatus === 'taken' ? 'border-rose-600' : 'border-slate-700'
              }`}
              placeholder="ripple_dev"
            />
          </div>
          
          <div>
            <label className="block text-xs font-semibold text-slate-400 mb-1.5">Email address</label>
            <input
              type="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full bg-slate-800 border border-slate-700 rounded-xl px-4 py-2.5 text-xs text-slate-200 focus:outline-none focus:ring-2 focus:ring-violet-500"
              placeholder="you@domain.com"
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

          <div>
            <label className="block text-xs font-semibold text-slate-400 mb-1.5">Confirm Password</label>
            <div className="relative">
              <input
                type={showConfirmPassword ? 'text' : 'password'}
                required
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                className="w-full bg-slate-800 border border-slate-700 rounded-xl pl-4 pr-10 py-2.5 text-xs text-slate-200 focus:outline-none focus:ring-2 focus:ring-violet-500"
                placeholder="••••••••"
              />
              <button
                type="button"
                onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-xs text-slate-400 hover:text-slate-200 font-semibold"
              >
                {showConfirmPassword ? 'Hide' : 'Show'}
              </button>
            </div>
          </div>

          <button
            type="submit"
            disabled={isLoading || isCheckingUsername}
            className="w-full bg-gradient-to-r from-violet-600 to-indigo-600 hover:from-violet-500 hover:to-indigo-500 text-white font-semibold text-xs py-3 rounded-xl transition duration-150 disabled:opacity-55"
          >
            {isLoading ? 'Creating Account...' : 'Sign Up'}
          </button>
        </form>

        <div className="text-center mt-6 text-xs text-slate-400">
          Already have an account?{' '}
          <button onClick={() => navigate('/login')} className="text-violet-400 hover:text-violet-300 font-semibold">
            Log In
          </button>
        </div>
      </div>
    </div>
  );
}
