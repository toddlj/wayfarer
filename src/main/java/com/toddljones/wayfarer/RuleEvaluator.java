package com.toddljones.wayfarer;

import jakarta.inject.Inject;
import jakarta.inject.Singleton;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;

@Singleton
public class RuleEvaluator {

    private static final Logger log = LoggerFactory.getLogger(RuleEvaluator.class);
    private final TelegramNotifier telegramNotifier;
    private final MapsRoutingService routingService;

    @Inject
    RuleEvaluator(TelegramNotifier telegramNotifier, MapsRoutingService routingService) {
        this.telegramNotifier = telegramNotifier;
        this.routingService = routingService;
    }

    public void evaluateRule(Configuration.Rule rule) {
        Duration currentDuration = routingService.fetchCurrentTransitTimeBetween(rule.origin(), rule.destination());
        Duration notificationThreshold = Duration.ofMinutes(rule.travelTime().notificationThresholdMinutes());
        if (currentDuration.compareTo(notificationThreshold) > 0) {
            log.info("Travel time {} from {} to {} is greater than {}. Sending notification.",
                    currentDuration, rule.origin().name(), rule.destination().name(), notificationThreshold);
            telegramNotifier.sendMessage(rule.user().telegramUserId(), formatMessage(rule, currentDuration));
        } else {
            log.info("Travel time {} from {} to {} is less than {}. Doing nothing.",
                    currentDuration, rule.origin().name(), rule.destination().name(), notificationThreshold);
        }
    }

    private static String formatMessage(Configuration.Rule rule, Duration currentDuration) {
        return "Travel time from %s to %s is greater than %s minutes. It is currently %s minutes."
                .formatted(
                        rule.origin().name(),
                        rule.destination().name(),
                        rule.travelTime().notificationThresholdMinutes(),
                        currentDuration.toMinutes());
    }
}
