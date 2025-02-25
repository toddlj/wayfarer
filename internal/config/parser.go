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

	return parseAndValidateConfig(data)
}

func parseAndValidateConfig(rawConfig []byte) (*Config, error) {
	var config Config
	if err := yaml.Unmarshal(rawConfig, &config); err != nil {
		return nil, err
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return &config, nil
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
		return 0, errInvalidDay
	}

	return weekday, nil
}

var errInvalidDay = errors.New("invalid day of the week")
