package scheduling

import (
	"log/slog"
	"time"
)

// Schedule represents a day and time to run the function
type Schedule struct {
	DayOfWeek time.Weekday // e.g., time.Monday
	Hour      int          // e.g., 14
	Minute    int          // e.g., 30
}

func ScheduleFunction(schedules []Schedule, timezone *time.Location, task func()) error {
	// Schedule each task in a separate Goroutine
	for _, schedule := range schedules {
		schedule := schedule // Capture range variable
		go func() {
			var scheduleNext func(allowToday bool)
			scheduleNext = func(allowToday bool) {
				nextRun := getNextScheduledTime(time.Now(), schedule, timezone, allowToday)
				slog.Info("Scheduled task", slog.Any("next_run", nextRun))

				time.AfterFunc(time.Until(nextRun), func() {
					go task()           // Run function in a separate Goroutine
					scheduleNext(false) // Recursively reschedule the next execution
				})
			}

			scheduleNext(true)
		}()
	}

	return nil
}

func getNextScheduledTime(now time.Time, schedule Schedule, timezone *time.Location, allowToday bool) time.Time {
	nowInTimezone := now.In(timezone)

	for days := 0; days < 8; days++ {
		if !allowToday && days <= 1 {
			// Start at 2 to avoid running the task today (or tomorrow for race conditions where it runs just before midnight)
			continue
		}
		nextRun := time.Date(nowInTimezone.Year(), nowInTimezone.Month(), nowInTimezone.Day()+days,
			schedule.Hour, schedule.Minute, 0, 0, timezone)
		if nextRun.Weekday() == schedule.DayOfWeek && nextRun.After(nowInTimezone) {
			return nextRun
		}
	}

	// Fallback, should never happen
	slog.Error("Failed to calculate next scheduled time, using current time plus 7 days")
	return now.Add(7 * 24 * time.Hour)
}
