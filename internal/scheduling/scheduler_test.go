package scheduling

import (
	"testing"
	"time"
)

func Test_GetNextScheduledTime(t *testing.T) {
	// Define test cases
	tests := []struct {
		name       string
		now        time.Time
		schedule   Schedule
		timezone   *time.Location
		allowToday bool
		expected   time.Time
	}{
		{
			name:       "Schedule on a future weekday",
			now:        time.Date(2025, 2, 10, 12, 0, 0, 0, time.UTC), // Monday
			schedule:   Schedule{DayOfWeek: time.Wednesday, Hour: 9, Minute: 0},
			timezone:   time.UTC,
			allowToday: false,
			expected:   time.Date(2025, 2, 12, 9, 0, 0, 0, time.UTC), // Wednesday at 09:00 UTC
		},
		{
			name:       "Schedule on the same day when allowToday is true",
			now:        time.Date(2025, 2, 10, 8, 0, 0, 0, time.UTC), // Monday 08:00 UTC
			schedule:   Schedule{DayOfWeek: time.Monday, Hour: 9, Minute: 0},
			timezone:   time.UTC,
			allowToday: true,
			expected:   time.Date(2025, 2, 10, 9, 0, 0, 0, time.UTC), // Monday at 09:00 UTC
		},
		{
			name:       "Skip today when allowToday is false",
			now:        time.Date(2025, 2, 10, 8, 0, 0, 0, time.UTC), // Monday 08:00 UTC
			schedule:   Schedule{DayOfWeek: time.Monday, Hour: 9, Minute: 0},
			timezone:   time.UTC,
			allowToday: false,
			expected:   time.Date(2025, 2, 17, 9, 0, 0, 0, time.UTC), // Next Monday at 09:00 UTC
		},
		{
			name:       "Wrap around to next week",
			now:        time.Date(2025, 2, 10, 12, 0, 0, 0, time.UTC), // Monday
			schedule:   Schedule{DayOfWeek: time.Sunday, Hour: 9, Minute: 0},
			timezone:   time.UTC,
			allowToday: false,
			expected:   time.Date(2025, 2, 16, 9, 0, 0, 0, time.UTC), // Next Sunday at 09:00 UTC
		},
		{
			name:       "Skip earlier today",
			now:        time.Date(2025, 2, 10, 12, 0, 0, 0, time.UTC), // Monday
			schedule:   Schedule{DayOfWeek: time.Monday, Hour: 9, Minute: 0},
			timezone:   time.UTC,
			allowToday: true,
			expected:   time.Date(2025, 2, 17, 9, 0, 0, 0, time.UTC), // Next Monday at 09:00 UTC
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getNextScheduledTime(tt.now, tt.schedule, tt.timezone, tt.allowToday)

			// Verify the exact scheduled time
			if !result.Equal(tt.expected) {
				t.Errorf("Test %s failed:\nExpected: %v\nGot:      %v", tt.name, tt.expected, result)
			}
		})
	}
}

func TestScheduleFunction_TaskExecutionAndRescheduling(t *testing.T) {
	// Save original function so we can restore it later.
	origNextFunc := getNextScheduledTimeFunction
	// Override getNextScheduledTimeFunc to always schedule the task 10ms in the future.
	getNextScheduledTimeFunction = func(now time.Time, schedule Schedule, timezone *time.Location, allowToday bool) time.Time {
		return time.Now().Add(10 * time.Millisecond)
	}
	defer func() { getNextScheduledTimeFunction = origNextFunc }()

	// Use a channel to signal task execution.
	executionCount := 0
	execCh := make(chan struct{}, 10)
	task := func() {
		executionCount++
		execCh <- struct{}{}
	}

	// Create a dummy schedule. Its values don’t matter because our override forces a 10ms delay.
	now := time.Now()
	schedules := []Schedule{
		{DayOfWeek: now.Weekday(), Hour: now.Hour(), Minute: now.Minute()},
	}
	loc := time.UTC

	// Start scheduling the task.
	if err := ScheduleFunction(schedules, loc, task); err != nil {
		t.Fatalf("ScheduleFunction returned error: %v", err)
	}

	// Allow some time for several task executions.
	timeout := time.After(100 * time.Millisecond)
Loop:
	for {
		select {
		case <-execCh:
			if executionCount >= 2 {
				break Loop
			}
		case <-timeout:
			break Loop
		}
	}

	// We expect the task to have been executed multiple times.
	if executionCount < 2 {
		t.Errorf("expected task to execute at least twice, but got %d", executionCount)
	}
}
