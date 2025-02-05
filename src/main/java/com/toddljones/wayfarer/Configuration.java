package com.toddljones.wayfarer;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.DayOfWeek;
import java.time.LocalTime;
import java.time.ZoneId;
import java.util.List;

public record Configuration(
        List<Rule> rules
) {

    public record Rule(
            int id,
            Location origin,
            Location destination,
            User user,
            @JsonProperty("travel_time") TravelTime travelTime,
            List<NotificationTime> times,
            ZoneId timezone
    ) {

        public record Location(
                String name,
                double latitude,
                double longitude
        ) {
        }

        public record User(
                @JsonProperty("telegram_user_id") String telegramUserId
        ) {
        }

        public record TravelTime(
                @JsonProperty("notification_threshold_minutes") int notificationThresholdMinutes
        ) {
        }

        public record NotificationTime(
                DayOfWeek day,
                LocalTime time
        ) {
        }

    }
}
