package com.toddljones.wayfarer;

import java.time.Instant;
import java.time.ZoneId;
import java.time.ZonedDateTime;

public class RuleEvaluationTimeCalculator {

    public Instant calculateNextEvaluationTime(Configuration.Rule rule, Instant after) {
        if (rule.times().isEmpty()) {
            throw new IllegalArgumentException("No times configured for rule " + rule.id());
        }
        ZoneId zoneId = rule.timezone();
        return rule.times().stream()
                .map(time -> {
                    ZonedDateTime nextTime = after.atZone(zoneId)
                            .with(time.day())
                            .with(time.time());
                    if (!nextTime.toInstant().isAfter(after)) {
                        return nextTime.plusWeeks(1).toInstant();
                    }
                    return nextTime.toInstant();
                })
                .min(Instant::compareTo)
                .orElseThrow();
    }
}
