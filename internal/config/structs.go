package config

// Location defines coordinates and name
type Location struct {
	Name      string  `yaml:"name"`
	Longitude float64 `yaml:"longitude"`
	Latitude  float64 `yaml:"latitude"`
}

// User defines the user receiving notifications
type User struct {
	TelegramUserID int64 `yaml:"telegram_user_id"`
}

// TravelTime defines the notification threshold
type TravelTime struct {
	NotificationThresholdMinutes int `yaml:"notification_threshold_minutes"`
}

// TimeSchedule defines a time and day pair
type TimeSchedule struct {
	Day  string `yaml:"day"`
	Time string `yaml:"time"`
}

// Rule represents one travel rule
type Rule struct {
	Id          int            `yaml:"id"`
	Origin      Location       `yaml:"origin"`
	Destination Location       `yaml:"destination"`
	User        User           `yaml:"user"`
	TravelTime  TravelTime     `yaml:"travel_time"`
	Times       []TimeSchedule `yaml:"times"`
	Timezone    string         `yaml:"timezone"`
}

// Config represents the full configuration
type Config struct {
	Rules []Rule `yaml:"rules"`
}
