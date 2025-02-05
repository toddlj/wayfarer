package com.toddljones.wayfarer;

import java.io.File;
import java.time.Instant;
import java.time.InstantSource;

public class Main {
    public static void main(String[] args) {
        // Dependency injection
        ConfigurationParser configurationParser = new ConfigurationParser();
        Configuration configuration = configurationParser.parse(new File("config.yaml"));
        MapsRoutingService routingService = new MapsRoutingService(new MapsRoutingService.RoutesClientProvider());
        TelegramNotifier telegramNotifier = new TelegramNotifier();
        RuleEvaluator ruleEvaluator = new RuleEvaluator(telegramNotifier, routingService);
        RuleEvaluationTimeCalculator ruleEvaluationTimeCalculator = new RuleEvaluationTimeCalculator();
        InstantSource instantSource = Instant::now;
        RuleEvaluationOrchestrator ruleEvaluationOrchestrator = new RuleEvaluationOrchestrator(ruleEvaluator, ruleEvaluationTimeCalculator, instantSource, configuration);

        // Execution
        ruleEvaluationOrchestrator.start();
    }
}