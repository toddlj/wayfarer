package com.toddljones.wayfarer;

public class Main {
    public static void main(String[] args) {
        // Dependency injection
        MapsRoutingService routingService = new MapsRoutingService(new MapsRoutingService.RoutesClientProvider());

        // Configuration
        WorldCoordinates origin = new WorldCoordinates(51.5074, -0.1278);      // London
        WorldCoordinates destination = new WorldCoordinates(48.8566, 2.3522);  // Paris

        // Execution
        java.time.Duration duration = routingService.fetchCurrentTransitTimeBetween(
                origin, destination
        );

        // Output
        System.out.println("Duration: " + duration);
    }
}