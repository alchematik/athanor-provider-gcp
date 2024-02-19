package api_gateway

import (
	"context"
	"fmt"
	"regexp"

	gcpapigateway "cloud.google.com/go/apigateway/apiv1"
	apigateway "github.com/alchematik/athanor-provider-gcp/gen/provider/api_gateway"
	"github.com/alchematik/athanor-provider-gcp/gen/provider/identifier"
	"github.com/googleapis/gax-go/v2"

	"cloud.google.com/go/apigateway/apiv1/apigatewaypb"
	sdkerrors "github.com/alchematik/athanor-go/sdk/errors"
	"github.com/alchematik/athanor-go/sdk/provider/value"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
)

var (
	apiConfigRe = regexp.MustCompile(`projects\/(.+)\/locations\/global\/apis\/(.*)\/configs\/(.*)`)
)

func NewHandler(ctx context.Context) (*apigateway.ApiGatewayHandler, error) {
	gcp, err := gcpapigateway.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	c := &client{
		GCP: gcp,
	}

	return &apigateway.ApiGatewayHandler{
		ApiGatewayGetter:  c,
		ApiGatewayUpdator: c,
		ApiGatewayCreator: c,
		ApiGatewayDeleter: c,
	}, nil
}

type client struct {
	GCP GCP
}

type GCP interface {
	CreateGateway(ctx context.Context, req *apigatewaypb.CreateGatewayRequest, opts ...gax.CallOption) (*gcpapigateway.CreateGatewayOperation, error)
	DeleteGateway(ctx context.Context, req *apigatewaypb.DeleteGatewayRequest, opts ...gax.CallOption) (*gcpapigateway.DeleteGatewayOperation, error)
	GetGateway(ctx context.Context, req *apigatewaypb.GetGatewayRequest, opts ...gax.CallOption) (*apigatewaypb.Gateway, error)
	UpdateGateway(ctx context.Context, req *apigatewaypb.UpdateGatewayRequest, opts ...gax.CallOption) (*gcpapigateway.UpdateGatewayOperation, error)
}

func (c *client) GetApiGateway(ctx context.Context, id identifier.ApiGatewayIdentifier) (apigateway.ApiGateway, error) {
	res, err := c.GCP.GetGateway(ctx, &apigatewaypb.GetGatewayRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/gateways/%s", id.Project, id.Location, id.GatewayId),
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return apigateway.ApiGateway{}, sdkerrors.NewErrorNotFound()
		}
		return apigateway.ApiGateway{}, err
	}

	matches := apiConfigRe.FindStringSubmatch(res.ApiConfig)
	if len(matches) < 4 {
		return apigateway.ApiGateway{}, fmt.Errorf("invalid API config ID in response: %q", res.ApiConfig)
	}

	gw := apigateway.ApiGateway{
		Identifier: id,
		Config: apigateway.Config{
			ApiConfig: identifier.ApiConfigIdentifier{
				Api: identifier.ApiIdentifier{
					ApiId: matches[2],
					// TODO: is it safe to assume that the project is the same?
					Project: id.Project,
					// The project ID returned in the response is
					// Project: matches[1],
				},
				ApiConfigId: matches[3],
			},
			DisplayName: res.DisplayName,
			Labels:      res.Labels,
		},
		Attrs: apigateway.Attrs{
			Create:          res.CreateTime.String(),
			Update:          res.UpdateTime.String(),
			State:           res.State.String(),
			DefaultHostname: res.DefaultHostname,
		},
	}

	return gw, nil
}

