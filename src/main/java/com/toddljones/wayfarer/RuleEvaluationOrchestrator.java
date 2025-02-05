package com.toddljones.wayfarer;

import jakarta.inject.Inject;
import jakarta.inject.Singleton;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Instant;
import java.time.InstantSource;
import java.util.Map;
import java.util.concurrent.*;

@Singleton
public class RuleEvaluationOrchestrator {

    private static final Logger log = LoggerFactory.getLogger(RuleEvaluationOrchestrator.class);
    private final RuleEvaluator ruleEvaluator;
    private final ScheduledExecutorService schedulerThread;
    private final ExecutorService threadPool;
    private final RuleEvaluationTimeCalculator ruleEvaluationTimeCalculator;
    private final InstantSource instantSource;
    private final Configuration configuration;
    private final Map<Integer, Instant> nextEvaluationTimes;

    @Inject
    RuleEvaluationOrchestrator(RuleEvaluator ruleEvaluator,
                               RuleEvaluationTimeCalculator ruleEvaluationTimeCalculator,
                               InstantSource instantSource,
                               Configuration configuration) {
        this.ruleEvaluator = ruleEvaluator;
        this.ruleEvaluationTimeCalculator = ruleEvaluationTimeCalculator;
        this.instantSource = instantSource;
        this.configuration = configuration;
        schedulerThread = Executors.newSingleThreadScheduledExecutor();
        threadPool = Executors.newCachedThreadPool();
        nextEvaluationTimes = new ConcurrentHashMap<>(configuration.rules().size());
    }

    public void start() {
        for (Configuration.Rule rule : configuration.rules()) {
            computeAndStoreNextEvaluationTime(rule, instantSource.instant());
        }
        schedulerThread.scheduleAtFixedRate(this::evaluateScheduledRulesAndReschedule, 100, 1000, TimeUnit.MILLISECONDS);
    }

    public void stop() {
        schedulerThread.shutdown();
        threadPool.shutdown();
    }

    private void computeAndStoreNextEvaluationTime(Configuration.Rule rule, Instant evaluateAfter) {
        Instant next = ruleEvaluationTimeCalculator.calculateNextEvaluationTime(rule, evaluateAfter);
        log.info("Next evaluation time for rule {}: {}", rule.id(), next);
        nextEvaluationTimes.put(rule.id(), next);
    }

    private void evaluateScheduledRulesAndReschedule() {
        try {
            Instant now = instantSource.instant();
            for (Configuration.Rule rule : configuration.rules()) {
                if (now.isAfter(nextEvaluationTimes.get(rule.id()))) {
                    threadPool.submit(() -> {
                        ruleEvaluator.evaluateRule(rule);
                    });
                    computeAndStoreNextEvaluationTime(rule, now);
                }
            }
        } catch (Exception e) {
            log.error("Error scheduling notifications", e);
        }
    }

}
