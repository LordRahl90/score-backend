package servers

import (
	"net/http"

	"sybo/domains/users"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Server server class defining the essential dependencies
type Server struct {
	db          *gorm.DB
	Router      *gin.Engine
	userService users.UserServicer
}

func New(db *gorm.DB) (*Server, error) {
	userSvc, err := users.New(db)
	if err != nil {
		return nil, err
	}
	server := &Server{
		db:          db,
		Router:      gin.Default(),
		userService: userSvc,
	}
	server.Router.Use(CORSMiddleware())
	server.userRoutes()
	return server, nil
}

func created(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusCreated, data)
}

func success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}

func serverError(ctx *gin.Context, code int, err error) {
	ctx.JSON(code, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}

// CORSMiddleware cors middleware for dev/testing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Access-Control-Allow-Headers, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
