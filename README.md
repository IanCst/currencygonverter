# Currency Tray App

A Windows system tray application that displays current USD and BRL exchange rates using the CurrencyAPI service.

## Features

- **System Tray Integration**: Runs quietly in the system tray
- **Real-time Rates**: Fetches live USD and BRL exchange rates
- **Automatic Updates**: Updates currency data daily at noon
- **Manual Refresh**: Click to refresh rates on demand
- **Clean Interface**: Simple menu with current rates and controls

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
- **Auto-Update**: The app automatically fetches new rates at noon daily
- **Exit**: Click "Quit" in the context menu

## API

This application uses the [CurrencyAPI](https://currencyapi.com/) service to fetch real-time exchange rates. The API key is embedded in the application.

## Files

- `main.go` - Application entry point
- `currency_api.go` - CurrencyAPI client implementation
- `tray.go` - System tray functionality and UI
- `test_api.go` - API testing utilities

## Requirements

- Windows OS (for system tray functionality)
- Internet connection (for fetching exchange rates)
- Go 1.25+ (for building from source)
Currency converter CLI Program 
