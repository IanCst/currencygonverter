package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/getlantern/systray"
)

type TrayApp struct {
	currencyData *CurrencyData
	mutex        sync.RWMutex
	updateChan   chan bool
	lastUpdate   time.Time
	
	// Menu items
	refreshItem   *systray.MenuItem
	statusItem    *systray.MenuItem
	quitItem      *systray.MenuItem
}

func NewTrayApp() *TrayApp {
	return &TrayApp{
		updateChan: make(chan bool, 1),
	}
}

func (app *TrayApp) onReady() {
	// Set the icon and tooltip
	systray.SetTitle("Currency Tracker")
	systray.SetTooltip("USD/BRL Currency Tracker - Updates every 3 hours")
	
	// Try to load the custom exchange icon, fallback to default if it fails
	icon := getDefaultIcon() // Start with default
	
	if customIcon, err := ioutil.ReadFile("icon/exchange.png"); err == nil {
		// Successfully loaded custom icon
		icon = customIcon
		log.Println("Using custom exchange icon")
	} else {
		log.Printf("Using default icon (custom icon load failed: %v)", err)
	}
	
	// Set the icon with error handling
	if len(icon) > 0 {
		systray.SetIcon(icon)
	}

	// Create menu items
	app.statusItem = systray.AddMenuItem("Loading...", "Current exchange rates (updates every 3 hours)")
	app.statusItem.Disable()
	
	systray.AddSeparator()
	
	app.refreshItem = systray.AddMenuItem("Refresh Now", "Fetch latest exchange rates immediately")
	
	systray.AddSeparator()
	
	app.quitItem = systray.AddMenuItem("Quit", "Exit the application")

	// Initial currency fetch
	go app.updateCurrency()

	// Handle menu clicks
	go app.handleMenuClicks()
}

func (app *TrayApp) onExit() {
	// Cleanup
	close(app.updateChan)
}

func (app *TrayApp) handleMenuClicks() {
	for {
		select {
		case <-app.refreshItem.ClickedCh:
			go app.updateCurrency()
			
		case <-app.quitItem.ClickedCh:
			systray.Quit()
			return
		}
	}
}

func (app *TrayApp) updateCurrency() {
	// Update menu to show loading
	app.statusItem.SetTitle("Loading...")
	systray.SetTooltip("Fetching currency data...")

	// Fetch new data
	data, err := fetchCurrencyRates()
	if err != nil {
		log.Printf("Error fetching currency data: %v", err)
		app.statusItem.SetTitle("Error loading data")
		systray.SetTooltip("Error fetching currency data")
		return
	}

	// Update stored data and last update time
	app.mutex.Lock()
	app.currencyData = data
	app.lastUpdate = time.Now()
	app.mutex.Unlock()

	// Update menu display
	app.updateDisplay()
	
	log.Println("Currency data updated successfully")
}

func (app *TrayApp) updateDisplay() {
	app.mutex.RLock()
	defer app.mutex.RUnlock()
	
	if app.currencyData == nil {
		return
	}

	// Update menu item text
	statusText := fmt.Sprintf("USD: $%.4f | BRL: R$%.4f", 
		app.currencyData.USD, app.currencyData.BRL)
	app.statusItem.SetTitle(statusText)

	// Calculate next update time (3 hours from last update)
	var tooltipText string
	if !app.lastUpdate.IsZero() {
		nextUpdate := app.lastUpdate.Add(3 * time.Hour)
		timeUntilNext := time.Until(nextUpdate)
		
		// Debug logging
		log.Printf("Debug: lastUpdate=%v, nextUpdate=%v, timeUntilNext=%v", 
			app.lastUpdate.Format("15:04:05"), nextUpdate.Format("15:04:05"), timeUntilNext)
		
		if timeUntilNext > 0 {
			// Next update is in the future
			durationStr := formatDuration(timeUntilNext)
			log.Printf("Debug: formatted duration = '%s'", durationStr)
			// Simplified tooltip format to avoid encoding issues
			tooltipText = fmt.Sprintf("USD: %.4f | BRL: %.4f\nLast: %s | Next: %s",
				app.currencyData.USD, app.currencyData.BRL, 
				app.lastUpdate.Format("15:04"),
				durationStr)
		} else {
			// Next update is overdue
			tooltipText = fmt.Sprintf("USD: %.4f | BRL: %.4f\nLast: %s | Next: due now",
				app.currencyData.USD, app.currencyData.BRL,
				app.lastUpdate.Format("15:04"))
		}
	} else {
		// No last update time available
		tooltipText = fmt.Sprintf("USD: %.4f | BRL: %.4f\nLast: just now | Next: 3 hours",
			app.currencyData.USD, app.currencyData.BRL)
	}
	
	// Debug: log the exact tooltip text being set
	log.Printf("Debug: Setting tooltip to: '%s'", tooltipText)
	log.Printf("Debug: Tooltip length: %d bytes", len(tooltipText))
	
	// Test with hardcoded simple text first
	testTooltip := "USD: 1.0000 | BRL: 5.4321\nLast: 21:30 | Next: 2 hours"
	log.Printf("Debug: Testing with hardcoded tooltip: '%s'", testTooltip)
	systray.SetTooltip(testTooltip)
}

func (app *TrayApp) startPeriodicScheduler() {
	go func() {
		// Create a ticker for every 3 hours
		ticker := time.NewTicker(3 * time.Hour)
		defer ticker.Stop()
		
		log.Println("Starting periodic scheduler: updating every 3 hours")
		
		for {
			select {
			case <-ticker.C:
				log.Println("Periodic update triggered (3 hours interval)")
				go app.updateCurrency()
				
			case <-app.updateChan:
				log.Println("Periodic scheduler stopped")
				return
			}
		}
	}()
}

func (app *TrayApp) Run() {
	// Start the periodic scheduler (every 3 hours)
	app.startPeriodicScheduler()
	
	// Run the system tray
	systray.Run(app.onReady, app.onExit)
}

// formatDuration formats a duration into a human-readable string
func formatDuration(d time.Duration) string {
	// Handle negative durations
	if d < 0 {
		return "overdue"
	}
	
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	
	// Handle very short durations
	if hours == 0 && minutes == 0 {
		return "now"
	}
	
	// Use simple format with only basic ASCII
	if hours > 0 {
		if minutes > 0 {
			return fmt.Sprintf("%d hours %d min", hours, minutes)
		}
		return fmt.Sprintf("%d hours", hours)
	}
	return fmt.Sprintf("%d min", minutes)
}

// getDefaultIcon returns a simple default icon as byte array
func getDefaultIcon() []byte {
	// A simple 16x16 PNG icon with currency symbol
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x10,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x91, 0x68, 0x36, 0x00, 0x00, 0x00,
		0x1F, 0x49, 0x44, 0x41, 0x54, 0x28, 0x53, 0x63, 0x60, 0x18, 0x05, 0xA3,
		0x60, 0x14, 0x8C, 0x02, 0x08, 0x00, 0x00, 0x00, 0x04, 0x10, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
	}
}
