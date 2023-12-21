// package main

// import (
// 	"io"
// 	"log"
// 	"net/http"
// )

// var customTransport = http.DefaultTransport

// func handleRequest(w http.ResponseWriter, r *http.Request) {
// 	// Create a new HTTP request with the same method, URL, and body as the original request
// 	targetURL := "https://info.cern.ch/"

// 	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
// 	if err != nil {
// 		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
// 		return
// 	}

// 	// Copy the headers from the original request to the proxy request
// 	for name, values := range r.Header {
// 		for _, value := range values {
// 			proxyReq.Header.Add(name, value)
// 		}
// 	}

// 	// Send the proxy request using the custom transport
// 	resp, err := customTransport.RoundTrip(proxyReq)
// 	if err != nil {
// 		http.Error(w, "Error sending proxy request", http.StatusInternalServerError)
// 		return
// 	}

// 	// Copy the headers from the proxy response to the original response
// 	for name, values := range resp.Header {
// 		for _, value := range values {
// 			w.Header().Add(name, value)
// 		}
// 	}

// 	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5501")
// 	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
// 	w.Header().Set("Access-Control-Allow-Headers", "application/json")

// 	// Set the status code of the original response to the status code of the proxy response
// 	w.WriteHeader(resp.StatusCode)

// 	// Copy the body of the proxy response to the original response
// 	io.Copy(w, resp.Body)
// }

// func main() {

// 	http.HandleFunc("/qa.example.com/", handleRequest)

// 	error := http.ListenAndServe(":80", nil)
// 	if error != nil {
// 		log.Fatal("Error starting proxy server: ", error)
// 	}

// }
