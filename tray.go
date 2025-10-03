package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/getlantern/systray"
)

type TrayApp struct {
	currencyData *CurrencyData
	mutex        sync.RWMutex
	updateChan   chan bool
	
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
	systray.SetTooltip("USD/BRL Currency Tracker")
	
	// Create a simple icon (you can replace this with a proper icon file)
	icon := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x10,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x91, 0x68, 0x36, 0x00, 0x00, 0x00,
		0x1F, 0x49, 0x44, 0x41, 0x54, 0x28, 0x53, 0x63, 0x60, 0x18, 0x05, 0xA3,
		0x60, 0x14, 0x8C, 0x02, 0x08, 0x00, 0x00, 0x00, 0x04, 0x10, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
	}
	systray.SetIcon(icon)

	// Create menu items
	app.statusItem = systray.AddMenuItem("Loading...", "Current exchange rates")
	app.statusItem.Disable()
	
	systray.AddSeparator()
	
	app.refreshItem = systray.AddMenuItem("Refresh Now", "Fetch latest exchange rates")
	
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

	// Update stored data
	app.mutex.Lock()
	app.currencyData = data
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

	// Update tooltip with more detailed info
	tooltipText := fmt.Sprintf("USD: $%.4f\nBRL: R$%.4f\nLast Updated: %s",
		app.currencyData.USD, app.currencyData.BRL, app.currencyData.LastUpdated)
	systray.SetTooltip(tooltipText)
}

func (app *TrayApp) startNoonScheduler() {
	go func() {
		for {
			now := time.Now()
			
			// Calculate next noon (12:00 PM)
			nextNoon := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
			
			// If it's already past noon today, schedule for tomorrow
			if now.After(nextNoon) {
				nextNoon = nextNoon.Add(24 * time.Hour)
			}
			
			// Calculate duration until next noon
			duration := nextNoon.Sub(now)
			
			log.Printf("Scheduling next update for: %v (in %v)", nextNoon, duration)
			
			// Wait until noon
			timer := time.NewTimer(duration)
			select {
			case <-timer.C:
				log.Println("Noon update triggered")
				go app.updateCurrency()
			case <-app.updateChan:
				timer.Stop()
				return
			}
		}
	}()
}

func (app *TrayApp) Run() {
	// Start the noon scheduler
	app.startNoonScheduler()
	
	// Run the system tray
	systray.Run(app.onReady, app.onExit)
}