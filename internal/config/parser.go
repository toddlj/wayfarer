package config

import (
	"errors"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

var ErrInvalidDay = errors.New("invalid day of the week")

var ErrInvalidTimeFormat = errors.New("invalid time format")

func (cfg *Config) Validate() error {
	for _, rule := range cfg.Rules {
		// Check ID
		if rule.Id <= 0 {
			return errors.New("id must be greater than 0")
		}

		// Check if Origin and Destination are defined
		if rule.Origin.Name == "" || rule.Destination.Name == "" {
			return errors.New("origin and destination must have a name")
		}

		// Check if User is defined
		if rule.User.TelegramUserID == 0 {
			return errors.New("user must have a Telegram user ID")
		}

		// Ensure TravelTime is specified
		if rule.TravelTime.NotificationThresholdMinutes <= 0 {
			return errors.New("notification_threshold_minutes must be greater than 0")
		}

		// Ensure Times are not empty
		if len(rule.Times) == 0 {
			return errors.New("at least one time must be specified")
		}

		// Validate each time schedule
		for _, t := range rule.Times {
			// Validate day
			if _, err := ParseWeekday(t.Day); err != nil {
				return err
			}

			// Validate time format
			if _, err := time.Parse("15:04", t.Time); err != nil {
				return ErrInvalidTimeFormat
			}
		}

		// Validate timezone
		if _, err := time.LoadLocation(rule.Timezone); err != nil {
			return err
		}
	}

	// All checks passed
	return nil
}

func ParseWeekday(day string) (time.Weekday, error) {
	dayMap := map[string]time.Weekday{
		"SUNDAY":    time.Sunday,
		"MONDAY":    time.Monday,
		"TUESDAY":   time.Tuesday,
		"WEDNESDAY": time.Wednesday,
		"THURSDAY":  time.Thursday,
		"FRIDAY":    time.Friday,
		"SATURDAY":  time.Saturday,
	}

	weekday, ok := dayMap[day]
	if !ok {
		return 0, ErrInvalidDay
	}

	return weekday, nil
}
