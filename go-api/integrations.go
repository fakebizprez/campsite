package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Integration stub responses
type SlackResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ID      string `json:"id,omitempty"`
}

type LinearResponse struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

type EmailResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

type TranscriptionResponse struct {
	Transcript string  `json:"transcript"`
	Confidence float64 `json:"confidence"`
	Status     string  `json:"status"`
}

type SearchResponse struct {
	Results []interface{} `json:"results"`
	Total   int           `json:"total"`
	Message string        `json:"message"`
}

func setupIntegrationRoutes(r *gin.RouterGroup) {
	// Integration stubs - these replace external services
	integrations := r.Group("/integrations")
	{
		// Slack integration stub
		slack := integrations.Group("/slack")
		{
			slack.POST("/send", func(c *gin.Context) {
				var req struct {
					Message string `json:"message"`
					Channel string `json:"channel"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, SlackResponse{
					Success: true,
					Message: "Slack integration disabled - using stub",
					ID:      generateStubID(),
				})
			})

			slack.POST("/channels", func(c *gin.Context) {
				var req struct {
					Name string `json:"name"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"id":   generateStubID(),
					"name": req.Name,
				})
			})
		}

		// Linear integration stub
		linear := integrations.Group("/linear")
		{
			linear.POST("/issues", func(c *gin.Context) {
				var req struct {
					Title       string `json:"title"`
					Description string `json:"description"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, LinearResponse{
					ID:    generateStubID(),
					Title: req.Title,
					URL:   "https://linear.stub/issue/" + generateStubID(),
				})
			})

			linear.POST("/issues/:id/comments", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"synced": 1,
					"status": "success",
				})
			})
		}

		// Email service stub
		email := integrations.Group("/email")
		{
			email.POST("/send", func(c *gin.Context) {
				var req struct {
					To       string `json:"to"`
					Subject  string `json:"subject"`
					Body     string `json:"body"`
					Template string `json:"template,omitempty"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				// Log email in development
				if gin.Mode() == gin.DebugMode {
					println("\n" + "===========================================")
					println("EMAIL STUB")
					println("To:", req.To)
					println("Subject:", req.Subject)
					if req.Template != "" {
						println("Template:", req.Template)
					}
					println("-------------------------------------------")
					println(req.Body)
					println("===========================================\n")
				}

				c.JSON(http.StatusOK, EmailResponse{
					MessageID: "stub_" + generateStubID(),
					Status:    "sent",
				})
			})
		}

		// Transcription service stub
		transcription := integrations.Group("/transcription")
		{
			transcription.POST("/transcribe", func(c *gin.Context) {
				var req struct {
					AudioURL string `json:"audio_url"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, TranscriptionResponse{
					Transcript: "[Transcription service disabled - placeholder text]",
					Confidence: 0.0,
					Status:     "completed",
				})
			})

			transcription.GET("/vtt/:id", func(c *gin.Context) {
				vtt := `WEBVTT

00:00:00.000 --> 00:05:00.000
[Transcription service disabled - placeholder text]`

				c.Header("Content-Type", "text/vtt")
				c.String(http.StatusOK, vtt)
			})
		}

		// Search service stub
		search := integrations.Group("/search")
		{
			search.POST("/index", func(c *gin.Context) {
				var req struct {
					Type    string      `json:"type"`
					ID      string      `json:"id"`
					Content interface{} `json:"content"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"indexed": true,
					"type":    req.Type,
					"id":      req.ID,
				})
			})

			search.GET("/query", func(c *gin.Context) {
				query := c.Query("q")
				
				c.JSON(http.StatusOK, SearchResponse{
					Results: []interface{}{},
					Total:   0,
					Message: "Search service disabled - query was: " + query,
				})
			})

			search.DELETE("/:type/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"deleted": true,
					"type":    c.Param("type"),
					"id":      c.Param("id"),
				})
			})
		}

		// Figma integration stub
		figma := integrations.Group("/figma")
		{
			figma.GET("/files/:key", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"name":          "Stub Design File",
					"thumbnail_url": "/placeholder-thumbnail.png",
					"last_modified": time.Now().Format(time.RFC3339),
				})
			})

			figma.GET("/files/:key/nodes/:nodeId", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"image_url": "/placeholder-frame.png",
					"status":    "rendered",
				})
			})
		}

		// Cal.com integration stub
		calendar := integrations.Group("/calendar")
		{
			calendar.POST("/bookings", func(c *gin.Context) {
				var req struct {
					EventType string `json:"event_type"`
					DateTime  string `json:"datetime"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"id":     generateStubID(),
					"url":    "https://cal.stub/booking/" + generateStubID(),
					"status": "confirmed",
				})
			})
		}
	}
}

func generateStubID() string {
	// Simple stub ID generation
	return "stub_" + fmt.Sprintf("%d", time.Now().UnixNano()%1000000)
}