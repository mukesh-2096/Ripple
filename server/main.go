package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"ripple/analytics"
	"ripple/auth"
	"ripple/db"
	"ripple/diff"
	"ripple/email"
	"ripple/loadtest"
	"ripple/proxy"
)

func enableCors(w *http.ResponseWriter) {
	origin := os.Getenv("CLIENT_URL")
	if origin == "" {
		origin = "*"
	}
	(*w).Header().Set("Access-Control-Allow-Origin", origin)
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Target-URL, X-Network-Profile")
}

func sendJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func handleProxyRequest(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}

	targetURL := r.Header.Get("X-Target-URL")
	profileName := r.Header.Get("X-Network-Profile")

	if targetURL == "" {
		http.Error(w, "Missing X-Target-URL header", http.StatusBadRequest)
		return
	}

	resp, duration, err := proxy.HandleProxy(w, r, profileName, targetURL)

	// Log to database in background
	statusCode := 0
	if resp != nil {
		statusCode = resp.StatusCode
	}
	if db.DB != nil {
		go func() {
			db.DB.Create(&db.RequestLog{
				Method:       r.Method,
				URL:          targetURL,
				StatusCode:   statusCode,
				ResponseTime: float64(duration.Milliseconds()),
				Profile:      profileName,
			})
		}()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	// Copy response headers and body
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	resp.Body.Close()
}

type DiffRequest struct {
	JSON1 string `json:"json1"`
	JSON2 string `json:"json2"`
}

func handleDiffRequest(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DiffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	diffs, err := diff.CompareJSON([]byte(req.JSON1), []byte(req.JSON2))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(diffs)
}

type LoadTestRequest struct {
	URL         string            `json:"url"`
	Method      string            `json:"method"`
	Headers     map[string]string `json:"headers"`
	Body        string            `json:"body"`
	Concurrency int               `json:"concurrency"`
	Requests    int               `json:"requests"`
}

func handleLoadTest(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoadTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Concurrency <= 0 {
		req.Concurrency = 10
	}
	if req.Requests <= 0 {
		req.Requests = 100
	}

	report := loadtest.RunLoadTest(req.URL, req.Method, req.Headers, []byte(req.Body), req.Concurrency, req.Requests)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func handleAnalytics(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	analytics.HandleGetAnalytics(w, r)
}

type SignupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func handleSignup(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	if r.Method != "POST" {
		sendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if db.DB == nil {
		sendJSONError(w, "Database connection offline", http.StatusInternalServerError)
		return
	}

	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		sendJSONError(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		sendJSONError(w, "Hashing failed", http.StatusInternalServerError)
		return
	}

	// Check if username or email already exists
	var existingUser db.User
	if db.DB != nil {
		if err := db.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
			sendJSONError(w, "Username is already taken", http.StatusConflict)
			return
		}
		if err := db.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			sendJSONError(w, "Email is already registered", http.StatusConflict)
			return
		}
	}

	user := db.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hash,
	}

	if db.DB != nil {
		if err := db.DB.Create(&user).Error; err != nil {
			sendJSONError(w, "Database error creating user", http.StatusInternalServerError)
			return
		}
	}

	// Trigger Automated Email (Mandate 3)
	go func() {
		subject := "Welcome to Ripple!"
		body := fmt.Sprintf("Hi %s,\n\nWelcome to Ripple - the Go-based network simulation & API load testing platform!\nYour account was successfully created.\n\nBest,\nThe Ripple Team", user.Username)
		email.SendEmail(user.Email, subject, body)
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})
}

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	if r.Method != "POST" {
		sendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if db.DB == nil {
		sendJSONError(w, "Database connection offline", http.StatusInternalServerError)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user db.User
	err := db.DB.Where("username = ? OR email = ?", req.UsernameOrEmail, req.UsernameOrEmail).First(&user).Error
	if err != nil {
		sendJSONError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		sendJSONError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(user.ID, user.Username)
	if err != nil {
		sendJSONError(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":    token,
		"username": user.Username,
	})
}

type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func handleContact(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	if r.Method != "POST" {
		sendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if db.DB == nil {
		sendJSONError(w, "Database connection offline", http.StatusInternalServerError)
		return
	}

	var req ContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Message == "" {
		sendJSONError(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	msg := db.ContactMessage{
		Name:    req.Name,
		Email:   req.Email,
		Message: req.Message,
	}

	if err := db.DB.Create(&msg).Error; err != nil {
		sendJSONError(w, "Failed to submit message", http.StatusInternalServerError)
		return
	}

	// Trigger Automated Email (Mandate 3)
	go func() {
		subject := "We received your message!"
		body := fmt.Sprintf("Hi %s,\n\nThank you for reaching out to us. We have received your contact message and our team will get back to you shortly.\n\nYour message:\n\"%s\"\n\nBest,\nThe Ripple Team", msg.Name, msg.Message)
		email.SendEmail(msg.Email, subject, body)
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Message submitted successfully",
	})
}

func handleCheckUsername(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if db.DB == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"available": true})
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Missing username parameter", http.StatusBadRequest)
		return
	}

	var count int64
	db.DB.Model(&db.User{}).Where("username = ?", username).Count(&count)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"available": count == 0,
	})
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	portFlag := flag.String("port", "5000", "Port to run the backend server on")
	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		port = *portFlag
	}

	// Setup PostgreSQL DSN
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Default DSN
		dsn = "host=localhost user=postgres password=postgres dbname=ripple port=5432 sslmode=disable"
	}

	log.Printf("Connecting to database...")
	err := db.InitDB(dsn)
	if err != nil {
		log.Printf("WARNING: Failed to connect to PostgreSQL database: %v. Running in log-only mock mode.", err)
	} else {
		log.Println("Database connection initialized and migrated successfully.")
	}

	http.HandleFunc("/api/proxy", handleProxyRequest)
	http.HandleFunc("/api/diff", handleDiffRequest)
	http.HandleFunc("/api/loadtest", handleLoadTest)
	http.HandleFunc("/api/analytics", handleAnalytics)
	http.HandleFunc("/api/auth/signup", handleSignup)
	http.HandleFunc("/api/auth/login", handleLogin)
	http.HandleFunc("/api/auth/check-username", handleCheckUsername)
	http.HandleFunc("/api/contact", handleContact)

	// Healthy check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
