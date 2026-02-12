package http

import (
	"net/http"
	"strings"

	"github.com/JAROBOTAI/jaro/internal/core/ports"
	"github.com/gin-gonic/gin"
)

// Server is the HTTP adapter that exposes the orchestrator through REST API.
// It follows Hexagonal Architecture by depending only on the Orchestrator port interface.
// This is a Primary Adapter (driving side) that receives external requests.
type Server struct {
	orchestrator ports.Orchestrator
}

// NewServer creates a new HTTP server with the given orchestrator.
// Purpose: Factory function for creating the HTTP API adapter with dependency injection.
// Inputs:
//   - orch: Implementation of the Orchestrator port for handling business logic
// Outputs:
//   - *Server: Initialized HTTP server ready to handle requests
func NewServer(orch ports.Orchestrator) *Server {
	return &Server{
		orchestrator: orch,
	}
}

// Run starts the HTTP server on the specified address.
// Purpose: Configures routes and starts the Gin HTTP server.
//          This method blocks until the server is stopped or encounters an error.
// Inputs:
//   - addr: Network address to listen on (e.g., ":8080", "0.0.0.0:3000")
// Outputs:
//   - error: Returns error if server fails to start or encounters fatal error
func (s *Server) Run(addr string) error {
	// Create Gin router with default middleware (logger, recovery)
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", s.healthCheckHandler)

	// Task management endpoints
	router.POST("/tasks", s.createTaskHandler)
	router.GET("/tasks/:id", s.getTaskStatusHandler)

	// Start server
	return router.Run(addr)
}

// healthCheckHandler handles health check requests.
// Purpose: Provides a simple endpoint for monitoring and load balancers.
// Inputs:
//   - c: Gin context containing request and response information
// Outputs: None (responds with JSON status)
func (s *Server) healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "jaro-orchestrator",
		"version": "v1.0",
	})
}

// CreateTaskRequest represents the expected JSON payload for creating a task.
type CreateTaskRequest struct {
	Input  string `json:"input" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
}

// createTaskHandler handles POST /tasks requests to create new tasks.
// Purpose: Receives user input, creates a task via orchestrator, and returns task details.
//          This is the primary entry point for submitting work to the JARO system.
// Inputs:
//   - c: Gin context containing request body with Input and UserID
// Outputs: JSON response with task_id and status (201 Created) or error (400/500)
func (s *Server) createTaskHandler(c *gin.Context) {
	var req CreateTaskRequest

	// Parse and validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Call orchestrator to create task
	task, err := s.orchestrator.StartTask(c.Request.Context(), req.Input, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create task",
			"details": err.Error(),
		})
		return
	}

	// Return successful response
	c.JSON(http.StatusCreated, gin.H{
		"task_id": task.ID,
		"status": task.Status,
		"created_at": task.CreatedAt,
		"user_id": task.UserID,
		"input": task.Input,
	})
}

// getTaskStatusHandler handles GET /tasks/:id requests to retrieve task status.
// Purpose: Allows clients to query the current state and progress of a task.
//          Returns the complete task object including status, artifacts, and metadata.
// Inputs:
//   - c: Gin context with task ID in URL parameter (:id)
// Outputs: JSON response with full task object (200 OK) or error (404/500)
func (s *Server) getTaskStatusHandler(c *gin.Context) {
	// Extract task ID from URL parameter
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "task_id is required",
		})
		return
	}

	// Call orchestrator to get task status
	task, err := s.orchestrator.GetTaskStatus(c.Request.Context(), taskID)
	if err != nil {
		// Check if it's a "not found" error by checking if error contains "not found"
		errMsg := err.Error()
		if strings.Contains(errMsg, "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "task not found",
				"task_id": taskID,
			})
			return
		}

		// Other errors are internal server errors
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get task status",
			"details": err.Error(),
		})
		return
	}

	// Return full task object
	c.JSON(http.StatusOK, task)
}
