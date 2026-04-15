package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/SamuelMarcano15/boleta-backend/internal/repository"
	"github.com/SamuelMarcano15/boleta-backend/internal/routes"
)

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("ERROR: La variable DATABASE_URL no está definida")
	}
	
	repository.ConnectDB(dbURL)

	r := gin.Default()

	// Llamamos a nuestro configurador de rutas
	routes.SetupRouter(r)

	// Dejamos el ping para pruebas de salud (Healthcheck)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Servidor corriendo en http://localhost:%s", port)
	
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error crítico al arrancar el servidor: %v", err)
	}
}