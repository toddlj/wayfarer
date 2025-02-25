package googlemaps

import (
	routing "cloud.google.com/go/maps/routing/apiv2"
	"cloud.google.com/go/maps/routing/apiv2/routingpb"
	"context"
	"errors"
	"fmt"
	"github.com/googleapis/gax-go/v2"
	"github.com/googleapis/gax-go/v2/callctx"
	"google.golang.org/genproto/googleapis/type/latlng"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"

	"google.golang.org/api/option"
)

type RoutesClient interface {
	ComputeRoutes(ctx context.Context, req *routingpb.ComputeRoutesRequest, opts ...gax.CallOption) (*routingpb.ComputeRoutesResponse, error)
	Close() error
}

type MapsRoutingService struct {
	client RoutesClient
}

func NewMapsRoutingService(googleApiBaseUrl string, googleApiKey string) (*MapsRoutingService, error) {
	if googleApiBaseUrl != "" {
		slog.Warn("Using insecure connection to custom Google Maps API", slog.String("url", googleApiBaseUrl))
		client, err := routing.NewRoutesClient(context.Background(),
			option.WithEndpoint(googleApiBaseUrl),
			option.WithoutAuthentication(),
			option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))
		if err != nil {
			return nil, fmt.Errorf("failed to create Routes client: %w", err)
		}

		return &MapsRoutingService{client: client}, nil
	}

	client, err := routing.NewRoutesClient(context.Background(), option.WithAPIKey(googleApiKey))
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

	ctx := callctx.SetHeaders(context.Background(), callctx.XGoogFieldMaskHeader, "routes.duration")
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
