package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/justinas/alice"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(rate.Every(time.Minute), 30) // Allow 3 requests per minute

var customTransport = http.DefaultTransport

// IPMiddleware is a middleware that restricts access based on IP address
func IPMiddleware(allowedIPs map[string]bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the request's IP is in the allowedIPs set
			clientIP := strings.Split(r.RemoteAddr, ":")[0]
			if !allowedIPs[clientIP] {
				http.Error(w, "Forbidden", http.StatusForbidden)
				fmt.Printf("Access denied for IP: %s\n", clientIP)
				return
			}

			// Call the next handler if the IP is allowed
			next.ServeHTTP(w, r)
		})
	}
}

// rateLimiterMiddleware is a middleware that enforces rate limiting
func rateLimiterMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if limiter.Allow() {
				// Process the request
			} else {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// handleRequest is the main handler for incoming requests
func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Modify the response body

	// Create a new HTTP request with the same method, URL, and body as the original request
	targetURL := "http://localhost:3000/"

	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}

	// Copy the headers from the original request to the proxy request
	for name, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}

	// Send the proxy request using the custom transport
	resp, err := customTransport.RoundTrip(proxyReq)
	if err != nil {
		http.Error(w, "Error sending proxy request", http.StatusInternalServerError)
		return
	}

	// Copy the headers from the proxy response to the original response
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5501")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "application/json")

	// Set the status code of the original response to the status code of the proxy response
	w.WriteHeader(resp.StatusCode)

	// Copy the body of the proxy response to the original response
	io.Copy(w, resp.Body)
}

func main() {
	// Default handler for root endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
		fmt.Println("Host:", r.Host)
	})

	allowedIPs := map[string]bool{
		"127.0.0.1": true, // Add your allowed IP addresses here
		// Add more IP addresses as needed
	}

	// Create a new middleware chain
	middlewareChain := alice.New(
		IPMiddleware(allowedIPs),
		rateLimiterMiddleware(),
	)

	// Attach middleware chain to the specific path
	http.Handle("/qa.example.com/", middlewareChain.ThenFunc(handleRequest))

	// Start the server
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Error starting proxy server: ", err)
	}
}
