package main

import (
	"fmt"
	"log"

	"funfillers/backend/config"
	"funfillers/backend/routes"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize PostgreSQL DB connection
	config.InitDB()

	r := routes.SetupRouter(cfg)

	fmt.Println("==================================================")
	fmt.Printf("🚀 FUNFILLERS Go Backend listening on port %s\n", cfg.Port)
	fmt.Println("   Health check: http://localhost:" + cfg.Port + "/api/health")
	fmt.Println("==================================================")

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
