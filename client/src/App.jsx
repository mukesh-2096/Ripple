import React, { useState, useEffect } from 'react';
import { Routes, Route, Navigate, useNavigate } from 'react-router-dom';
import LandingPage from './components/LandingPage';
import Login from './components/Login';
import Signup from './components/Signup';
import ContactPage from './components/ContactPage';
import RequestBuilder from './components/RequestBuilder';
import NetworkSelector from './components/NetworkSelector';
import ResponsePanel from './components/ResponsePanel';
import DiffView from './components/DiffView';
import LoadTester from './components/LoadTester';
import AnalyticsDashboard from './components/AnalyticsDashboard';

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL || 'http://localhost:5000';

export default function App() {
  const [userToken, setUserToken] = useState(localStorage.getItem('token') || '');
  const [username, setUsername] = useState(localStorage.getItem('username') || '');

  const navigate = useNavigate();

  const handleLoginSuccess = (token, user) => {
    localStorage.setItem('token', token);
    localStorage.setItem('username', user);
    setUserToken(token);
    setUsername(user);
    navigate('/dashboard');
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    setUserToken('');
    setUsername('');
    navigate('/');
  };

  return (
    <Routes>
      <Route path="/" element={<LandingPage />} />
      <Route 
        path="/login" 
        element={userToken ? <Navigate to="/dashboard" replace /> : <Login onLoginSuccess={handleLoginSuccess} />} 
      />
      <Route 
        path="/signup" 
        element={userToken ? <Navigate to="/dashboard" replace /> : <Signup onSignupSuccess={() => navigate('/login')} />} 
      />
      <Route path="/contact" element={<ContactPage />} />
      <Route 
        path="/dashboard" 
        element={
          userToken ? (
            <Dashboard 
              userToken={userToken} 
              username={username} 
              onLogout={handleLogout} 
            />
          ) : (
            <Navigate to="/login" replace />
          )
        } 
      />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}

// Inner component for the protected Dashboard workspace view
function Dashboard({ userToken, username, onLogout }) {
  const [activeTab, setActiveTab] = useState('proxy');
  const [selectedProfile, setSelectedProfile] = useState('5G');
  
  // Proxy state
  const [isLoadingProxy, setIsLoadingProxy] = useState(false);
  const [responseA, setResponseA] = useState(null);
  const [responseB, setResponseB] = useState(null);
  const [durationA, setDurationA] = useState(undefined);
  const [durationB, setDurationB] = useState(undefined);
  const [diffs, setDiffs] = useState(null);

  // Load test state
  const [isLoadingLoadTest, setIsLoadingLoadTest] = useState(false);
  const [loadTestReport, setLoadTestReport] = useState(null);

  // Analytics state
  const [analyticsData, setAnalyticsData] = useState(null);

  const fetchAnalytics = async () => {
    try {
      const res = await fetch(`${BACKEND_URL}/api/analytics`);
      if (res.ok) {
        const data = await res.json();
        setAnalyticsData(data);
      }
    } catch (e) {
      console.error('Failed to fetch analytics:', e);
    }
  };

  useEffect(() => {
    fetchAnalytics();
  }, []);

  const handleProxySend = async (reqDetails) => {
    setIsLoadingProxy(true);
    setResponseA(null);
    setResponseB(null);
    setDurationA(undefined);
    setDurationB(undefined);
    setDiffs(null);

    // 1. Send Request A through proxy with selected network simulation
    let respAContent = '';
    try {
      const start = Date.now();
      const res = await fetch(`${BACKEND_URL}/api/proxy`, {
        method: reqDetails.method,
        headers: {
          ...reqDetails.headers,
          'X-Target-URL': reqDetails.url,
          'X-Network-Profile': selectedProfile,
          'Authorization': `Bearer ${userToken}`
        },
        body: reqDetails.method !== 'GET' ? reqDetails.body : undefined,
      });
      setDurationA(Date.now() - start);
      respAContent = await res.text();
      setResponseA(respAContent);
    } catch (e) {
      setResponseA(`Error: ${e.message}`);
    }

    // 2. Send Request B directly to the URL (or proxy with 5G as baseline)
    let respBContent = '';
    try {
      const start = Date.now();
      const res = await fetch(`${BACKEND_URL}/api/proxy`, {
        method: reqDetails.method,
        headers: {
          ...reqDetails.headers,
          'X-Target-URL': reqDetails.url,
          'X-Network-Profile': '5G', // baseline
          'Authorization': `Bearer ${userToken}`
        },
        body: reqDetails.method !== 'GET' ? reqDetails.body : undefined,
      });
      setDurationB(Date.now() - start);
      respBContent = await res.text();
      setResponseB(respBContent);
    } catch (e) {
      setResponseB(`Error: ${e.message}`);
    }

    // 3. Compute Diff if both are valid JSON
    if (respAContent && respBContent) {
      try {
        const res = await fetch(`${BACKEND_URL}/api/diff`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            json1: respAContent,
            json2: respBContent,
          }),
        });
        if (res.ok) {
          const diffResult = await res.json();
          setDiffs(diffResult);
        }
      } catch (e) {
        console.error('Diff error:', e);
      }
    }

    setIsLoadingProxy(false);
    fetchAnalytics();
  };

  const handleRunLoadTest = async (testDetails) => {
    setIsLoadingLoadTest(true);
    setLoadTestReport(null);

    try {
      const res = await fetch(`${BACKEND_URL}/api/loadtest`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          url: testDetails.url,
          method: 'GET',
          concurrency: testDetails.concurrency,
          requests: testDetails.requests,
        }),
      });
      if (res.ok) {
        const report = await res.json();
        setLoadTestReport(report);
      }
    } catch (e) {
      console.error('Load test failed:', e);
    }
    setIsLoadingLoadTest(false);
  };

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans">
      {/* Header */}
      <header className="border-b border-slate-900 bg-slate-950/80 backdrop-blur-md sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-6 py-4 flex justify-between items-center">
          <div className="flex items-center gap-3">
            <div className="h-9 w-9 rounded-xl bg-gradient-to-tr from-violet-600 to-indigo-600 flex items-center justify-center font-bold text-white shadow-lg shadow-violet-500/20">
              R
            </div>
            <span className="font-extrabold text-xl tracking-tight bg-gradient-to-r from-white via-slate-200 to-slate-400 bg-clip-text text-transparent">
              Ripple
            </span>
          </div>

          <div className="flex items-center gap-6">
            <nav className="flex gap-1.5 bg-slate-900 border border-slate-850 p-1.5 rounded-xl">
              {[
                { id: 'proxy', label: 'Proxy & Simulation' },
                { id: 'loadtest', label: 'Load Tester' },
                { id: 'analytics', label: 'Analytics' },
              ].map(tab => (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={`px-4 py-2 rounded-lg text-xs font-semibold transition-all duration-150 ${
                    activeTab === tab.id 
                      ? 'bg-slate-800 text-white shadow-md' 
                      : 'text-slate-400 hover:text-slate-200'
                  }`}
                >
                  {tab.label}
                </button>
              ))}
            </nav>

            <div className="h-6 w-px bg-slate-800" />

            <div className="flex items-center gap-3">
              <span className="text-xs text-slate-450">
                Logged in as <strong className="text-slate-250 font-semibold">{username}</strong>
              </span>
              <button
                onClick={onLogout}
                className="text-xs text-rose-450 hover:text-rose-350 bg-rose-955/20 hover:bg-rose-955/40 px-3 py-1.5 rounded-xl border border-rose-900/50 transition-colors"
              >
                Log Out
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content Area */}
      <main className="flex-1 max-w-7xl w-full mx-auto px-6 py-8">
        {activeTab === 'proxy' && (
          <div className="space-y-6">
            <NetworkSelector 
              selectedProfile={selectedProfile} 
              onSelectProfile={setSelectedProfile} 
            />
            <RequestBuilder 
              onSendRequest={handleProxySend} 
              isLoading={isLoadingProxy} 
            />
            <ResponsePanel 
              response1={responseA} 
              response2={responseB} 
              duration1={durationA} 
              duration2={durationB} 
            />
            <DiffView diffs={diffs} />
          </div>
        )}

        {activeTab === 'loadtest' && (
          <LoadTester 
            onRunLoadTest={handleRunLoadTest} 
            isLoading={isLoadingLoadTest} 
            report={loadTestReport} 
          />
        )}

        {activeTab === 'analytics' && (
          <AnalyticsDashboard 
            data={analyticsData} 
            onRefresh={fetchAnalytics} 
          />
        )}
      </main>

      {/* Footer */}
      <footer className="border-t border-slate-900 py-6 bg-slate-950">
        <div className="max-w-7xl mx-auto px-6 text-center text-xs text-slate-500">
          Built with ❤️ at a Hackathon Sprint. Licensed under MIT.
        </div>
      </footer>
    </div>
  );
}
