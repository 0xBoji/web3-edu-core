package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/0xBoji/web3-edu-core/config"
	"github.com/0xBoji/web3-edu-core/internal/api/middleware"
	"github.com/0xBoji/web3-edu-core/internal/api/v1/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/0xBoji/web3-edu-core/docs" // This is required for swagger
)

// Server represents the API server
type Server struct {
	router *gin.Engine
}

// NewServer creates a new server
func NewServer() *Server {
	// Set gin mode
	gin.SetMode(config.ServerSetting.RunMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Apply global middleware
	router.Use(middleware.CorsMiddleware())

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup v1 API routes
	routes.RegisterRoutes(router)

	return &Server{
		router: router,
	}
}

// Run starts the server
func (s *Server) Run() {
	addr := fmt.Sprintf("%s:%d", config.ServerSetting.Host, config.ServerSetting.HttpPort)
	server := &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  config.ServerSetting.ReadTimeout,
		WriteTimeout: config.ServerSetting.WriteTimeout,
	}

	log.Printf("Server is running on %s", addr)
	server.ListenAndServe()
}
