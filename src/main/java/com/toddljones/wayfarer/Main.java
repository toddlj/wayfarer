package com.toddljones.wayfarer;

import org.apache.commons.cli.*;

import java.io.File;
import java.time.Instant;
import java.time.InstantSource;

public class Main {
    public static void main(String[] args) {
        // Arguments
        Options options = new Options();
        options.addRequiredOption(null, "config-file", true, "Path to the config file");
        CommandLine cmd = parseCommandLineArguments(args, options);
        File configFile = new File(cmd.getOptionValue("config-file"));

        // Dependency injection
        ConfigurationParser configurationParser = new ConfigurationParser();
        Configuration configuration = configurationParser.parse(configFile);
        MapsRoutingService routingService = new MapsRoutingService(new MapsRoutingService.RoutesClientProvider());
        TelegramNotifier telegramNotifier = new TelegramNotifier();
        RuleEvaluator ruleEvaluator = new RuleEvaluator(telegramNotifier, routingService);
        RuleEvaluationTimeCalculator ruleEvaluationTimeCalculator = new RuleEvaluationTimeCalculator();
        InstantSource instantSource = Instant::now;
        RuleEvaluationOrchestrator ruleEvaluationOrchestrator = new RuleEvaluationOrchestrator(ruleEvaluator, ruleEvaluationTimeCalculator, instantSource, configuration);

        // Execution
        ruleEvaluationOrchestrator.start();
    }

    private static CommandLine parseCommandLineArguments(String[] args, Options options) {
        CommandLineParser parser = new DefaultParser();
        CommandLine cmd;
        try {
            cmd = parser.parse(options, args);
        } catch (ParseException e) {
            throw new RuntimeException("Error parsing command line arguments", e);
        }
        return cmd;
    }
}