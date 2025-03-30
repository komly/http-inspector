package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type RequestInfo struct {
	Method     string              `json:"method"`
	URL        string              `json:"url"`
	Headers    map[string][]string `json:"headers"`
	Body       string              `json:"body"`
	RemoteAddr string              `json:"remote_addr"`
	Host       string              `json:"host"`
	RequestURI string              `json:"request_uri"`
	Protocol   string              `json:"protocol"`
	Timestamp  string              `json:"timestamp"`
	BasicAuth  *BasicAuthInfo      `json:"basic_auth,omitempty"`
}

type BasicAuthInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handleRequest)

	log.Printf("Starting HTTP Inspector on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
	}
	r.Body.Close()

	// Create request info
	reqInfo := RequestInfo{
		Method:     r.Method,
		URL:        r.URL.String(),
		Headers:    r.Header,
		Body:       string(body),
		RemoteAddr: r.RemoteAddr,
		Host:       r.Host,
		RequestURI: r.RequestURI,
		Protocol:   r.Proto,
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	// Parse Basic Auth if present
	if auth := r.Header.Get("Authorization"); auth != "" {
		if strings.HasPrefix(auth, "Basic ") {
			decoded, err := base64.StdEncoding.DecodeString(auth[6:])
			if err == nil {
				parts := strings.SplitN(string(decoded), ":", 2)
				if len(parts) == 2 {
					reqInfo.BasicAuth = &BasicAuthInfo{
						Username: parts[0],
						Password: parts[1],
					}
				}
			}
		}
	}

	// Marshal to JSON for pretty printing
	jsonData, err := json.MarshalIndent(reqInfo, "", "  ")
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Log the request details
	log.Printf("\n=== New Request ===\n%s\n===============\n", string(jsonData))

	// Always return 200 OK
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Request logged successfully")
}
