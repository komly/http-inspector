package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
)

type RequestInfo struct {
	Method       string              `json:"method"`
	URL          string              `json:"url,omitempty"`
	Headers      map[string][]string `json:"headers,omitempty"`
	Body         string              `json:"body,omitempty"`
	RemoteAddr   string              `json:"remote_addr"`
	Host         string              `json:"host,omitempty"`
	RequestURI   string              `json:"request_uri,omitempty"`
	Protocol     string              `json:"protocol"`
	Timestamp    string              `json:"timestamp"`
	BasicAuth    *BasicAuthInfo      `json:"basic_auth,omitempty"`
	GrpcMethod   string              `json:"grpc_method,omitempty"`
	GrpcMetadata metadata.MD         `json:"grpc_metadata,omitempty"`
}

type BasicAuthInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// gRPC interceptor для логирования запросов
func logInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	p, _ := peer.FromContext(ctx)
	addr := "unknown"
	if p != nil {
		addr = p.Addr.String()
	}

	reqInfo := RequestInfo{
		Method:       "gRPC",
		Protocol:     "gRPC",
		RemoteAddr:   addr,
		Timestamp:    time.Now().Format(time.RFC3339),
		GrpcMethod:   info.FullMethod,
		GrpcMetadata: md,
	}

	// Попытаемся сериализовать тело запроса
	if reqBytes, err := json.Marshal(req); err == nil {
		reqInfo.Body = string(reqBytes)
	}

	// Логируем запрос
	jsonData, _ := json.MarshalIndent(reqInfo, "", "  ")
	log.Printf("\n=== New gRPC Request ===\n%s\n===============\n", string(jsonData))

	// Вызываем оригинальный обработчик
	resp, err := handler(ctx, req)

	if err != nil {
		log.Printf("gRPC error: %v", err)
		return nil, err
	}

	return resp, nil
}

func main() {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "9090"
	}

	// Запускаем HTTP сервер
	go func() {
		http.HandleFunc("/", handleRequest)
		log.Printf("Starting HTTP Inspector on port %s...", httpPort)
		if err := http.ListenAndServe(":"+httpPort, nil); err != nil {
			log.Fatal(err)
		}
	}()

	// Запускаем gRPC сервер
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(logInterceptor),
	)

	// Включаем reflection для возможности использования grpcurl
	reflection.Register(s)

	log.Printf("Starting gRPC Inspector on port %s...", grpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
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
	log.Printf("\n=== New HTTP Request ===\n%s\n===============\n", string(jsonData))

	// Always return 200 OK
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Request logged successfully")
}
