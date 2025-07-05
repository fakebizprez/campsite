# Migration Guide: From AWS+Ruby to Open-Source+Go

This guide walks you through migrating from the original Campsite AWS-dependent Ruby application to the new Go-based API with local storage alternatives.

## Overview of Changes

### AWS Services Removed

| Original AWS Service | Replacement | Implementation |
|---------------------|-------------|----------------|
| **S3 (File Storage)** | Local Filesystem | Files stored in `./uploads` with direct HTTP endpoints |
| **ECS (Data Exports)** | Local Background Jobs | ZIP generation with Go's standard library |
| **Transcribe (Audio)** | Placeholder/Disabled | Can be replaced with Whisper or similar |

### Cloud Services Replaced

| Original Service | Replacement | Implementation |
|-----------------|-------------|----------------|
| **Pusher (Real-time)** | WebSocket Server | Native Go WebSocket implementation |
| **Imgix (CDN)** | Local File Serving | Direct file serving with basic content-type detection |

## Migration Steps

### Phase 1: Ruby API AWS Removal (Completed)

The Ruby API has been modified to remove AWS dependencies:

1. **Removed AWS SDK gems:**
   ```ruby
   # Removed from Gemfile:
   # gem "aws-sdk-s3"
   # gem "aws-sdk-ecs"  
   # gem "aws-sdk-transcribeservice"
   ```

2. **Added local storage replacement:**
   - `config/initializers/local_storage.rb` - Local file storage abstraction
   - `app/jobs/data_export_zip_job.rb` - Local ZIP creation job
   - `app/controllers/api/v1/uploads_controller.rb` - File upload endpoints

3. **Updated models:**
   - `DataExport` - Uses local background jobs instead of ECS
   - `DataExportResource` - Uses local storage exceptions
   - `PresignedPostFields` - Still works with local storage adapter

### Phase 2: Go API Foundation (Completed)

Created a new Go-based API in `go-api/` directory:

1. **Core functionality:**
   - Organizations, Projects, Posts CRUD operations
   - Local file storage with upload/download endpoints
   - WebSocket server for real-time communication
   - MySQL database with GORM ORM

2. **Development setup:**
   - Docker support for containerized deployment
   - Makefile for build automation
   - Test suite with SQLite in-memory database
   - Environment configuration

## Deployment Options

### Option 1: Keep Ruby API with Local Storage

If you prefer to stay with Ruby:

1. **Update Gemfile** (already done):
   ```bash
   cd api
   bundle install
   ```

2. **Configure local storage:**
   ```bash
   # Files will be stored in api/storage/uploads/
   mkdir -p storage/uploads
   ```

3. **Update credentials** to remove AWS settings:
   ```bash
   VISUAL="code --wait" bin/rails credentials:edit --environment development
   ```

4. **Run the server:**
   ```bash
   bundle exec rails server
   ```

### Option 2: Migrate to Go API

For the full Go migration:

1. **Setup Go API:**
   ```bash
   cd go-api
   make setup
   ```

2. **Configure database:**
   ```bash
   cp .env.example .env
   # Edit .env with your MySQL credentials
   ```

3. **Run the Go server:**
   ```bash
   make run
   # or for development with auto-reload:
   make dev
   ```

### Option 3: Hybrid Approach

Run both APIs during migration:

1. **Ruby API** on port 3001 (existing functionality)
2. **Go API** on port 8080 (new features)
3. **Nginx proxy** to route requests based on endpoints

## Data Migration

### Database Schema

The Go API uses simplified models. To migrate data:

1. **Export from Rails:**
   ```ruby
   # In Rails console
   organizations = Organization.all.map { |o| { name: o.name, slug: o.slug } }
   File.write('organizations.json', organizations.to_json)
   ```

2. **Import to Go API:**
   ```bash
   # Using the Go API endpoints
   curl -X POST http://localhost:8080/api/v1/organizations \
     -H "Content-Type: application/json" \
     -d '{"name":"Example Org","slug":"example"}'
   ```

### File Storage Migration

Move existing S3 files to local storage:

1. **Download from S3:**
   ```bash
   aws s3 sync s3://your-bucket ./local-files
   ```

2. **Copy to Go API uploads:**
   ```bash
   cp -r ./local-files/* ./go-api/uploads/
   ```

## Configuration Changes

### Environment Variables

**Ruby API** (.env or credentials):
```bash
# Remove these AWS-related variables:
# AWS_ACCESS_KEY_ID
# AWS_SECRET_ACCESS_KEY  
# AWS_S3_BUCKET

# Add local storage path (optional):
LOCAL_STORAGE_PATH=./storage/uploads
```

**Go API** (.env):
```bash
PORT=8080
DATABASE_URL=user:password@tcp(localhost:3306)/campsite_go
UPLOADS_DIR=./uploads
MAX_UPLOAD_SIZE=100MB
```

