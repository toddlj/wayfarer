package config

import (
	"strings"
	"testing"
)

// validConfig returns a baseline valid configuration.
func validConfig() Config {
	return Config{
		Rules: []Rule{{
			Id: 1,
			Origin: Location{
				Name:      "10 Downing Street",
				Longitude: -0.1276,
				Latitude:  51.503,
			},
			Destination: Location{
				Name:      "Palace of Westminster",
				Longitude: -0.1246,
				Latitude:  51.498,
			},
			User: User{
				TelegramUserID: 123456789,
			},
			TravelTime: TravelTime{
				NotificationThresholdMinutes: 10,
			},
			Times: []TimeSchedule{{
				Day:  "MONDAY",
				Time: "09:00",
			}},
			Timezone: "Europe/London",
		}},
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid config",
			cfg:     validConfig(),
			wantErr: false,
		},
		{
			name: "invalid rule id",
			cfg: func() Config {
				cfg := validConfig()
				cfg.Rules[0].Id = 0
				return cfg
			}(),
			wantErr: true,
			errMsg:  "id must be greater than 0",
		},
		{
			name: "missing origin name",
			cfg: func() Config {
				cfg := validConfig()
				cfg.Rules[0].Origin.Name = ""
				return cfg
			}(),
			wantErr: true,
			errMsg:  "origin and destination must have a name",
		},
		{
			name: "missing destination name",
			cfg: func() Config {
				cfg := validConfig()
				cfg.Rules[0].Destination.Name = ""
				return cfg
			}(),
			wantErr: true,
			errMsg:  "origin and destination must have a name",
		},
		{
			name: "missing telegram user id",
			cfg: func() Config {
				cfg := validConfig()
				cfg.Rules[0].User.TelegramUserID = 0
				return cfg
			}(),
			wantErr: true,
			errMsg:  "user must have a Telegram user ID",
		},
		{
			name: "notification threshold minutes not positive",
			cfg: func() Config {
				cfg := validConfig()
				cfg.Rules[0].TravelTime.NotificationThresholdMinutes = 0
				return cfg
			}(),
			wantErr: true,
			errMsg:  "notification_threshold_minutes must be greater than 0",
		},
		{
			name: "empty times list",
			cfg: func() Config {
				cfg := validConfig()
				cfg.Rules[0].Times = []TimeSchedule{}
				return cfg
			}(),
			wantErr: true,
			errMsg:  "at least one time must be specified",
		},
		{
			name: "invalid weekday",
			cfg: func() Config {
				cfg := validConfig()
				cfg.Rules[0].Times[0].Day = "FUNDAY"
				return cfg
			}(),
			wantErr: true,
			errMsg:  errInvalidDay.Error(),
		},
		{
			name: "invalid time format",
			cfg: func() Config {
				cfg := validConfig()
				cfg.Rules[0].Times[0].Time = "9 AM"
				return cfg
			}(),
			wantErr: true,
			errMsg:  errInvalidTimeFormat.Error(),
		},
		{
			name: "invalid timezone",
			cfg: func() Config {
				cfg := validConfig()
				cfg.Rules[0].Timezone = "Invalid/Timezone"
				return cfg
			}(),
			wantErr: true,
			errMsg:  "unknown time zone",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.validate()
			if tt.wantErr && err == nil {
				t.Fatalf("expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantErr && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error to contain %q, got %q", tt.errMsg, err.Error())
				}
			}
		})
	}
}
