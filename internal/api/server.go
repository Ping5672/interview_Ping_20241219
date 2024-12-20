package api

import (
	"interview_Ping_20241219/internal/services"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	server := &Server{
		router: gin.Default(),
	}
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// Add CORS middleware
	s.router.Use(corsMiddleware())

	// Register routes
	services.RegisterPlayerRoutes(s.router)
	services.RegisterLevelRoutes(s.router)
	services.RegisterRoomRoutes(s.router)
	services.RegisterReservationRoutes(s.router)
	services.RegisterChallengeRoutes(s.router)
	services.RegisterLogRoutes(s.router)
	services.RegisterPaymentRoutes(s.router)
}

func (s *Server) Router() *gin.Engine {
	return s.router
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
