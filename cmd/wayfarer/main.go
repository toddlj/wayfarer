package main

import (
	"flag"
	"fmt"
	"google.golang.org/genproto/googleapis/type/latlng"
	"log/slog"
	"os"
	"time"
	"wayfarer/internal/config"
	"wayfarer/internal/googlemaps"
	"wayfarer/internal/scheduling"
	"wayfarer/internal/telegram"
)

func main() {
	// Configure logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	logger = logger.With("app", "wayfarer")
	slog.SetDefault(logger)

	// Load command-line arguments
	configFilePath := flag.String("config-file", "config.yaml", "Path of config file")
	flag.Parse()

	// Load environment variables
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramBotToken == "" {
		slog.Error("TELEGRAM_BOT_TOKEN environment variable must be set")
		os.Exit(1)
	}
	telegramApiBaseUrl := os.Getenv("TELEGRAM_API_BASE_URL")
	if telegramApiBaseUrl == "" {
		telegramApiBaseUrl = "https://api.telegram.org"
	}
	googleApiKey := os.Getenv("GOOGLE_API_KEY")
	if googleApiKey == "" {
		slog.Error("GOOGLE_API_KEY environment variable must be set")
		os.Exit(1)
	}
	googleApiBaseUrl := os.Getenv("GOOGLE_API_BASE_URL")

	// Load configuration
	cfg, err := config.LoadConfig(*configFilePath)
	if err != nil {
		slog.Error("Failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	// Initialize clients
	telegramClient := telegram.NewClient(telegramApiBaseUrl, telegramBotToken)
	mapsRoutingService, err := googlemaps.NewMapsRoutingService(googleApiBaseUrl, googleApiKey)
	if err != nil {
		slog.Error("Failed to initialize Google Maps client", slog.Any("error", err))
		os.Exit(1)
	}

	// Start scheduling tasks
	for _, rule := range cfg.Rules {
		err := scheduleRuleEvaluations(telegramClient, mapsRoutingService, rule)
		if err != nil {
			slog.Error("Failed to schedule rule", slog.Any("error", err), slog.Any("rule_id", rule.Id))
			os.Exit(1)
		}
	}

	// Keep the main thread alive
	select {}
}

func scheduleRuleEvaluations(telegramClient *telegram.Client, mapsRoutingService *googlemaps.MapsRoutingService, rule config.Rule) error {
	// Convert config as needed
	// Already validated in config.validate()
	schedules := make([]scheduling.Schedule, 0, len(rule.Times))
	for _, t := range rule.Times {
		weekday, _ := config.ParseWeekday(t.Day)
		timeOfDay, _ := time.Parse("15:04", t.Time)
		schedule := scheduling.Schedule{
			DayOfWeek: weekday,
			Hour:      timeOfDay.Hour(),
			Minute:    timeOfDay.Minute(),
		}
		schedules = append(schedules, schedule)
	}
	timezone, _ := time.LoadLocation(rule.Timezone)
	origin := &latlng.LatLng{Latitude: rule.Origin.Latitude, Longitude: rule.Origin.Longitude}
	destination := &latlng.LatLng{Latitude: rule.Destination.Latitude, Longitude: rule.Destination.Longitude}

	return scheduling.ScheduleFunction(schedules, timezone, func() {
		routeDuration, err := mapsRoutingService.FetchCurrentTransitTimeBetween(origin, destination)
		if err != nil {
			slog.Error("Failed to fetch transit time", slog.Any("error", err), slog.Any("rule_id", rule.Id))
			return
		}
		if routeDuration.Minutes() > float64(rule.TravelTime.NotificationThresholdMinutes) {
			slog.Info("Travel time exceeds threshold", slog.Any("rule_id", rule.Id), slog.Any("duration", routeDuration))
			message := fmt.Sprintf("Travel time between %s and %s is greater than %d minutes: currently scheduled to take %.0f minutes",
				rule.Origin.Name, rule.Destination.Name, rule.TravelTime.NotificationThresholdMinutes, routeDuration.Minutes())
			err := telegramClient.SendMessage(rule.User.TelegramUserID, message)
			if err != nil {
				slog.Error("Failed to send message", slog.Any("error", err), slog.Any("rule_id", rule.Id))
			}
		}
	})
}
