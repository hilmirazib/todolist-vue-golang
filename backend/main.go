package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"my-todolist/db"
	"my-todolist/middleware"
	"my-todolist/models"
	"my-todolist/routes"
)

func main() {
	_ = godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	db.Connect(&models.Todo{})

	r := gin.New()
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		sugar.Infow("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"ip", c.ClientIP(),
			"ua", c.Request.UserAgent(),
		)
	})

	origins := strings.Split(os.Getenv("CORS_ALLOW_ORIGINS"), ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	r.Use(cors.New(cors.Config{
		// AllowOrigins: []string{"*"}, // sesuaikan di prod
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(middleware.RateLimit())

	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	routes.Register(r)
	log.Printf("listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
