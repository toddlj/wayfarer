package com.toddljones.wayfarer;

import java.io.File;
import java.time.Duration;

public class Main {
    public static void main(String[] args) {
        // Dependency injection
        MapsRoutingService routingService = new MapsRoutingService(new MapsRoutingService.RoutesClientProvider());
        TelegramNotifier telegramNotifier = new TelegramNotifier();
        ConfigurationParser configurationParser = new ConfigurationParser();

        // Configuration
        Configuration configuration = configurationParser.parse(new File("config.yaml"));

        // Execution
        Configuration.Rule rule = configuration.rules().getFirst();
        java.time.Duration duration = routingService.fetchCurrentTransitTimeBetween(rule.origin(), rule.destination());
        if (duration.compareTo(Duration.ofMinutes(rule.travelTime().notificationThresholdMinutes())) > 0) {
            telegramNotifier.sendMessage(rule.user().telegramUserId(), "Travel time from %s to %s is greater than %s minutes. It is currently %s minutes".formatted(rule.origin().name(), rule.destination().name(), rule.travelTime().notificationThresholdMinutes(), duration.toMinutes()));
        }
    }
}