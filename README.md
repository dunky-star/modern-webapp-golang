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

## 🔄 Middleware Stack

The application uses a layered middleware approach (applied in order):

1. **Security Headers** - Adds security headers (X-Content-Type-Options, X-Frame-Options, X-XSS-Protection, Referrer-Policy) to all responses
2. **Request Logging** - Logs all HTTP requests to `output/logs/access.log` with rotating files (5MB/2 weeks max)
3. **Session Management** - Cookie-based sessions using `alexedwards/scs/v2` with 24-hour lifetime and secure, HTTP-only cookies
4. **CSRF Protection** - Validates CSRF tokens (32-byte, constant-time comparison) for non-safe HTTP methods
5. **CSRF Token Generation** - Generates and injects tokens into templates for GET requests

**Additional Features:**
- **Template System**: Caching in production, auto-reload in dev, automatic CSRF token injection, HTML escaping
- **Logging**: Structured logging with timestamps, separate loggers for application events and HTTP requests

## 🛡️ Production Features

- **HTTP Timeouts**: 
  - Read timeout: 10 seconds
  - Write timeout: 30 seconds
  - Idle timeout: 1 minute
- **Error Handling**: Comprehensive error handling and logging
- **Environment-Aware**: Different behaviors for dev, stage, and prod environments

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**Dunky Star**

- GitHub: [@dunky-star](https://github.com/dunky-star)

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Feel free to check the [issues page](https://github.com/dunky-star/modern-webapp-golang/issues).

---

⭐ Star this repo if you find it helpful!
