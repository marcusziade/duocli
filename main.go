package main

import (
	"duocli/cmd"
	"duocli/internal/database"
	"log"
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	
	// Execute CLI
	cmd.Execute()
}