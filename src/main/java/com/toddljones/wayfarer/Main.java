package com.toddljones.wayfarer;

import java.io.File;

public class Main {
    public static void main(String[] args) {
        // Dependency injection
        MapsRoutingService routingService = new MapsRoutingService(new MapsRoutingService.RoutesClientProvider());
        TelegramNotifier telegramNotifier = new TelegramNotifier();
        ConfigurationParser configurationParser = new ConfigurationParser();
        RuleEvaluator ruleEvaluator = new RuleEvaluator(telegramNotifier, routingService);

        // Configuration
        Configuration configuration = configurationParser.parse(new File("config.yaml"));

        // Execution
        Configuration.Rule rule = configuration.rules().getFirst();
        ruleEvaluator.evaluateRule(rule);
    }
}