package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Organization struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Slug string `json:"slug" gorm:"uniqueIndex"`
}

type Project struct {
	ID             uint         `json:"id" gorm:"primaryKey"`
	Name           string       `json:"name"`
	Slug           string       `json:"slug"`
	OrganizationID uint         `json:"organization_id"`
	Organization   Organization `json:"organization" gorm:"foreignKey:OrganizationID"`
}

type Post struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	ProjectID uint    `json:"project_id"`
	Project   Project `json:"project" gorm:"foreignKey:ProjectID"`
}

var db *gorm.DB

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Default to SQLite for easy local development
		dsn = "./campsite.db"
	}

	var err error
	
	// Choose driver based on DSN format
	if strings.Contains(dsn, "mysql") || strings.Contains(dsn, "tcp(") {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else {
		// Use SQLite for local development
		if strings.HasPrefix(dsn, "sqlite://") {
			dsn = strings.TrimPrefix(dsn, "sqlite://")
		}
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	}
	
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	db.AutoMigrate(&Organization{}, &Project{}, &Post{})
	log.Println("Database connected and migrated")
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	initDB()

	// Create WebSocket hub
	hub := newHub()
	go hub.run()

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware for development
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Setup WebSocket routes
	setupWebSocketRoutes(r, hub)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "campsite-go-api"})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Organizations
		v1.GET("/organizations", getOrganizations)
		v1.POST("/organizations", createOrganization)
		v1.GET("/organizations/by-slug/:slug", getOrganization)

		// Projects
		v1.GET("/organizations/:org_slug/projects", getProjects)
		v1.POST("/organizations/:org_slug/projects", createProject)
		v1.GET("/organizations/:org_slug/projects/by-slug/:slug", getProject)
		
		// Posts
		v1.GET("/organizations/:org_slug/projects/:project_slug/posts", getPosts)
		v1.POST("/organizations/:org_slug/projects/:project_slug/posts", createPost)
		v1.GET("/posts/:id", getPost)

		// File uploads (replacing S3)
		v1.POST("/uploads", uploadFile)
		v1.GET("/uploads/:key", serveFile)
		
		// Integration stubs for external services
		setupIntegrationRoutes(v1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Campsite Go API on port %s", port)
	log.Fatal(r.Run(":" + port))
}

// Organization handlers
func getOrganizations(c *gin.Context) {
	var organizations []Organization
	if err := db.Find(&organizations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, organizations)
}

func createOrganization(c *gin.Context) {
	var org Organization
	if err := c.ShouldBindJSON(&org); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&org).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, org)
}

func getOrganization(c *gin.Context) {
	slug := c.Param("slug")
	var org Organization
	if err := db.Where("slug = ?", slug).First(&org).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}
	c.JSON(http.StatusOK, org)
}

// Project handlers
func getProjects(c *gin.Context) {
	orgSlug := c.Param("org_slug")
	var projects []Project
	
	if err := db.Joins("Organization").Where("Organization.slug = ?", orgSlug).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, projects)
}

func createProject(c *gin.Context) {
	orgSlug := c.Param("org_slug")
	
	// Find organization
	var org Organization
	if err := db.Where("slug = ?", orgSlug).First(&org).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	var project Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project.OrganizationID = org.ID
	if err := db.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, project)
}

func getProject(c *gin.Context) {
	orgSlug := c.Param("org_slug")
	projectSlug := c.Param("slug")
	
	var project Project
	if err := db.Joins("Organization").Where("Organization.slug = ? AND Project.slug = ?", orgSlug, projectSlug).First(&project).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	
	c.JSON(http.StatusOK, project)
}

// Post handlers
func getPosts(c *gin.Context) {
	orgSlug := c.Param("org_slug")
	projectSlug := c.Param("project_slug")
	
	var posts []Post
	if err := db.Joins("Project").Joins("JOIN organizations ON projects.organization_id = organizations.id").
		Where("organizations.slug = ? AND Project.slug = ?", orgSlug, projectSlug).
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, posts)
}

func createPost(c *gin.Context) {
	orgSlug := c.Param("org_slug")
	projectSlug := c.Param("project_slug")
	
	// Find project
	var project Project
	if err := db.Joins("Organization").Where("Organization.slug = ? AND Project.slug = ?", orgSlug, projectSlug).First(&project).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.ProjectID = project.ID
	if err := db.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

func getPost(c *gin.Context) {
	id := c.Param("id")
	
	var post Post
	if err := db.Preload("Project.Organization").First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	
	c.JSON(http.StatusOK, post)
}

// File upload handlers (replacing S3)
func uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Create uploads directory if it doesn't exist
	os.MkdirAll("./uploads", 0755)

	// Save the file
	filename := generateFileName(file.Filename)
	filepath := "./uploads/" + filename
	
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"filename": filename,
		"url":      "/api/v1/uploads/" + filename,
		"size":     file.Size,
	})
}

func serveFile(c *gin.Context) {
	key := c.Param("key")
	filepath := "./uploads/" + key
	
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	
	c.File(filepath)
}

func generateFileName(originalName string) string {
	// Simple timestamp-based filename generation
	// In production, you'd want something more sophisticated
	return originalName // For now, keep original name (could cause conflicts)
}