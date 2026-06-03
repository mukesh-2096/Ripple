# 🚀 Ripple Project: Complete Setup & Execution Guide

This document provides a comprehensive, step-by-step setup guide to get the Ripple application (Go backend + React frontend) running on your local machine from scratch.

---

## 📂 1. Directory Structure Verification

Before running any commands, verify that your project folder structure matches the following layout:

```text
ripple/
├── server/                  # Go Backend
│   ├── main.go              # Main server entrypoint (ports, routing, handlers)
│   ├── auth/                # Hashing (bcrypt) & JWT generation
│   ├── db/                  # GORM initialization, models, and migrations
│   ├── email/               # SMTP configuration and sent email logs
│   ├── proxy/               # Reverse proxy implementation & latency injection
│   ├── loadtest/            # Go-concurrency load tester logic
│   ├── simulation/          # Preconfigured network latency profiles
│   ├── .env.example         # Template for server environment variables
│   └── .env                 # Server configuration (To be created)
├── client/                  # React Frontend
│   ├── src/                 # React components and layouts
│   ├── .env.example         # Template for client environment variables
│   ├── .env                 # Client configuration (To be created)
│   ├── vite.config.js       # Vite build & server config
│   └── package.json         # Frontend dependencies and npm scripts
├── load_test.js             # k6 automated load testing script (Mandate 4)
├── .gitignore               # Excludes vendor, node_modules, and .env files
└── SETUP_GUIDE.md           # This setup guide
```

---

## 📥 2. Step-by-Step Setup Instructions

### Step 2.1: Clone the Repository
Open your terminal (PowerShell, Command Prompt, or Git Bash) and execute:
```cmd
git clone https://github.com/mukesh-2096/ripple.git
cd ripple
```

---

### Step 2.2: Setup the PostgreSQL Database

Ripple uses GORM with a PostgreSQL database. You need to create a database before launching the backend to avoid running in mock-only mode.

#### Option A: Using pgAdmin (GUI)
1. Open **pgAdmin** and connect to your local database server.
2. Right-click on **Databases** -> **Create** -> **Database...**
3. Set the database name to: `ripple`
4. Set the Owner to: `postgres`
5. Click **Save**.

#### Option B: Using Command Line / terminal
If you have PostgreSQL binaries in your system path:
```cmd
createdb -U postgres ripple
```

---

### Step 2.3: Configure Environment Variables (`.env` files)

You must create `.env` files in both the `server` and `client` folders to connect them.

#### 1. Backend `.env` (`server/.env`)
Create a file named `.env` in the `server/` directory:
* **Path:** `ripple/server/.env`
* **Content:**
```env
PORT=5000
DATABASE_URL=postgresql://postgres:YOUR_PASSWORD@localhost:5432/ripple
CLIENT_URL=http://localhost:5173
```
> [!IMPORTANT]
> Replace `YOUR_PASSWORD` with your actual local PostgreSQL password.

#### 2. Frontend `.env` (`client/.env`)
Create a file named `.env` in the `client/` directory:
* **Path:** `ripple/client/.env`
* **Content:**
```env
VITE_BACKEND_URL=http://localhost:5000
```

---

### Step 2.4: Initialize & Start the Go Backend

1. Open a terminal and navigate to the `server` directory:
   ```cmd
   cd server
   ```
2. Initialize dependencies:
   ```cmd
   go mod tidy
   ```
3. Start the backend server:
   ```cmd
   go run main.go
   ```
   *(If Go is not in your environment variables, run using the absolute path: `& "C:\Program Files\Go\bin\go.exe" run main.go`)*

#### Expected Output:
```text
2026/06/03 15:00:00 Connecting to database...
2026/06/03 15:00:01 Database connection initialized and migrated successfully.
2026/06/03 15:00:01 Server listening on :5000
```
> [!NOTE]
> If you see `WARNING: Failed to connect to PostgreSQL database... Running in log-only mock mode`, make sure your PostgreSQL service is running, the database `ripple` exists, and your password in `server/.env` is correct.

---

### Step 2.5: Initialize & Start the React Frontend

1. Open a **second, separate** terminal window.
2. Navigate to the `client` directory:
   ```cmd
   cd client
   ```
3. Install frontend dependencies:
   ```cmd
   npm install
   ```
4. Start the Vite development server:
   ```cmd
   npm run dev
   ```

#### Expected Output:
```text
  VITE v8.0.16  ready in 320ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
```

---

## 🔍 3. Verifying the Connection

1. Open your browser and navigate to `http://localhost:5173/`.
2. Look at the status indicator badge on the landing page:
   - **🟢 Connected to Backend**: This confirms the frontend successfully fetched the `/health` API from your running Go backend.
   - **🔴 Not Connected to Backend**: Check that the backend server is running on port 5000, and your frontend `.env` has the correct `VITE_BACKEND_URL`.

---

## ⚡ 4. Verification & Testing

### Test User Signup & Login
- Navigate to the **Signup** page (`/signup`).
- Try signing up. If you enter an existing username or email, the backend will return a specific error payload (e.g. `Username already exists` or `Email already exists`).
- Toggle the password visibility icons to inspect input values.
- Navigate to the **Login** page (`/login`) to login with your registered credentials.

### Run Concurrency Load Testing (Mandate 4)
We've provided a [k6](https://k6.io/) load test script to verify that the Go server handles traffic spikes gracefully under high concurrency.
1. Download and install [k6](https://k6.io/docs/get-started/installation/).
2. Open a terminal in the root `ripple/` directory and run:
   ```cmd
   k6 run load_test.js
   ```
3. View the summary statistics table displayed in the console.

---

## 🛠 Troubleshooting

### Port 5000 or 5173 is Already in Use
If the terminal complains that a port is already in use:
- **Windows (PowerShell):** Find and kill the process using the port.
  ```powershell
  # Find the PID using port 5000
  Get-NetTCPConnection -LocalPort 5000 | Select-Object OwningProcess
  
  # Kill the process
  Stop-Process -Id <PID> -Force
  ```
- Alternatively, you can change the ports in your `.env` files.
