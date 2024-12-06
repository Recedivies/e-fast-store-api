package api

import (
	"fmt"

	"github.com/Roixys/e-fast-store-api/config"
	"github.com/Roixys/e-fast-store-api/exception"
	"github.com/Roixys/e-fast-store-api/token"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Server serves HTTP requests.
type Server struct {
	DB         *gorm.DB
	config     config.Config
	router     *gin.Engine
	tokenMaker token.Maker
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(configuration config.Config, DB *gorm.DB) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(configuration.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     configuration,
		DB:         DB,
		tokenMaker: tokenMaker,
	}

	if server.config.Environment != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.New()
	router.Use(CORSMiddleware())
	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/users/login", server.loginUser)
	router.POST("/users/register", server.createUser)

	authRoutes := router.Group("/")
	authRoutes.Use(authMiddleware(server.tokenMaker))
	{
		authRoutes.GET("/products", server.getListProduct)
		authRoutes.GET("/carts", server.getCartProduct)
		authRoutes.POST("/carts", server.createCartProduct)
		authRoutes.DELETE("/carts/:product_id", server.deleteCartProduct)
		authRoutes.POST("/payments", server.createPayment)
	}

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) start(address string) error {
	return server.router.Run(address)
}

func RunGinServer(configuration config.Config, DB *gorm.DB) {
	server, err := NewServer(configuration, DB)
	exception.FatalIfNeeded(err, "cannot create server")

	err = server.start(configuration.HTTPServerAddress)
	exception.FatalIfNeeded(err, "cannot start server")
}
