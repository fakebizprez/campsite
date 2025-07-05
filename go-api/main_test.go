package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	// Use in-memory SQLite for testing
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	// Auto-migrate the schema
	testDB.AutoMigrate(&Organization{}, &Project{}, &Post{})
	
	return testDB
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	
	// Set up test database
	db = setupTestDB()
	
	r := gin.New()
	
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "campsite-go-api"})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/organizations", getOrganizations)
		v1.POST("/organizations", createOrganization)
		v1.GET("/organizations/:slug", getOrganization)
	}
	
	return r
}

func TestHealthCheck(t *testing.T) {
	router := setupTestRouter()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "campsite-go-api", response["service"])
}

func TestCreateOrganization(t *testing.T) {
	router := setupTestRouter()
	
	org := Organization{
		Name: "Test Organization",
		Slug: "test-org",
	}
	
	jsonData, _ := json.Marshal(org)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/organizations", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response Organization
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Test Organization", response.Name)
	assert.Equal(t, "test-org", response.Slug)
	assert.NotZero(t, response.ID)
}

func TestGetOrganizations(t *testing.T) {
	router := setupTestRouter()
	
	// First, create an organization
	org := Organization{Name: "Test Org", Slug: "test"}
	db.Create(&org)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/organizations", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response []Organization
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, "Test Org", response[0].Name)
}

func TestGetOrganizationBySlug(t *testing.T) {
	router := setupTestRouter()
	
	// First, create an organization
	org := Organization{Name: "Test Org", Slug: "test-slug"}
	db.Create(&org)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/organizations/test-slug", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response Organization
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Test Org", response.Name)
	assert.Equal(t, "test-slug", response.Slug)
}

func TestGetNonExistentOrganization(t *testing.T) {
	router := setupTestRouter()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/organizations/non-existent", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Organization not found", response["error"])
}