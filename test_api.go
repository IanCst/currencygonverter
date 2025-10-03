package main

import (
	"fmt"
	"log"
)

func testAPI() {
	log.Println("Testing Currency API...")
	
	data, err := fetchCurrencyRates()
	if err != nil {
		log.Fatalf("API test failed: %v", err)
	}
	
	fmt.Println("API Test Results:")
	fmt.Println(data.String())
}

// Uncomment the following to run the test:
// func main() {
// 	testAPI()
// }