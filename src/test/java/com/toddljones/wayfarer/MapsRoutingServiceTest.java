package com.toddljones.wayfarer;

import com.google.maps.routing.v2.*;
import com.google.protobuf.Duration;
import com.google.type.LatLng;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.ArgumentCaptor;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.io.IOException;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class MapsRoutingServiceTest {

    @Mock
    private RoutesClient routesClient;

    @Mock
    private MapsRoutingService.RoutesClientProvider routesClientProvider;

    private MapsRoutingService mapsRoutingService;

    @BeforeEach
    void setUp() throws IOException {
        when(routesClientProvider.get()).thenReturn(routesClient);
        mapsRoutingService = new MapsRoutingService(routesClientProvider);
    }

    @Test
    void shouldReturnCorrectDurationForValidRoute() {
        // Given
        WorldCoordinates origin = new WorldCoordinates(51.5074, -0.1278);      // London
        WorldCoordinates destination = new WorldCoordinates(48.8566, 2.3522);  // Paris

        java.time.Duration expectedDuration = java.time.Duration.ofHours(4);
        when(routesClient.computeRoutes(any(ComputeRoutesRequest.class)))
                .thenReturn(createResponse(expectedDuration));

        // When
        java.time.Duration actualDuration = mapsRoutingService.fetchCurrentTransitTimeBetween(origin, destination);

        // Then
        assertThat(actualDuration).isEqualTo(expectedDuration);

        ArgumentCaptor<ComputeRoutesRequest> requestCaptor = ArgumentCaptor.forClass(ComputeRoutesRequest.class);
        verify(routesClient).computeRoutes(requestCaptor.capture());
        assertThat(requestCaptor.getValue()).satisfies(request -> {
            LatLng originLatLng = request.getOrigin().getLocation().getLatLng();
            LatLng destLatLng = request.getDestination().getLocation().getLatLng();
            assertThat(originLatLng.getLatitude()).isEqualTo(origin.latitude());
            assertThat(originLatLng.getLongitude()).isEqualTo(origin.longitude());
            assertThat(destLatLng.getLatitude()).isEqualTo(destination.latitude());
            assertThat(destLatLng.getLongitude()).isEqualTo(destination.longitude());
            assertThat(request.getTravelMode()).isEqualTo(RouteTravelMode.TRANSIT);
            assertThat(request.getComputeAlternativeRoutes()).isFalse();
        });
    }

    private static ComputeRoutesResponse createResponse(java.time.Duration expectedDuration) {
        long expectedDurationSeconds = expectedDuration.getSeconds();

        Route route = Route.newBuilder()
                .setDuration(Duration.newBuilder().setSeconds(expectedDurationSeconds).build())
                .build();
        return ComputeRoutesResponse.newBuilder()
                .addRoutes(route)
                .build();
    }

}