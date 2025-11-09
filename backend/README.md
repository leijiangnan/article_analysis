# Article Analysis System - Backend

A Go-based backend service for article analysis management system with OpenAI integration.

## Features

- Article upload and storage (TXT format)
- Article categorization by author
- Full-text search functionality
- AI-powered article analysis using OpenAI GPT
- Analysis results storage and retrieval
- RESTful API design

## Tech Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: MySQL 8.0
- **AI Service**: OpenAI GPT API
- **Logging**: Zap
- **Configuration**: Viper

## Project Structure

```
backend/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── config/             # Configuration management
│   ├── handler/            # HTTP handlers
│   ├── middleware/          # Middleware components
│   ├── model/              # Data models
│   ├── repository/         # Data access layer
│   └── service/            # Business logic layer
├── pkg/
│   └── logger/             # Logging utilities
├── scripts/
│   └── init.sql            # Database initialization
├── config.yaml             # Configuration file
├── go.mod                   # Go module dependencies
└── README.md               # This file
```

## Quick Start

### Prerequisites

- Go 1.21 or higher
- MySQL 8.0 or higher
- OpenAI API key

### Installation

1. Clone the repository
2. Install dependencies:
```bash
cd backend
go mod download
```

3. Configure the application:
```bash
cp config.yaml.example config.yaml
# Edit config.yaml with your database and OpenAI settings
```

4. Initialize the database:
```bash
mysql -u root -p < scripts/init.sql
```

5. Run the application:
```bash
go run cmd/main.go
```

## Configuration

Edit `config.yaml` file:

```yaml
database:
  host: localhost
  port: 3306
  username: root
  password: your_password
  database: article_analysis

server:
  port: 8080
  mode: debug

openai:
  api_key: your_openai_api_key
  model: gpt-3.5-turbo
  max_tokens: 2000

log:
  level: info
```

## API Endpoints

### Article Management

- **POST** `/api/v1/articles/upload` - Upload article file
- **GET** `/api/v1/articles` - Get article list with pagination and search
- **GET** `/api/v1/articles/:id` - Get article details

### Article Analysis

- **POST** `/api/v1/articles/:id/analyze` - Submit article for AI analysis
- **GET** `/api/v1/articles/:id/analysis` - Get analysis results
- **GET** `/api/v1/analysis/status/:task_id` - Get analysis task status

## Development

### Running in Development Mode

```bash
go run cmd/main.go
```

### Building for Production

```bash
go build -o article-analysis cmd/main.go
```

## License

MIT License