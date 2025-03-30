# HTTP/gRPC Inspector

A utility for logging and debugging HTTP and gRPC requests. Perfect for quickly inspecting what requests are being sent to your server.

## Features

- HTTP Request Logging:
  - Method (GET, POST, etc.)
  - URL
  - Headers
  - Request Body
  - Basic Auth (automatically decoded)
  - Client IP Address
  - Timestamp

- gRPC Request Logging:
  - Service Method
  - Metadata (headers)
  - Request Body
  - Client IP Address
  - Timestamp

## Quick Start

1. Download the binary:
```bash
curl -L https://github.com/komly/http-inspector/releases/download/v1.2.0/http-inspector-linux-amd64 -o http-inspector
chmod +x http-inspector
```

2. Run:
```bash
# Default: HTTP on 8080, gRPC on 9090
./http-inspector

# Or specify custom ports
HTTP_PORT=8888 GRPC_PORT=9999 ./http-inspector
```

## Usage Examples

### HTTP Requests

```bash
# Simple GET request
curl http://localhost:8080/test

# POST request with JSON
curl -X POST -H "Content-Type: application/json" -d '{"test":"data"}' http://localhost:8080/api

# Request with Basic Auth
curl -u user:pass http://localhost:8080/
```

### gRPC Requests

The utility supports gRPC reflection, so you can use grpcurl for testing:

```bash
# List available services
grpcurl -plaintext localhost:9090 list

# Send request with metadata
grpcurl -v -H 'authorization: Bearer token123' -d '{"key": "value"}' -plaintext localhost:9090 your.service.Method
```

## Configuration

Environment variables:
- `HTTP_PORT`: HTTP server port (default: 8080)
- `GRPC_PORT`: gRPC server port (default: 9090)

## gRPC Settings

- Keepalive parameters:
  - MinTime: 5 seconds
  - MaxConnectionIdle: 15 seconds
  - MaxConnectionAge: 30 seconds
  - Timeout: 1 second

## Running in Container

```bash
docker run -p 8080:8080 -p 9090:9090 http-inspector
```

## Running as systemd Service

1. Create `/etc/systemd/system/http-inspector.service`:
```ini
[Unit]
Description=HTTP/gRPC Inspector Service
After=network.target

[Service]
ExecStart=/path/to/http-inspector
Restart=always
Environment=HTTP_PORT=8080
Environment=GRPC_PORT=9090

[Install]
WantedBy=multi-user.target
```

2. Enable the service:
```bash
sudo systemctl enable http-inspector
sudo systemctl start http-inspector
```

## Example Output

### HTTP Request
```json
{
  "method": "POST",
  "url": "/api/test",
  "headers": {
    "Content-Type": ["application/json"],
    "Authorization": ["Basic dXNlcjpwYXNz"]
  },
  "body": "{\"test\":\"data\"}",
  "remote_addr": "127.0.0.1:52431",
  "host": "localhost:8080",
  "protocol": "HTTP/1.1",
  "timestamp": "2024-03-30T12:34:56Z",
  "basic_auth": {
    "username": "user",
    "password": "pass"
  }
}
```

### gRPC Request
```json
{
  "method": "gRPC",
  "protocol": "gRPC",
  "remote_addr": "127.0.0.1:52432",
  "timestamp": "2024-03-30T12:34:57Z",
  "grpc_method": "/service.v1.TestService/Method",
  "grpc_metadata": {
    "authorization": ["Bearer token123"]
  },
  "body": "{\"key\":\"value\"}"
}
``` 