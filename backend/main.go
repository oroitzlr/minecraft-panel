package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oroitz-lago-ramos/minecraft-panel/internal/auth"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/joho/godotenv"
)

func main() {
	// Charger les variables d'environnement
	godotenv.Load()
	
	// 1. Connexion MongoDB
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("❌ MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ MongoDB ne répond pas: %v", err)
	}
	log.Println("✅ MongoDB connecté")

	db := client.Database("minecraft_panel")

	// 2. Services
	authService := auth.NewService(db)
	authService.EnsureAdminExists()

	// 3. Handlers
	authHandler := auth.NewHandler(authService)

	// 4. Gin + CORS
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 5. Routes publiques
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/logout", authHandler.Logout)
	}

	// 6. Routes protégées
	api := r.Group("/api")
	api.Use(auth.AuthMiddleware())
	{
		api.GET("/auth/me", authHandler.Me)
		api.GET("/server/status", func(c *gin.Context) {
			c.JSON(200, gin.H{"online": false, "players": 0})
		})
		api.POST("/server/start", auth.AdminOnly(), func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "démarrage en cours..."})
		})
		api.POST("/server/stop", auth.AdminOnly(), func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "arrêt en cours..."})
		})
	}

	// 7. Lancer
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 Serveur sur :%s", port)
	r.Run(":" + port)
}