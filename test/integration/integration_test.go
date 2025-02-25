package integration

import (
	"bytes"
	"cloud.google.com/go/maps/routing/apiv2/routingpb"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/durationpb"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"
)

func Test_IntegrationTest(t *testing.T) {
	// Static test data
	telegramUserId := int64(444444444)
	currentRouteDurationMinutes := 10
	routeDurationThresholdMinutes := 8
	telegramToken := "TOKENTOKENTOKEN"
	googleApiKey := "FAKEKEYFAKEKEYFAKEKEY"

	// Create config file
	config := generateConfig(telegramUserId, routeDurationThresholdMinutes)
	filename := saveConfigFile(t, config)
	defer removeConfigFile(t, filename)

	// Start mock Telegram API Server
	var telegramRequests []TelegramMessage
	mockServer := httptest.NewServer(http.HandlerFunc(handleTelegramCall(t, telegramToken, &telegramRequests)))
	defer mockServer.Close()

	// Start mock Google API server
	port := startGoogleServer(t, currentRouteDurationMinutes)

	// Run the application
	cmd := exec.Command("../../wayfarer", "--config-file", "test_data/test_config.yaml")
	cmd.Env = append(os.Environ(),
		"TELEGRAM_BOT_TOKEN="+telegramToken,
		"TELEGRAM_API_BASE_URL="+mockServer.URL,
		"GOOGLE_API_KEY="+googleApiKey,
		"GOOGLE_API_BASE_URL=localhost:"+port,
	)
	printStdErrAsync(t, cmd)
	defer killProcess(t, cmd)
	errChan := make(chan error, 1)
	go func() {
		errChan <- cmd.Run() // Runs and waits for completion
	}()

	// Wait for Telegram message to be sent
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	timeout := time.After(70 * time.Second)
Poll:
	for pollCounter := 1; ; pollCounter++ {
		select {
		case <-ticker.C:
			t.Logf("Checking for sent Telegram message (attempt %d)...", pollCounter)
			if len(telegramRequests) > 0 {
				break Poll
			}
		case <-timeout:
			t.Error("Timeout reached: No Telegram message sent.")
			break Poll
		case err := <-errChan:
			if err != nil {
				t.Errorf("Process finished with error: %v", err)
			} else {
				t.Log("Process finished successfully.")
			}
			break Poll
		}
	}

	// Verify the Telegram message was sent correctly
	if len(telegramRequests) < 1 {
		t.Fatalf("Expected at least 1 request, but got %d", len(telegramRequests))
	}
	expectedMessage := TelegramMessage{ChatID: telegramUserId,
		Text: fmt.Sprintf("Travel time between 10 Downing Street and Palace of Westminster is greater than %d minutes: currently scheduled to take %d minutes",
			routeDurationThresholdMinutes, currentRouteDurationMinutes)}
	if telegramRequests[0] != expectedMessage {
		t.Errorf("Unexpected request body: got %+v, expected %+v", telegramRequests[0], expectedMessage)
	} else {
		t.Logf("Request body verified: %+v", telegramRequests[0])
	}
}

func killProcess(t *testing.T, cmd *exec.Cmd) {
	if cmd.Process != nil {
		err := cmd.Process.Kill()
		if err != nil {
			t.Errorf("Failed to kill process: %v", err)
		} else {
			t.Log("Process killed successfully")
		}
	}
}

type TelegramMessage struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func printStdErrAsync(t *testing.T, cmd *exec.Cmd) {
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf
	var mu sync.Mutex
	go func() {
		var printedUpTo int
		for {
			if len(stderrBuf.String()[printedUpTo:]) > 0 {
				mu.Lock()
				t.Logf("%s", stderrBuf.String()[printedUpTo:])
				printedUpTo = len(stderrBuf.String())
				mu.Unlock()
			}
		}
	}()
}

func generateConfig(telegramUserId int64, routeDurationThresholdMinutes int) string {
	now := time.Now()
	today := strings.ToUpper(now.Weekday().String())
	tomorrow := strings.ToUpper(now.Add(24 * time.Hour).Weekday().String())
	var notifications string
	for i := 1; i <= 10; i++ {
		notificationTime := now.Add(time.Duration(i) * time.Minute).Format(time.TimeOnly)[:5]
		notifications += fmt.Sprintf("      - day: %s\n        time: %s\n", today, notificationTime)
		notifications += fmt.Sprintf("      - day: %s\n        time: %s\n", tomorrow, notificationTime)
	}
	config := fmt.Sprintf(`rules:
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
      telegram_user_id: %d
    travel_time:
      notification_threshold_minutes: %d
    times:
%s    timezone: UTC`, telegramUserId, routeDurationThresholdMinutes, notifications)
	return config
}

func saveConfigFile(t *testing.T, config string) string {
	err := os.Mkdir("test_data", 0755)
	if err != nil && !os.IsExist(err) {
		t.Fatalf("Error creating directory: %s", err)

	}
	filename := "test_data/test_config.yaml"
	err = os.WriteFile(filename, []byte(config), 0644)
	if err != nil {
		t.Fatalf("Error writing file: %s", err)
	}
	return filename
}

func removeConfigFile(t *testing.T, filename string) {
	err := os.Remove(filename)
	if err != nil {
		t.Fatalf("Error removing file: %s", err)
	}
}

func handleTelegramCall(t *testing.T, telegramToken string, telegramRequests *[]TelegramMessage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate URL
		expectedURL := fmt.Sprintf("/bot%s/sendMessage", telegramToken)
		if r.URL.Path != expectedURL {
			t.Errorf("Unexpected URL: got %s, expected %s", r.URL.Path, expectedURL)
		}

		// Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Fatalf("Failed to close request body: %v", err)
			}
		}(r.Body)

		// Unmarshal JSON request body
		var msg TelegramMessage
		if err := json.Unmarshal(body, &msg); err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		// Store the request for later verification
		*telegramRequests = append(*telegramRequests, msg)

		// Respond with OK response
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(`{"ok":true}`))
		if err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}
}

type MockRoutingServer struct {
	routingpb.RoutesServer
	t                    *testing.T
	routeDurationMinutes int
}

func (s *MockRoutingServer) ComputeRoutes(_ context.Context, req *routingpb.ComputeRoutesRequest) (*routingpb.ComputeRoutesResponse, error) {
	s.t.Logf("Received ComputeRoutesRequest: %+v", req)

	duration := time.Duration(s.routeDurationMinutes) * time.Minute
	durationProto := durationpb.New(duration)

	return &routingpb.ComputeRoutesResponse{
		Routes: []*routingpb.Route{
			{
				Duration: durationProto,
			},
		},
	}, nil
}

func startGoogleServer(t *testing.T, routeDurationMinutes int) (port string) {
	server := grpc.NewServer()
	routingpb.RegisterRoutesServer(server, &MockRoutingServer{t: t, routeDurationMinutes: routeDurationMinutes})

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	reflection.Register(server)

	_, port, err = net.SplitHostPort(listener.Addr().String())
	if err != nil {
		log.Fatalf("failed to split host and port: %v", err)
	}

	t.Logf("Mock gRPC Google API server listening on localhost:%s...", port)
	go func() {
		if err := server.Serve(listener); err != nil {
			t.Errorf("failed to serve: %v", err)
			return
		}
	}()

	return port
}
