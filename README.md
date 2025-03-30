# HTTP Request Inspector

A simple utility to inspect and log all incoming HTTP requests. Perfect for debugging HTTP requests in a containerized environment.

## Features

- Logs all HTTP request details including:
  - Method
  - URL
  - Headers
  - Body
  - Remote Address
  - Host
  - Request URI
  - Protocol
  - Timestamp
- Always returns 200 OK response
- JSON formatted logs
- Containerized application

## Quick Start

1. Build and run the container:
```bash
docker build -t http-inspector .
docker run -p 8080:8080 http-inspector
```

2. Send requests to inspect:
```bash
curl -X POST -H "Content-Type: application/json" -d '{"test":"data"}' http://localhost:8080
```

3. Check the logs:
```bash
docker logs <container_id>
```

## Environment Variables

- `PORT`: Port to listen on (default: 8080)

## Example Output

The utility will log each request in JSON format:

```json
{
  "method": "POST",
  "url": "/",
  "headers": {
    "Content-Type": ["application/json"],
    "User-Agent": ["curl/7.64.1"]
  },
  "body": "{\"test\":\"data\"}",
  "remote_addr": "172.17.0.1:12345",
  "host": "localhost:8080",
  "request_uri": "/",
  "protocol": "HTTP/1.1",
  "timestamp": "2024-03-21T10:00:00Z"
}
``` 