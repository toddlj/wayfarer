package config

import (
	"errors"
	"time"
)

var errInvalidDay = errors.New("invalid day of the week")

var errInvalidTimeFormat = errors.New("invalid time format")

func (cfg *Config) validate() error {
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

		// validate each time schedule
		for _, t := range rule.Times {
			// validate day
			if _, err := ParseWeekday(t.Day); err != nil {
				return err
			}

			// validate time format
			if _, err := time.Parse("15:04", t.Time); err != nil {
				return errInvalidTimeFormat
			}
		}

		// validate timezone
		if _, err := time.LoadLocation(rule.Timezone); err != nil {
			return err
		}
	}

	// All checks passed
	return nil
}
