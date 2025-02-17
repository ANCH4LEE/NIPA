package main

import (
	"back/internal/config"
	"back/internal/handlers"
	"back/internal/ticket"
	"context"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("Database Connection String:", cfg.GetConnectionString())

	db, err := ticket.NewPostgresDatabase(cfg.GetConnectionString())
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	repo := ticket.NewTicketRepo(db)
	if repo == nil {
		log.Fatal("TicketRepo initialization failed!")
	}

	h := handlers.NewTicketHandlers(repo)

	go func() {
		for {
			time.Sleep(10 * time.Second)
			if db == nil {
				log.Println("Database is nil, skipping Ping()")
				return
			}
			if err := db.Ping(); err != nil {
				log.Printf("Database connection lost: %v", err)
				if reconnErr := db.Reconnect(cfg.GetConnectionString()); reconnErr != nil {
					log.Printf("Failed to reconnect: %v", reconnErr)
				} else {
					log.Printf("Successfully reconnected to the database")
				}
			}
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	configCors := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:4000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(configCors))

	r.Use(TimeoutMiddleware(5 * time.Second))

	r.GET("/health", h.HealthCheck)

	v1 := r.Group("/api/v1")
	{
		tickets := v1.Group("/tickets")
		{
			tickets.GET("", h.GetTickets)
			tickets.POST("", h.CreateTicket)
			tickets.PUT("/:id", h.UpdateTicketStatus)
		}
	}

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Printf("Failed to run server: %v", err)
	}

}
