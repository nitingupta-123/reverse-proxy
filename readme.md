# Reverse Proxy with Middleware in Go

This Go program demonstrates a simple reverse proxy server with middleware using the `net/http` package. The proxy server intercepts incoming HTTP requests, modifies them, and forwards them to a target URL. The implementation includes middleware for IP address filtering and rate limiting.

## Overview

- The program defines an IP middleware that restricts access based on a predefined set of allowed IP addresses.
- It includes a rate limiter middleware that enforces a maximum number of requests per minute.
- The `handleRequest` function acts as the main handler for incoming requests. It modifies the response body and forwards the request to a target URL.
- The server runs on port 80 and responds to requests on the root ("/") and "/qa.example.com/" paths.

## Code Description

### Middleware Functions

#### IPMiddleware

```go
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
