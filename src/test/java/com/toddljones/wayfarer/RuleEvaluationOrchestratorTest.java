package com.toddljones.wayfarer;

import org.awaitility.Awaitility;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.time.*;
import java.util.List;
import java.util.concurrent.atomic.AtomicReference;

import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class RuleEvaluationOrchestratorTest {

    @Mock
    private RuleEvaluator ruleEvaluator;

    @Mock
    private InstantSource instantSource;

    private RuleEvaluationOrchestrator ruleEvaluationOrchestrator;

    @BeforeEach
    void setUp() {
        Configuration configuration = new Configuration(
                List.of(
                        new Configuration.Rule(
                                1,
                                new Configuration.Rule.Location("Home", -0.1276, 51.503),
                                new Configuration.Rule.Location("Work", -0.1246, 51.498),
                                new Configuration.Rule.User("12345"),
                                new Configuration.Rule.TravelTime(30),
                                List.of(new Configuration.Rule.NotificationTime(DayOfWeek.MONDAY, LocalTime.parse("15:15"))),
                                ZoneId.of("UTC")
                        ),
                        new Configuration.Rule(
                                2,
                                new Configuration.Rule.Location("Home", -0.1276, 51.503),
                                new Configuration.Rule.Location("Work", -0.1246, 51.498),
                                new Configuration.Rule.User("12345"),
                                new Configuration.Rule.TravelTime(30),
                                List.of(new Configuration.Rule.NotificationTime(DayOfWeek.TUESDAY, LocalTime.parse("16:15"))),
                                ZoneId.of("UTC")
                        )
                )
        );
        ruleEvaluationOrchestrator = new RuleEvaluationOrchestrator(ruleEvaluator, new RuleEvaluationTimeCalculator(), instantSource, configuration);
    }

    @AfterEach
    void tearDown() {
        ruleEvaluationOrchestrator.stop();
    }

    @Test
    void triggersRuleEvaluations() {
        // Given
        AtomicReference<Instant> now = new AtomicReference<>(Instant.parse("2021-01-04T15:14:00Z"));
        when(instantSource.instant()).thenAnswer(_ -> now.get());

        // Start
        ruleEvaluationOrchestrator.start();

        // When 1
        now.set(Instant.parse("2021-01-04T15:15:01Z"));

        // Then 1
        Awaitility.await().atMost(Duration.ofSeconds(3))
                .until(() -> !mockingDetails(ruleEvaluator).getInvocations().isEmpty());
        verify(ruleEvaluator).evaluateRule(argThat(rule -> rule.id() == 1));
        verifyNoMoreInteractions(ruleEvaluator);
        reset(ruleEvaluator);

        // When 2
        now.set(Instant.parse("2021-01-05T16:15:01Z"));

        // Then 2
        Awaitility.await().atMost(Duration.ofSeconds(3)).until(() -> !mockingDetails(ruleEvaluator).getInvocations().isEmpty());
        verify(ruleEvaluator).evaluateRule(argThat(rule -> rule.id() == 2));
        verifyNoMoreInteractions(ruleEvaluator);
    }
}