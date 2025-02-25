package config

import (
	"os"
	"reflect"
	"testing"
)

func TestLoadConfig_HappyPath(t *testing.T) {
	// Given
	validYAML := `
rules:
  - id: 1
    origin:
      name: 10 Downing Street
      longitude: -0.1276
      latitude: 51.503
    destination:
      name: Palace of Westminster
      longitude: -0.1246
      latitude: 51.498
    user:
      telegram_user_id: 444455555
    travel_time:
      notification_threshold_minutes: 8
    times:
      - day: MONDAY
        time: 08:00
    timezone: Europe/London
`
	file := writeToFile(t, validYAML)
	defer removeFile(t, file)

	// When
	actual, err := LoadConfig(file)
	if err != nil {
		t.Fatalf("Error loading config: %s", err)
	}

	// Then
	expected := &Config{
		Rules: []Rule{
			{
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
					TelegramUserID: 444455555,
				},
				TravelTime: TravelTime{
					NotificationThresholdMinutes: 8,
				},
				Times: []TimeSchedule{
					{
						Day:  "MONDAY",
						Time: "08:00",
					},
				},
				Timezone: "Europe/London",
			},
		},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected: %+v\nGot: %+v", expected, actual)
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("non_existent_file.yaml")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	// Given
	invalidYAML := "telegram: user_id: 12345: syntax_error"
	file := writeToFile(t, invalidYAML)
	defer removeFile(t, file)

	// When
	_, err := LoadConfig(file)
	if err == nil {
		t.Fatal("expected an error due to malformed YAML, got nil")
	}
}

func writeToFile(t *testing.T, contents string) string {
	tempFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(tempFile.Name(), []byte(contents), 0644)
	if err != nil {
		t.Fatal(err)
	}
	return tempFile.Name()
}

func removeFile(t *testing.T, name string) {
	err := os.Remove(name)
	if err != nil {
		t.Fatalf("Error removing file: %s", err)
	}
}
