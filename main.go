package main

import (
	"log"
	"os"
)

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Currency Tray Application...")

	// Initialize database connection (optional - app will work without DB)
	log.Println("Initializing database connection...")
	if err := InitDatabase(); err != nil {
		log.Printf("Warning: Database initialization failed: %v", err)
		log.Println("Application will continue without database functionality")
		log.Println("To enable database: check PostgreSQL is running and credentials in .env file")
	} else {
		log.Println("Database connection successful")
	}

	// Ensure database connection is closed on exit
	defer CloseDatabase()

	// Create and run the tray application
	app := NewTrayApp()
	
	// Handle any startup errors
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Application crashed: %v", r)
			os.Exit(1)
		}
	}()

	log.Println("Currency Tray App initialized. Starting system tray...")
	
	// This will block until the application is quit
	app.Run()
	
	log.Println("Currency Tray Application terminated.")
}
