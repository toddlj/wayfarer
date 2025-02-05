package com.toddljones.wayfarer;

import com.google.maps.routing.v2.ComputeRoutesRequest;
import com.google.maps.routing.v2.ComputeRoutesResponse;
import com.google.maps.routing.v2.Location;
import com.google.maps.routing.v2.RouteTravelMode;
import com.google.maps.routing.v2.RoutesClient;
import com.google.maps.routing.v2.RoutesSettings;
import com.google.maps.routing.v2.Waypoint;
import com.google.type.LatLng;
import jakarta.inject.Inject;
import jakarta.inject.Singleton;

import java.io.IOException;
import java.time.Duration;
import java.util.Map;

@Singleton
public final class MapsRoutingService {

    @Singleton
    public static class RoutesClientProvider {
        private final String googleApiKey;

        public RoutesClientProvider() {
            this.googleApiKey = System.getenv("GOOGLE_API_KEY");
            if (googleApiKey == null) {
                throw new IllegalStateException("GOOGLE_API_KEY environment variable must be set");
            }
        }

        RoutesClient get() throws IOException {
            return RoutesClient.create(
                    RoutesSettings.newBuilder()
                            .setApiKey(googleApiKey)
                            .setHeaderProvider(() -> Map.of("X-Goog-FieldMask", "routes.duration"))
                            .build()
            );
        }
    }

    private final RoutesClientProvider routesClientProvider;

    @Inject
    MapsRoutingService(RoutesClientProvider routesClientProvider) {
        this.routesClientProvider = routesClientProvider;
    }

    public Duration fetchCurrentTransitTimeBetween(WorldCoordinates origin, WorldCoordinates destination) {
        try (RoutesClient routesClient = routesClientProvider.get()) {
            ComputeRoutesRequest request = ComputeRoutesRequest.newBuilder()
                    .setOrigin(Waypoint.newBuilder().setLocation(
                            Location.newBuilder().setLatLng(
                                    LatLng.newBuilder()
                                            .setLatitude(origin.latitude())
                                            .setLongitude(origin.longitude())
                                            .build())))
                    .setDestination(Waypoint.newBuilder()
                            .setLocation(
                                    Location.newBuilder()
                                            .setLatLng(
                                                    LatLng.newBuilder()
                                                            .setLatitude(destination.latitude())
                                                            .setLongitude(destination.longitude()))))
                    .setTravelMode(RouteTravelMode.TRANSIT)
                    .setComputeAlternativeRoutes(false)
                    .build();
            ComputeRoutesResponse response = routesClient.computeRoutes(request);
            long durationSeconds = response.getRoutes(0).getDuration().getSeconds();
            return Duration.ofSeconds(durationSeconds);
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }
}
