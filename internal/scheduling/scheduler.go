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

// Can be overridden in tests
var getNextScheduledTimeFunction = getNextScheduledTime

// ScheduleFunction holidays supplied in the format "2026-01-01"
func ScheduleFunction(schedules []Schedule, timezone *time.Location, holidays []string, task func()) error {
	holidayMap := make(map[string]bool)
	for _, h := range holidays {
		holidayMap[h] = true
	}

	for _, schedule := range schedules {
		schedule := schedule // Capture range variable
		go func() {
			var scheduleNext func(allowToday bool)
			scheduleNext = func(allowToday bool) {
				nextRun := getNextScheduledTimeFunction(time.Now(), schedule, timezone, allowToday, holidayMap)
				slog.Info("Scheduled task", slog.Any("next_run", nextRun))

				time.AfterFunc(time.Until(nextRun), func() {
					go task()           // Run the task in a separate Goroutine
					scheduleNext(false) // Recursively reschedule the next execution
				})
			}

			scheduleNext(true)
		}()
	}

	return nil
}

func getNextScheduledTime(now time.Time, schedule Schedule, timezone *time.Location, allowToday bool, holidays map[string]bool) time.Time {
	nowInTimezone := now.In(timezone)

	for days := 0; days < 366; days++ {
		if !allowToday && days <= 1 {
			// Start at 2 to avoid running the task today (or tomorrow for race conditions where it runs just before midnight)
			continue
		}
		nextRun := time.Date(nowInTimezone.Year(), nowInTimezone.Month(), nowInTimezone.Day()+days,
			schedule.Hour, schedule.Minute, 0, 0, timezone)
		if nextRun.Weekday() == schedule.DayOfWeek && nextRun.After(nowInTimezone) {
			// Check if this date is a holiday
			dateKey := nextRun.Format("2006-01-02")
			if holidays[dateKey] {
				slog.Debug("Skipping scheduled run due to holiday", slog.String("date", dateKey))
				continue
			}
			return nextRun
		}
	}

	slog.Error("Failed to calculate next scheduled time within a year")
	return now.Add(7 * 24 * time.Hour)
}
