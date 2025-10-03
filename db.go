package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

var db *sql.DB

// getDefaultConfig returns default database configuration
// You can override these with environment variables
func getDefaultConfig() *DatabaseConfig {
	// Try to load .env file if it exists
	_ = godotenv.Load()
	
	return &DatabaseConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", ""),
		DBName:   getEnvOrDefault("DB_NAME", "currency_tracker"),
		SSLMode:  getEnvOrDefault("DB_SSL_MODE", "disable"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// InitDatabase initializes the database connection and creates the table
func InitDatabase() error {
	config := getDefaultConfig()
	
	log.Printf("Attempting to connect to PostgreSQL: %s@%s:%s/%s", 
		config.User, config.Host, config.Port, config.DBName)
	
	// Build connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	// Open database connection
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database (check your credentials): %v", err)
	}

	// Create the table if it doesn't exist
	if err = createTable(); err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	log.Printf("Database connection initialized successfully with user: %s", config.User)
	return nil
}

// createTable creates the currency_rates table if it doesn't exist
func createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS currency_rates (
		id SERIAL PRIMARY KEY,
		request_time TIMESTAMP WITH TIME ZONE NOT NULL,
		brl_to_usd_rate DECIMAL(10, 6) NOT NULL,
		usd_value DECIMAL(10, 6) NOT NULL DEFAULT 1.000000,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_currency_rates_request_time ON currency_rates(request_time);
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	log.Println("Currency rates table created/verified successfully")
	return nil
}

// SaveCurrencyRate saves a currency rate record to the database
func SaveCurrencyRate(requestTime time.Time, brlRate float64, usdRate float64) error {
	if db == nil {
		return fmt.Errorf("database connection not initialized")
	}

	query := `
	INSERT INTO currency_rates (request_time, brl_to_usd_rate, usd_value)
	VALUES ($1, $2, $3)
	`

	_, err := db.Exec(query, requestTime, brlRate, usdRate)
	if err != nil {
		return fmt.Errorf("failed to save currency rate: %v", err)
	}

	log.Printf("Saved currency rate: BRL=%.6f, USD=%.6f at %v", brlRate, usdRate, requestTime)
	return nil
}

// GetRecentRates retrieves the most recent currency rates from the database
func GetRecentRates(limit int) ([]CurrencyRateRecord, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	query := `
	SELECT id, request_time, brl_to_usd_rate, usd_value, created_at
	FROM currency_rates
	ORDER BY request_time DESC
	LIMIT $1
	`

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent rates: %v", err)
	}
	defer rows.Close()

	var rates []CurrencyRateRecord
	for rows.Next() {
		var rate CurrencyRateRecord
		err := rows.Scan(&rate.ID, &rate.RequestTime, &rate.BRLRate, &rate.USDRate, &rate.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		rates = append(rates, rate)
	}

	return rates, nil
}

// CurrencyRateRecord represents a record from the currency_rates table
type CurrencyRateRecord struct {
	ID          int
	RequestTime time.Time
	BRLRate     float64
	USDRate     float64
	CreatedAt   time.Time
}

// CloseDatabase closes the database connection
func CloseDatabase() {
	if db != nil {
		db.Close()
		log.Println("Database connection closed")
	}
}