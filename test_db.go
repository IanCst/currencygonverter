package main

import (
	"log"
	"time"
)

func testDatabase() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Testing database connection...")
	
	// Test database initialization
	if err := InitDatabase(); err != nil {
		log.Printf("Database connection failed: %v", err)
		return
	}
	
	log.Println("Database connected successfully!")
	
	// Test saving a currency rate
	testTime := time.Now()
	testBRL := 5.4321
	testUSD := 1.0000
	
	log.Printf("Testing save operation: BRL=%.4f, USD=%.4f", testBRL, testUSD)
	if err := SaveCurrencyRate(testTime, testBRL, testUSD); err != nil {
		log.Printf("Save operation failed: %v", err)
	} else {
		log.Println("Save operation successful!")
	}
	
	// Test reading recent rates
	log.Println("Testing read operation...")
	rates, err := GetRecentRates(5)
	if err != nil {
		log.Printf("Read operation failed: %v", err)
	} else {
		log.Printf("Read operation successful! Found %d records", len(rates))
		for i, rate := range rates {
			log.Printf("  Record %d: BRL=%.6f, USD=%.6f, Time=%v", 
				i+1, rate.BRLRate, rate.USDRate, rate.RequestTime.Format("15:04:05"))
		}
	}
	
	CloseDatabase()
	log.Println("Database test completed.")
}

// Uncomment to run database test:
// func main() {
// 	testDatabase()
// }