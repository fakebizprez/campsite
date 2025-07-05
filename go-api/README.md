# Campsite Go API

A Go-based API server that replaces the Ruby on Rails backend of Campsite with local storage alternatives instead of AWS services.

## Features

- **Local File Storage**: Replaces AWS S3 with local filesystem storage
- **WebSocket Support**: Replaces Pusher with native WebSocket implementation
- **RESTful API**: Core endpoints for organizations, projects, and posts
- **Database Abstraction**: Uses GORM for database operations
- **Docker Support**: Containerized deployment
- **Test Coverage**: Unit tests with SQLite in-memory database

## Quick Start

### Prerequisites

- Go 1.21 or later
- MySQL/MariaDB (or SQLite for testing)
- Make (optional, for using Makefile commands)

### Installation

1. **Clone and setup:**
   ```bash
   cd go-api
   make setup
   ```

2. **Configure environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your database configuration
   ```

3. **Run the server:**
   ```bash
   make run
   ```

The API will be available at `http://localhost:8080`

### Development

For development with auto-reload:
```bash
# Install air for auto-reload
go install github.com/cosmtrek/air@latest

# Start development server
make dev
```

### Testing

Run the test suite:
```bash
make test
```

### Docker Deployment

Build and run with Docker:
```bash
make docker-build
make docker-run
```

## API Endpoints

### Organizations
- `GET /api/v1/organizations` - List all organizations
- `POST /api/v1/organizations` - Create a new organization
- `GET /api/v1/organizations/:slug` - Get organization by slug

### Projects
- `GET /api/v1/organizations/:org_slug/projects` - List projects in organization
- `POST /api/v1/organizations/:org_slug/projects` - Create new project
- `GET /api/v1/organizations/:org_slug/projects/:slug` - Get project by slug

### Posts
- `GET /api/v1/organizations/:org_slug/projects/:project_slug/posts` - List posts in project
- `POST /api/v1/organizations/:org_slug/projects/:project_slug/posts` - Create new post
- `GET /api/v1/organizations/:org_slug/projects/:project_slug/posts/:id` - Get post by ID

### File Uploads
- `POST /api/v1/uploads` - Upload a file (replaces S3 presigned posts)
- `GET /api/v1/uploads/:key` - Serve uploaded file

### WebSocket
- `GET /ws` - WebSocket connection for real-time updates
- `POST /api/v1/broadcast` - Broadcast message to connected clients

### Health Check
- `GET /health` - Service health status

## Architecture Changes

### Replaced AWS Services

1. **AWS S3** → Local filesystem storage
   - Files stored in `./uploads` directory
   - Direct upload/download endpoints
   - Configurable storage path

2. **AWS ECS** → Local background jobs
   - Data exports processed locally
   - ZIP file generation using Go's standard library
   - Background processing with goroutines

3. **AWS Transcribe** → Placeholder implementation
   - Service disabled by default
   - Can be replaced with open-source alternatives like Whisper

### Replaced Cloud Services

1. **Pusher** → Native WebSocket implementation
   - Real-time message broadcasting
   - Channel-based communication
   - Client connection management

2. **Imgix CDN** → Local file serving
   - Direct file serving from local storage
   - Basic content-type detection
   - Can be extended with image processing

## Database Schema

The Go API uses the same basic entities as the original Ruby application:

- **Organizations**: Top-level tenant/workspace
- **Projects**: Containers for posts and content
- **Posts**: Main content items
- **Users/Members**: Authentication and permissions (to be implemented)

## Configuration

Environment variables (see `.env.example`):

```bash
# Server
PORT=8080

# Database
DATABASE_URL=user:password@tcp(localhost:3306)/campsite_go

# File Storage
UPLOADS_DIR=./uploads
MAX_UPLOAD_SIZE=100MB

# CORS
CORS_ALLOWED_ORIGINS=*
```

## Migration from Ruby API

This Go API provides a subset of the original Ruby API functionality:

1. **Core Features Implemented:**
   - Organization management
   - Project management  
   - Post management
   - File uploads
   - Real-time messaging

2. **Features Not Yet Implemented:**
   - User authentication
   - Permissions/authorization
   - Advanced integrations
   - Call recording
   - Email notifications

3. **AWS Services Replaced:**
   - S3 → Local file storage
   - ECS → Local processing
   - Transcribe → Placeholder

## Development Roadmap

- [ ] User authentication and JWT tokens
- [ ] Permission system with RBAC
- [ ] Image processing pipeline  
- [ ] Background job queue
- [ ] Email notifications
- [ ] API rate limiting
- [ ] Logging and metrics
- [ ] Integration tests
- [ ] Documentation generation

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run the test suite
6. Submit a pull request

## License

This project maintains the same license as the original Campsite codebase.