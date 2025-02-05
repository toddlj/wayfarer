package com.toddljones.wayfarer;

import org.junit.jupiter.api.Test;

import java.time.DayOfWeek;
import java.time.Instant;
import java.time.LocalTime;
import java.time.ZoneId;
import java.util.List;

import static org.assertj.core.api.Assertions.assertThat;

class RuleEvaluationTimeCalculatorTest {

    private final RuleEvaluationTimeCalculator ruleEvaluationTimeCalculator = new RuleEvaluationTimeCalculator();

    @Test
    void calculatesNextOnSameDay() {
        // Given
        Configuration.Rule rule = createRule(List.of(
                new Configuration.Rule.NotificationTime(DayOfWeek.MONDAY, LocalTime.parse("15:15"))
        ), ZoneId.of("UTC"));
        Instant after = Instant.parse("2021-01-04T15:14:00Z");

        // When
        Instant actual = ruleEvaluationTimeCalculator.calculateNextEvaluationTime(rule, after);

        // Then
        assertThat(actual).isEqualTo(Instant.parse("2021-01-04T15:15:00Z"));
    }

    @Test
    void calculatesOnNextWeek() {
        // Given
        Configuration.Rule rule = createRule(List.of(
                new Configuration.Rule.NotificationTime(DayOfWeek.MONDAY, LocalTime.parse("15:15"))
        ), ZoneId.of("UTC"));
        Instant after = Instant.parse("2021-01-04T15:16:00Z");

        // When
        Instant actual = ruleEvaluationTimeCalculator.calculateNextEvaluationTime(rule, after);

        // Then
        assertThat(actual).isEqualTo(Instant.parse("2021-01-11T15:15:00Z"));
    }

    @Test
    void calculatesNextAmongMultipleRules() {
        // Given
        Configuration.Rule rule = createRule(List.of(
                new Configuration.Rule.NotificationTime(DayOfWeek.MONDAY, LocalTime.parse("15:15")),
                new Configuration.Rule.NotificationTime(DayOfWeek.MONDAY, LocalTime.parse("16:15"))
        ), ZoneId.of("UTC"));
        Instant after = Instant.parse("2021-01-04T15:45:00Z");

        // When
        Instant actual = ruleEvaluationTimeCalculator.calculateNextEvaluationTime(rule, after);

        // Then
        assertThat(actual).isEqualTo(Instant.parse("2021-01-04T16:15:00Z"));
    }

    @Test
    void throwsIfNoTimes() {
        // Given
        Configuration.Rule rule = createRule(List.of(), ZoneId.of("UTC"));
        Instant after = Instant.now();

        // When
        IllegalArgumentException exception = org.junit.jupiter.api.Assertions.assertThrows(IllegalArgumentException.class,
                () -> ruleEvaluationTimeCalculator.calculateNextEvaluationTime(rule, after));

        // Then
        assertThat(exception).hasMessageContaining("No times configured");
    }

    private Configuration.Rule createRule(List<Configuration.Rule.NotificationTime> times, ZoneId zone) {
        return new Configuration.Rule(1,
                new Configuration.Rule.Location("origin", 0, 0),
                new Configuration.Rule.Location("destination", 0, 0),
                new Configuration.Rule.User("user"),
                new Configuration.Rule.TravelTime(5),
                times,
                zone);
    }

}