package main

import (
	"log"
	"os"
)

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Currency Tray Application...")

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