# NetCheck API

A comprehensive web API built with Go and Gin for checking SSL certificates, HTTP/3 support, DNS information, IP details, and web server settings.

## Features

- **SSL Certificate Analysis**: Check certificate validity, issuer, expiration, and key details
- **HTTP/3 Support Detection**: Test if a domain supports HTTP/3 protocol
- **DNS Information**: Get A, AAAA, CNAME, MX, TXT, and NS records
- **IP Information**: Basic IP address validation and connection testing
- **Web Server Settings**: Analyze HTTP headers, server information, and response details
- **Comprehensive Check**: Get all information in a single request

## Installation

1. Clone or download the project
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Build the application:
   ```bash
   go build -o netcheck main.go
   ```
4. Run the server:
   ```bash
   ./netcheck
   ```

The API will be available at `http://localhost:8080`

## API Endpoints

### Health Check
- **GET** `/api/v1/health`
- Returns API status and version information

### SSL Certificate Check
- **GET** `/api/v1/ssl?domain=example.com`
- Returns SSL certificate information including validity, issuer, expiration date, and key details

### HTTP/3 Support Check
- **GET** `/api/v1/http3?domain=example.com`
- Tests if the domain supports HTTP/3 protocol

### DNS Information
- **GET** `/api/v1/dns?domain=example.com`
- Returns DNS records including A, AAAA, CNAME, MX, TXT, and NS records

### IP Information
- **GET** `/api/v1/ip?ip=8.8.8.8` or `/api/v1/ip?domain=example.com`
- Returns IP address information and validation
- Accepts both IP addresses and domain names (resolves domain to IPs)

### Web Server Settings
- **GET** `/api/v1/web-settings?domain=example.com`
- Returns HTTP headers, server information, response time, and other web server details

### Comprehensive Check
- **GET** `/api/v1/comprehensive?domain=example.com`
- Returns all available information (SSL, HTTP/3, DNS, web settings) in a single request

### Blocklist Check
- **GET** `/api/v1/blocklist?domain=example.com`
- Returns whether a domain is blocked by various DNS servers and blocklist services

### HSTS Check
- **GET** `/api/v1/hsts?domain=example.com`
- Returns HTTP Strict Transport Security (HSTS) configuration including max-age, includeSubDomains, and preload directives

## Example Usage

### Check SSL Certificate
```bash
curl "http://localhost:8080/api/v1/ssl?domain=google.com"
```

### Check HTTP/3 Support
```bash
curl "http://localhost:8080/api/v1/http3?domain=cloudflare.com"
```

### Get DNS Information
```bash
curl "http://localhost:8080/api/v1/dns?domain=github.com"
```

### Check Web Server Settings
```bash
curl "http://localhost:8080/api/v1/web-settings?domain=example.com"
```

### Comprehensive Check
```bash
curl "http://localhost:8080/api/v1/comprehensive?domain=google.com"
```

### Check Blocklist Status
```bash
curl "http://localhost:8080/api/v1/blocklist?domain=example.com"
```

### Check HSTS Configuration
```bash
curl "http://localhost:8080/api/v1/hsts?domain=google.com"
```

## Response Format

All API responses follow this format:

```json
{
  "success": true,
  "data": {
    // Response data here
  },
  "error": "Error message if any",
  "message": "Additional message if any"
}
```

## Error Handling

The API includes comprehensive error handling:
- Invalid domain names
- Network timeouts
- SSL connection failures
- DNS resolution errors
- HTTP request failures

## Dependencies

- [Gin](https://github.com/gin-gonic/gin) - Web framework
- [quic-go](https://github.com/quic-go/quic-go) - HTTP/3 support
- Standard Go libraries for networking and crypto

## License

This project is open source and available under the MIT License.
