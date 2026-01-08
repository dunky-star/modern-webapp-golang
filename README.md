# Modern Webapp Golang

> A production-ready, scalable web application built with Go, following industry best practices and clean architecture principles.

[![Go Version](https://img.shields.io/badge/Go-1.25.5-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## 🚀 Features

- **RESTful API** with JSON responses and health check endpoints
- **Template Rendering** using Go's `html/template` package with intelligent caching
- **Session Management** using industry-standard `alexedwards/scs` with secure cookie-based sessions
- **CSRF Protection** with cryptographically secure token generation and validation
- **Security Headers** middleware for XSS, clickjacking, and content-type protection
- **Request Logging** with rotating log files (size and age-based rotation)
- **Clean Architecture** with separation of concerns (cmd, pkg, internal)
- **Production-Ready** with configurable timeouts, error handling, and environment modes
- **Modular Design** with reusable packages and components

## 📋 Prerequisites

- Go 1.25.5 or higher
- Git

## 🛠️ Installation

```bash
# Clone the repository
git clone https://github.com/dunky-star/modern-webapp-golang.git
cd modern-webapp-golang

# Install dependencies
go mod download

# Build the application
go build -o bin/api ./cmd/api

# Run the application
./bin/api
```

## 🏃 Quick Start

```bash
# Run with default settings (port 3000)
go run ./cmd/api

# Run with custom port
go run ./cmd/api -port=8080

# Run in production mode
go run ./cmd/api -port=3000 -env=prod
```

## 📁 Project Structure

```
modern-webapp-golang/
├── cmd/
│   └── api/              # Application entry point (handlers, routes, middleware, CSRF)
├── pkg/
│   └── render/          # Reusable template rendering package
├── internal/
│   └── data/            # Internal application packages (template data structures)
├── web/                 # HTML templates
├── output/
│   └── logs/            # Rotating access logs
└── go.mod              # Go module definition
```

## 🌐 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Home page |
| `GET` | `/v1/health` | Health check with uptime and status |
| `GET` | `/v1/about` | About page |
| `GET` | `/favicon.ico` | Favicon handler |

### Health Check Response

```json
Version: 1.0.0
{
  "status": "available",
  "uptime": "2h15m30s",
  "timestamp": "2025-01-02T22:30:00Z"
}
```

## 🏗️ Architecture

This project follows the **Standard Go Project Layout**:

- **`cmd/`** - Main applications for this project
- **`pkg/`** - Library code that's ok to use by external applications
- **`internal/`** - Private application and library code
- **`web/`** - Web assets and templates
- **`output/`** - Generated output files (logs, etc.)

## 🔧 Configuration

The application supports the following command-line flags:

- `-port` - Server port (default: 3000)
- `-env` - Environment mode: `dev`, `stage`, or `prod` (default: `dev`)

### Environment Modes

- **`dev`** - Development mode: templates reload on every request, logs to console and file
- **`stage`** - Staging mode: template caching enabled, logs to file only
- **`prod`** - Production mode: template caching enabled, secure cookies, logs to file only

## 🛡️ Security Features

### CSRF Protection
- Cryptographically secure token generation (32-byte random tokens)
- HTTP-only cookies with SameSite=Strict
- Constant-time comparison to prevent timing attacks
- Supports both header (`X-CSRF-Token`) and form field (`csrf_token`) submission
- 12-hour token validity
- Automatic token injection into templates

### Session Management
- **Industry-standard implementation** using `alexedwards/scs/v2`
- **Cookie-based sessions** with secure, HTTP-only cookies
- **24-hour session lifetime** with automatic expiration
- **Environment-aware security**: Secure flag enabled in production (HTTPS only)
- **SameSite=Strict** protection against CSRF attacks
- **Stateless design** - session data stored client-side in encrypted cookies
- **Simple API**: `session.Put()`, `session.Get()`, `session.GetString()` for easy access

### Security Headers
- `X-Content-Type-Options: nosniff` - Prevents MIME type sniffing
- `X-Frame-Options: deny` - Prevents clickjacking attacks
- `X-XSS-Protection: 1; mode=block` - Enables XSS filtering
- `Referrer-Policy: strict-origin-when-cross-origin` - Controls referrer information

## 🔄 Middleware Stack

The application uses a layered middleware approach (applied in order):

1. **Security Headers** - Adds security headers to all responses
2. **Request Logging** - Logs all HTTP requests with method, path, status, and duration
3. **Session Management** - Loads and saves session data for each request
4. **CSRF Protection** - Validates CSRF tokens for non-safe HTTP methods
5. **CSRF Token Generation** - Generates and injects tokens for GET requests

## 📊 Logging

### Request Logging
- All HTTP requests are logged to `output/logs/access.log`
- Log format: `RemoteAddr Protocol Method Path StatusCode Duration`
- Rotating log files based on:
  - **Size**: Rotates when file exceeds 5MB
  - **Age**: Rotates when file is older than 2 weeks
- Rotated files are archived with timestamp: `access.log.YYYYMMDD-HHMMSS`
- In development mode, logs are also written to console

### Application Logging
- Structured logging with timestamps
- Separate loggers for application events and HTTP requests

## 🎨 Template System

- **Template Caching**: Templates are cached in production for performance
- **Development Mode**: Templates reload on every request for easy development
- **CSRF Token Injection**: Tokens are automatically injected into template context
- **Base Layout**: Shared base template with Bootstrap styling
- **HTML Escaping**: Automatic XSS protection via Go's template package

## 🛡️ Production Features

- **HTTP Timeouts**: 
  - Read timeout: 10 seconds
  - Write timeout: 30 seconds
  - Idle timeout: 1 minute
- **Error Handling**: Comprehensive error handling and logging
- **Environment-Aware**: Different behaviors for dev, stage, and prod environments

> **Note**: For detailed information on security features (CSRF, Sessions, Security Headers), logging, and templates, see their respective sections above.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**Dunky Star**

- GitHub: [@dunky-star](https://github.com/dunky-star)

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Feel free to check the [issues page](https://github.com/dunky-star/modern-webapp-golang/issues).

---

⭐ Star this repo if you find it helpful!
