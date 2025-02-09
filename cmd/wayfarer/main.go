package main

import (
	"fmt"
	"google.golang.org/genproto/googleapis/type/latlng"
	"log"
	"log/slog"
	"os"
	"strconv"
	"wayfarer/internal/googlemaps"
	"wayfarer/internal/telegram"
)

func main() {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatIDStr := "399031758"

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		slog.Error("Invalid chat ID", slog.Any("error", err))
		return
	}

	bot := telegram.NewClient(botToken)

	err = bot.SendMessage(chatID, "Hello from Go module!")
	if err != nil {
		slog.Error("Failed to send message", slog.Any("error", err))
	}

	// Initialize the MapsRoutingService
	routingService, err := googlemaps.NewMapsRoutingService()
	if err != nil {
		log.Fatalf("Error initializing routing service: %v", err)
	}

	// Define origin and destination locations
	origin := &latlng.LatLng{Latitude: 40.748817, Longitude: -73.985428}      // Example: New York City (Lat, Lng)
	destination := &latlng.LatLng{Latitude: 40.712776, Longitude: -74.005974} // Example: Another location in NYC

	// Fetch the current transit time
	duration, err := routingService.FetchCurrentTransitTimeBetween(origin, destination)
	if err != nil {
		slog.Error("Failed to fetch transit time", slog.Any("error", err))
	}

	err = routingService.Close()
	if err != nil {
		slog.Error("Failed to close routing service", slog.Any("error", err))
	}

	// Output the transit time in seconds
	fmt.Printf("Transit time: %v seconds\n", duration.Seconds())

}
