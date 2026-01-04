# Modern Webapp Golang

> A production-ready, scalable web application built with Go, following industry best practices and clean architecture principles.

[![Go Version](https://img.shields.io/badge/Go-1.25.5-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## 🚀 Features

- **RESTful API** with JSON responses and health check endpoints
- **Template Rendering** using Go's `html/template` package
- **Clean Architecture** with separation of concerns (cmd, pkg, internal)
- **Production-Ready** with configurable timeouts and error handling
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
│   └── api/              # Application entry point
│       ├── main.go       # Server configuration and startup
│       ├── handlers.go  # HTTP request handlers
│       └── routes.go    # Route definitions
├── pkg/
│   └── render/          # Reusable template rendering package
│       └── render.go
├── internal/            # Internal application packages
├── web/                 # HTML templates
│   ├── home.page.tmpl
│   └── about.page.tmpl
└── go.mod              # Go module definition
```

## 🌐 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Home page |
| `GET` | `/v1/health` | Health check with uptime and status |
| `GET` | `/v1/about` | About page |

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

## 🔧 Configuration

The application supports the following command-line flags:

- `-port` - Server port (default: 3000)
- `-env` - Environment mode: `dev`, `stage`, or `prod` (default: `dev`)

## 🛡️ Production Features

- **HTTP Timeouts**: Configurable read, write, and idle timeouts
- **Structured Logging**: Built-in logging with timestamps
- **Error Handling**: Comprehensive error handling and logging
- **Template Safety**: HTML escaping for XSS protection

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**Dunky Star**

- GitHub: [@dunky-star](https://github.com/dunky-star)

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Feel free to check the [issues page](https://github.com/dunky-star/modern-webapp-golang/issues).

---

⭐ Star this repo if you find it helpful!
