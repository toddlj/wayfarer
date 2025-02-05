package com.toddljones.wayfarer;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.time.Duration;

import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class RuleEvaluatorTest {

    @Mock
    private TelegramNotifier telegramNotifier;

    @Mock
    private MapsRoutingService routingService;

    private RuleEvaluator ruleEvaluator;

    @BeforeEach
    void setUp() {
        ruleEvaluator = new RuleEvaluator(telegramNotifier, routingService);
    }

    @Test
    void shouldSendNotificationWhenTravelTimeExceedsThreshold() {
        // Given
        Configuration.Rule rule = new Configuration.Rule(
                1,
                new Configuration.Rule.Location("Home", -0.1276, 51.503),
                new Configuration.Rule.Location("Work", -0.1246, 51.498),
                new Configuration.Rule.User("12345"),
                new Configuration.Rule.TravelTime(30),
                null,
                null);
        when(routingService.fetchCurrentTransitTimeBetween(rule.origin(), rule.destination()))
                .thenReturn(Duration.ofMinutes(45));

        // When
        ruleEvaluator.evaluateRule(rule);

        // Then
        String expectedMessage = "Travel time from Home to Work is greater than 30 minutes. It is currently 45 minutes.";
        verify(telegramNotifier).sendMessage("12345", expectedMessage);
    }

    @Test
    void shouldNotSendNotificationWhenTravelTimeUnderThreshold() {
        // Given
        Configuration.Rule rule = new Configuration.Rule(
                1,
                new Configuration.Rule.Location("Home", -0.1276, 51.503),
                new Configuration.Rule.Location("Work", -0.1246, 51.498),
                new Configuration.Rule.User("12345"),
                new Configuration.Rule.TravelTime(30),
                null,
                null);
        when(routingService.fetchCurrentTransitTimeBetween(rule.origin(), rule.destination()))
                .thenReturn(Duration.ofMinutes(30));

        // When
        ruleEvaluator.evaluateRule(rule);

        // Then
        verify(telegramNotifier, never()).sendMessage(anyString(), anyString());
    }

    @Test
    void shouldNotSendNotificationWhenTravelTimeEqualsThreshold() {
        // Given
        Configuration.Rule rule = new Configuration.Rule(
                1,
                new Configuration.Rule.Location("Home", -0.1276, 51.503),
                new Configuration.Rule.Location("Work", -0.1246, 51.498),
                new Configuration.Rule.User("12345"),
                new Configuration.Rule.TravelTime(30),
                null,
                null);
        when(routingService.fetchCurrentTransitTimeBetween(rule.origin(), rule.destination()))
                .thenReturn(Duration.ofMinutes(30));

        // When
        ruleEvaluator.evaluateRule(rule);

        // Then
        verify(telegramNotifier, never()).sendMessage(anyString(), anyString());
    }
}