### Docker Deployment

**Ruby API** (existing):
```dockerfile
# Use existing Dockerfile, but remove AWS credentials
```

**Go API** (new):
```bash
cd go-api
make docker-build
make docker-run
```

### Nginx Configuration

For production deployment with both APIs:

```nginx
upstream ruby_api {
    server 127.0.0.1:3001;
}

upstream go_api {
    server 127.0.0.1:8080;
}

server {
    listen 80;
    
    # Route to Go API for new endpoints
    location /api/v2 {
        proxy_pass http://go_api;
    }
    
    # WebSocket support
    location /ws {
        proxy_pass http://go_api;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
    
    # Route to Ruby API for existing endpoints
    location /api/v1 {
        proxy_pass http://ruby_api;
    }
    
    # File uploads to Go API
    location /uploads {
        proxy_pass http://go_api;
    }
}
```

## Feature Compatibility

### Implemented in Go API

- ✅ Organizations CRUD
- ✅ Projects CRUD  
- ✅ Posts CRUD
- ✅ File uploads/downloads
- ✅ WebSocket real-time messaging
- ✅ Health checks
- ✅ CORS support

### Not Yet Implemented in Go API

- ❌ User authentication
- ❌ Permissions/authorization
- ❌ Advanced integrations (Slack, Linear, etc.)
- ❌ Call recording
- ❌ Email notifications
- ❌ Image processing
- ❌ Search functionality

### Ruby API Still Needed For

- User authentication and sessions
- Complex integrations
- Legacy endpoints
- Email functionality
- Advanced features

## Testing

### Ruby API Tests

```bash
cd api
bundle exec rails test
```

### Go API Tests

```bash
cd go-api
make test
```

### Integration Testing

Test that both APIs work together:

```bash
# Start Ruby API
cd api && bundle exec rails server -p 3001

# Start Go API  
cd go-api && make run

# Test endpoints
curl http://localhost:3001/api/v1/health  # Ruby
curl http://localhost:8080/health         # Go
```

## Performance Considerations

### File Storage

- **Local storage** is faster than S3 for small files
- **No network latency** for file operations
- **Simpler backup** with standard filesystem tools
- **Limited scalability** compared to object storage

### Database

- **Go API** uses connection pooling with GORM
- **Ruby API** uses Rails ActiveRecord
- **Consider** using the same database for both APIs during transition

### WebSockets vs Pusher

- **Native WebSockets** eliminate external dependency
- **Lower latency** for real-time features
- **More complex** client management
- **No built-in** channel authentication (implement as needed)

## Monitoring and Logging

### Go API Logging

```go
// Built-in logging
log.Printf("Server started on port %s", port)

// Add structured logging with logrus if needed
```

### Health Checks

- **Ruby API**: `/health` endpoint (if implemented)
- **Go API**: `/health` endpoint (implemented)

### Metrics

Consider adding:
- Request/response times
- Database connection status
- File storage usage
- WebSocket connection count

## Security Considerations

### File Uploads

- **Validate file types** and sizes
- **Scan for malware** before serving
- **Use secure filenames** to prevent path traversal

### WebSocket Security

- **Implement authentication** for WebSocket connections
- **Validate origins** in production
- **Rate limit** connections per IP

### Database

- **Use environment variables** for credentials
- **Enable SSL** connections in production
- **Regular backups** of local database

## Backup and Recovery

### File Storage Backup

```bash
# Daily backup
tar -czf backup-$(date +%Y%m%d).tar.gz uploads/

# Sync to remote storage
rsync -av uploads/ user@backup-server:/backups/campsite/
```

### Database Backup

```bash
# MySQL backup
mysqldump campsite_go > backup-$(date +%Y%m%d).sql

# Automated backup script
0 2 * * * /path/to/backup-script.sh
```

## Troubleshooting

### Common Issues

1. **Port conflicts**: Ensure Ruby API and Go API use different ports
2. **Database connections**: Check MySQL credentials and network access
3. **File permissions**: Ensure uploads directory is writable
4. **CORS errors**: Configure allowed origins for API access

### Debug Commands

```bash
# Check Go API health
curl http://localhost:8080/health

# Test file upload
curl -X POST -F "file=@test.txt" http://localhost:8080/api/v1/uploads

# WebSocket test (using wscat)
wscat -c ws://localhost:8080/ws
```

## Next Steps

1. **Phase 3**: Implement user authentication in Go API
2. **Phase 4**: Add permission system and authorization
3. **Phase 5**: Migrate remaining features from Ruby API
4. **Phase 6**: Deprecate Ruby API and fully transition to Go

## Support

For migration assistance:
- Check the Go API README: `go-api/README.md`
- Review test files for API usage examples
- Check Docker configurations for deployment patterns