package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"remy/internal/config"
	"remy/internal/handlers"
	pubpkg "remy/internal/infrastructure/publisher"
	"remy/internal/logging"
	"remy/internal/middlewares"
	"remy/internal/response"
	"remy/internal/services"
)

func main() {
	config.LoadEnv()

	cfg, err := config.LoadConfig()
	if err != nil {
		logging.Logger.Fatalf("failed to load configuration: %v", err)
	}

	logging.InitLogger(cfg)

	db, err := config.SetupDatabase(cfg)
	if err != nil {
		logging.Logger.Fatalf("failed to connect to database: %v", err)
	}

	router := newRouter()
	setupRoutes(router, db)

	server := NewServer(router, cfg)
	startServer(server)
}

func newRouter() *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		logging.Logger.WithFields(logrus.Fields{
			"path":    c.Request.URL.Path,
			"method":  c.Request.Method,
			"status":  c.Writer.Status(),
			"latency": duration.String(),
		}).Info("incoming request")
	})
	router.Use(middlewares.GlobalErrorHandler())
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, response.APIResponse{
			Errors: []*response.APIError{
				{
					Status: http.StatusNotFound,
					Code:   "not_found",
					Title:  "Not Found",
					Detail: "The requested resource was not found.",
				},
			},
		})
	})
	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, response.APIResponse{
			Errors: []*response.APIError{
				{
					Status: http.StatusMethodNotAllowed,
					Code:   "method_not_allowed",
					Title:  "Method Not Allowed",
					Detail: "The requested method is not allowed for this resource.",
				},
			},
		})
	})

	return router
}

func NewServer(router *gin.Engine, cfg *config.Config) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}
}

func startServer(server *http.Server) {
	logging.Logger.Infof("Starting server on %s", server.Addr)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logging.Logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logging.Logger.Fatalf("Failed to shutdown server: %v", err)
	}

	logging.Logger.Info("Server gracefully stopped")
}

func setupRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := router.Group("/api/v1")

	publisher := pubpkg.NewInMemoryPublisher()

	noteService := services.NewNoteService(db, publisher)
	noteHandler := handlers.NewNoteHandler(noteService)

	v1.POST("/notes", noteHandler.Create)
	v1.GET("/notes", noteHandler.List)
}
