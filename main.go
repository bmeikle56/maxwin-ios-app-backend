package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"maxwin/config"
	"maxwin/database"
	"maxwin/handlers"
	"maxwin/middleware"
	"maxwin/mock"
	"maxwin/services"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	if err := config.LoadEnv(); err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	if config.UseMockDB() {
		log.Printf("using mock store (APP_ENV=%s)", os.Getenv("APP_ENV"))
	} else {
		if err := database.ConnectDB(); err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		defer func() {
			if err := database.Close(); err != nil {
				log.Printf("error closing DB: %v", err)
			}
		}()
		log.Printf("connected to postgres (APP_ENV=%s)", os.Getenv("APP_ENV"))
	}

	store := mock.NewStore()
	authHandlers := handlers.NewAuthHandlers(services.NewAuthService(store))
	sessionHandlers := handlers.NewSessionHandlers(services.NewSessionService(store))
	earningsHandlers := handlers.NewEarningsHandlers(services.NewEarningsService(store))

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     config.CORSOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"appEnv": os.Getenv("APP_ENV"),
			"mockDB": config.UseMockDB(),
		})
	})

	// Public auth endpoints — no JWT required.
	r.POST("/auth/signin", authHandlers.SignIn)
	r.POST("/auth/password-reset", authHandlers.PasswordReset)

	// Protected routes — require a per-user Bearer JWT.
	api := r.Group("/")
	api.Use(middleware.AuthMiddleware())
	{
		api.POST("/auth/signout", authHandlers.SignOut)
		api.DELETE("/auth/account", authHandlers.DeleteAccount)

		api.GET("/sessions", sessionHandlers.List)
		api.GET("/sessions/:id", sessionHandlers.Get)
		api.POST("/sessions", sessionHandlers.Create)
		api.PUT("/sessions/:id", sessionHandlers.Update)
		api.DELETE("/sessions/:id", sessionHandlers.Delete)

		api.GET("/earnings", earningsHandlers.Fetch)
	}

	port := config.Port()
	log.Printf("maxwin backend listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