func (c *client) CreateApiGateway(ctx context.Context, id identifier.ApiGatewayIdentifier, config apigateway.Config) (apigateway.ApiGateway, error) {
	apiConfigID, ok := config.ApiConfig.(identifier.ApiConfigIdentifier)
	if !ok {
		return apigateway.ApiGateway{}, fmt.Errorf("expected API config identifier for api_config, got %T", config.ApiConfig)
	}

	apiID, ok := apiConfigID.Api.(identifier.ApiIdentifier)
	if !ok {
		return apigateway.ApiGateway{}, fmt.Errorf("expected API identifier for api_config.api, got %T", config.ApiConfig)
	}

	op, err := c.GCP.CreateGateway(ctx, &apigatewaypb.CreateGatewayRequest{
		Parent:    fmt.Sprintf("projects/%s/locations/%s", id.Project, id.Location),
		GatewayId: id.GatewayId,
		Gateway: &apigatewaypb.Gateway{
			Labels:      config.Labels,
			DisplayName: config.DisplayName,
			ApiConfig:   fmt.Sprintf("projects/%s/locations/global/apis/%s/configs/%s", apiID.Project, apiID.ApiId, apiConfigID.ApiConfigId),
		},
	})
	if err != nil {
		return apigateway.ApiGateway{}, err
	}

	res, err := op.Wait(ctx)
	if err != nil {
		return apigateway.ApiGateway{}, err
	}

	matches := apiConfigRe.FindStringSubmatch(res.ApiConfig)
	if len(matches) < 4 {
		return apigateway.ApiGateway{}, fmt.Errorf("invalid API config ID in response: %q", res.ApiConfig)
	}

	return apigateway.ApiGateway{
		Identifier: id,
		Config: apigateway.Config{
			ApiConfig: identifier.ApiConfigIdentifier{
				Api: identifier.ApiIdentifier{
					ApiId:   matches[2],
					Project: matches[1],
				},
				ApiConfigId: matches[3],
			},
			DisplayName: res.DisplayName,
			Labels:      res.Labels,
		},
		Attrs: apigateway.Attrs{
			Create:          res.CreateTime.String(),
			Update:          res.UpdateTime.String(),
			State:           res.State.String(),
			DefaultHostname: res.DefaultHostname,
		},
	}, nil
}

func (c *client) UpdateApiGateway(ctx context.Context, id identifier.ApiGatewayIdentifier, config apigateway.Config, mask []value.UpdateMaskField) (apigateway.ApiGateway, error) {
	apiConfigID, ok := config.ApiConfig.(identifier.ApiConfigIdentifier)
	if !ok {
		return apigateway.ApiGateway{}, fmt.Errorf("expected API config identifier for api_config, got %T", config.ApiConfig)
	}

	apiID, ok := apiConfigID.Api.(identifier.ApiIdentifier)
	if !ok {
		return apigateway.ApiGateway{}, fmt.Errorf("expected API identifier for api_config.api, got %T", config.ApiConfig)
	}

	updateMask := &fieldmaskpb.FieldMask{}
	for _, m := range mask {
		switch m.Name {
		case "labels":
			updateMask.Paths = append(updateMask.Paths, "labels")
		case "display_name":
			updateMask.Paths = append(updateMask.Paths, "display_name")
		case "api_config":
			updateMask.Paths = append(updateMask.Paths, "api_config")
		}
	}

	op, err := c.GCP.UpdateGateway(ctx, &apigatewaypb.UpdateGatewayRequest{
		UpdateMask: updateMask,
		Gateway: &apigatewaypb.Gateway{
			Name:        fmt.Sprintf("projects/%s/locations/%s/gateways/%s", id.Project, id.Location, id.GatewayId),
			Labels:      config.Labels,
			DisplayName: config.DisplayName,
			ApiConfig:   fmt.Sprintf("projects/%s/locations/global/apis/%s/configs/%s", apiID.Project, apiID.ApiId, apiConfigID.ApiConfigId),
		},
	})
	if err != nil {
		return apigateway.ApiGateway{}, err
	}

	res, err := op.Wait(ctx)
	if err != nil {
		return apigateway.ApiGateway{}, err
	}

	matches := apiConfigRe.FindStringSubmatch(res.ApiConfig)
	if len(matches) < 4 {
		return apigateway.ApiGateway{}, fmt.Errorf("invalid API config ID in response: %q", res.ApiConfig)
	}

	return apigateway.ApiGateway{
		Identifier: id,
		Config: apigateway.Config{
			ApiConfig: identifier.ApiConfigIdentifier{
				Api: identifier.ApiIdentifier{
					ApiId:   matches[2],
					Project: matches[1],
				},
				ApiConfigId: matches[3],
			},
			DisplayName: res.DisplayName,
			Labels:      res.Labels,
		},
		Attrs: apigateway.Attrs{
			Create:          res.CreateTime.String(),
			Update:          res.UpdateTime.String(),
			State:           res.State.String(),
			DefaultHostname: res.DefaultHostname,
		},
	}, nil
}

func (c *client) DeleteApiGateway(ctx context.Context, id identifier.ApiGatewayIdentifier) error {
	op, err := c.GCP.DeleteGateway(ctx, &apigatewaypb.DeleteGatewayRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/gateways/%s", id.Project, id.Location, id.GatewayId),
	})
	if err != nil {
		return err
	}

	return op.Wait(ctx)
}
