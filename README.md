# Currency Tray App

A Windows system tray application that displays current USD and BRL exchange rates using the CurrencyAPI service.

## Features

- **System Tray Integration**: Runs quietly in the system tray
- **Real-time Rates**: Fetches live USD and BRL exchange rates
- **Automatic Updates**: Updates currency data every 3 hours automatically
- **Manual Refresh**: Click to refresh rates on demand
- **Clean Interface**: Simple menu with current rates and next update countdown
- **Custom Icon**: Uses your custom exchange.png icon from the icon folder
- **PostgreSQL Integration**: Automatically saves currency data with timestamps to database
- **Data Persistence**: Stores request time and BRL/USD exchange rates for historical tracking

## Installation

1. Make sure you have Go installed on your system
2. Clone or download this repository
3. Open terminal/PowerShell in the project directory
4. Run: `go build -o currency-tray.exe .`
5. Run the executable: `./currency-tray.exe`

## Usage

- **Start**: Double-click `currency-tray.exe` or run it from command line
- **View Rates**: Right-click the tray icon to see current exchange rates
- **Refresh**: Click "Refresh Now" to update rates immediately
- **Auto-Update**: The app automatically fetches new rates every 3 hours
- **Next Update**: Hover over the tray icon to see time until next update
- **Exit**: Click "Quit" in the context menu

## API

This application uses the [CurrencyAPI](https://currencyapi.com/) service to fetch real-time exchange rates. The API key is embedded in the application.

## Files

- `main.go` - Application entry point
- `currency_api.go` - CurrencyAPI client implementation
- `tray.go` - System tray functionality and UI
- `db.go` - PostgreSQL database operations
- `test_api.go` - API testing utilities
- `setup_database.sql` - Database setup script
- `.env.example` - Environment variables template
- `icon/exchange.png` - Custom tray icon

## Database Setup (Optional)

The app can optionally save currency data to a PostgreSQL database:

1. **Install PostgreSQL** (if not already installed)
2. **Create Database**: Run the provided SQL script:
   ```bash
   psql -U postgres -f setup_database.sql
   ```
3. **Configure Connection**: Set environment variables or use defaults:
   - `DB_HOST=localhost`
   - `DB_PORT=5432`
   - `DB_USER=postgres`
   - `DB_PASSWORD=password`
   - `DB_NAME=currency_tracker`

*Note: The app will work without a database connection and will simply log warnings if the database is unavailable.*

## Requirements

- Windows OS (for system tray functionality)
- Internet connection (for fetching exchange rates)
- Go 1.25+ (for building from source)
- PostgreSQL (optional, for data persistence)
Currency converter CLI Program 
