package googlemaps

import (
	routing "cloud.google.com/go/maps/routing/apiv2"
	"cloud.google.com/go/maps/routing/apiv2/routingpb"
	"context"
	"errors"
	"fmt"
	"github.com/googleapis/gax-go/v2/callctx"
	"google.golang.org/genproto/googleapis/type/latlng"
	"os"
	"time"

	"google.golang.org/api/option"
)

type MapsRoutingService struct {
	client *routing.RoutesClient
}

func NewMapsRoutingService() (*MapsRoutingService, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		return nil, errors.New("GOOGLE_API_KEY environment variable must be set")
	}

	client, err := routing.NewRoutesClient(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Routes client: %w", err)
	}

	return &MapsRoutingService{client: client}, nil
}

func (s *MapsRoutingService) Close() error {
	return s.client.Close()
}

func (s *MapsRoutingService) FetchCurrentTransitTimeBetween(origin, destination *latlng.LatLng) (time.Duration, error) {
	req := &routingpb.ComputeRoutesRequest{
		Origin:      &routingpb.Waypoint{LocationType: &routingpb.Waypoint_Location{Location: &routingpb.Location{LatLng: origin}}},
		Destination: &routingpb.Waypoint{LocationType: &routingpb.Waypoint_Location{Location: &routingpb.Location{LatLng: destination}}},
		TravelMode:  routingpb.RouteTravelMode_TRANSIT,
	}

	ctx := callctx.SetHeaders(context.Background(), "x-goog-fieldmask", "routes.duration")
	resp, err := s.client.ComputeRoutes(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("API request to compute routes failed: %w", err)
	}

	if len(resp.Routes) == 0 {
		return 0, errors.New("no routes found")
	}

	duration := resp.Routes[0].Duration
	return time.Duration(duration.Seconds) * time.Second, nil
}
