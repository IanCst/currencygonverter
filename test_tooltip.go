package main

import (
	"fmt"
	"time"
)

func testTooltipFormatting() {
	fmt.Println("Testing tooltip formatting...")
	
	// Test formatDuration function
	testCases := []time.Duration{
		2*time.Hour + 30*time.Minute, // 2h 30m
		45*time.Minute,               // 45m
		5*time.Minute,                // 5m
		0*time.Minute,                // 0m
		-30*time.Minute,              // negative
	}
	
	for _, duration := range testCases {
		formatted := formatDuration(duration)
		fmt.Printf("Duration: %v -> Formatted: '%s'\n", duration, formatted)
	}
	
	// Test full tooltip generation
	fmt.Println("\nTesting full tooltip generation...")
	lastUpdate := time.Now().Add(-1 * time.Hour) // 1 hour ago
	nextUpdate := lastUpdate.Add(3 * time.Hour)
	timeUntilNext := time.Until(nextUpdate)
	
	tooltipText := fmt.Sprintf("USD: $%.4f\nBRL: R$%.4f\nLast Updated: %s\nNext Update: in %s",
		1.0000, 5.4321, 
		lastUpdate.Format("15:04:05"),
		formatDuration(timeUntilNext))
		
	fmt.Println("Sample tooltip:")
	fmt.Println(tooltipText)
}

// Uncomment to run tooltip test:
// func main() {
// 	testTooltipFormatting()
// }