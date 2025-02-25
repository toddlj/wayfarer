package googlemaps

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/maps/routing/apiv2/routingpb"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/genproto/googleapis/type/latlng"
	"google.golang.org/protobuf/types/known/durationpb"
)

type fakeRoutesClient struct {
	computeRoutesFunc func(ctx context.Context, req *routingpb.ComputeRoutesRequest, opts ...gax.CallOption) (*routingpb.ComputeRoutesResponse, error)
}

func (f *fakeRoutesClient) ComputeRoutes(ctx context.Context, req *routingpb.ComputeRoutesRequest, opts ...gax.CallOption) (*routingpb.ComputeRoutesResponse, error) {
	return f.computeRoutesFunc(ctx, req, opts...)
}

func (f *fakeRoutesClient) Close() error {
	return nil
}

func TestFetchCurrentTransitTimeBetween_HappyPath(t *testing.T) {
	// Given
	fakeClient := &fakeRoutesClient{
		computeRoutesFunc: func(ctx context.Context, req *routingpb.ComputeRoutesRequest, opts ...gax.CallOption) (*routingpb.ComputeRoutesResponse, error) {
			return &routingpb.ComputeRoutesResponse{
				Routes: []*routingpb.Route{
					{
						Duration: durationpb.New(600 * time.Second),
					},
				},
			}, nil
		},
	}

	service := &MapsRoutingService{client: fakeClient}
	origin := &latlng.LatLng{Latitude: 51.503, Longitude: -0.1276}
	destination := &latlng.LatLng{Latitude: 51.498, Longitude: -0.1246}

	// When
	duration, err := service.FetchCurrentTransitTimeBetween(origin, destination)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Then
	expected := 600 * time.Second
	if duration != expected {
		t.Errorf("expected duration %v, got %v", expected, duration)
	}
}

func TestFetchCurrentTransitTimeBetween_ClientError(t *testing.T) {
	// Given
	fakeClient := &fakeRoutesClient{
		computeRoutesFunc: func(ctx context.Context, req *routingpb.ComputeRoutesRequest, opts ...gax.CallOption) (*routingpb.ComputeRoutesResponse, error) {
			return nil, errors.New("client error")
		},
	}

	service := &MapsRoutingService{client: fakeClient}
	origin := &latlng.LatLng{Latitude: 51.503, Longitude: -0.1276}
	destination := &latlng.LatLng{Latitude: 51.498, Longitude: -0.1246}

	// When
	_, err := service.FetchCurrentTransitTimeBetween(origin, destination)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Then
	if !strings.Contains(err.Error(), "client error") {
		t.Errorf("expected error to contain 'client error', got %q", err.Error())
	}
}

func TestFetchCurrentTransitTimeBetween_NoRoutesFound(t *testing.T) {
	// Given
	fakeClient := &fakeRoutesClient{
		computeRoutesFunc: func(ctx context.Context, req *routingpb.ComputeRoutesRequest, opts ...gax.CallOption) (*routingpb.ComputeRoutesResponse, error) {
			return &routingpb.ComputeRoutesResponse{
				Routes: []*routingpb.Route{},
			}, nil
		},
	}

	service := &MapsRoutingService{client: fakeClient}
	origin := &latlng.LatLng{Latitude: 51.503, Longitude: -0.1276}
	destination := &latlng.LatLng{Latitude: 51.498, Longitude: -0.1246}

	// When
	_, err := service.FetchCurrentTransitTimeBetween(origin, destination)
	if err == nil {
		t.Fatal("expected error for no routes found, got nil")
	}

	// Then
	expectedError := "no routes found"
	if err.Error() != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, err.Error())
	}
}
