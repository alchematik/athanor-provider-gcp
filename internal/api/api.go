package api

import (
	"context"
	"fmt"

	"github.com/alchematik/athanor-provider-gcp/gen/provider/api"
	"github.com/alchematik/athanor-provider-gcp/gen/provider/identifier"

	"cloud.google.com/go/apigateway/apiv1"
	"cloud.google.com/go/apigateway/apiv1/apigatewaypb"
	sdkerrors "github.com/alchematik/athanor-go/sdk/errors"
	"github.com/alchematik/athanor-go/sdk/provider/value"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
)

func NewHandler(ctx context.Context) (*api.ApiHandler, error) {
	gcp, err := apigateway.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	c := &client{
		GCP: gcp,
	}
	return &api.ApiHandler{
		ApiGetter:  c,
		ApiCreator: c,
		ApiUpdator: c,
		ApiDeleter: c,
	}, nil
}

type client struct {
	GCP GCP
}

type GCP interface {
	// CreateApiConfig(ctx context.Context, req *apigatewaypb.CreateApiConfigRequest, opts ...gax.CallOption) (*apigateway.CreateApiConfigOperation, error)
	// DeleteApiConfig(ctx context.Context, req *apigatewaypb.DeleteApiConfigRequest, opts ...gax.CallOption) (*apigateway.DeleteApiConfigOperation, error)
	// GetApiConfig(ctx context.Context, req *apigatewaypb.GetApiConfigRequest, opts ...gax.CallOption) (*apigatewaypb.ApiConfig, error)
	// CreateGateway(ctx context.Context, req *apigatewaypb.CreateGatewayRequest, opts ...gax.CallOption) (*apigateway.CreateGatewayOperation, error)
	// DeleteGateway(ctx context.Context, req *apigatewaypb.DeleteGatewayRequest, opts ...gax.CallOption) (*apigateway.DeleteGatewayOperation, error)
	// GetGateway(ctx context.Context, req *apigatewaypb.GetGatewayRequest, opts ...gax.CallOption) (*apigatewaypb.Gateway, error)
	CreateApi(ctx context.Context, req *apigatewaypb.CreateApiRequest, opts ...gax.CallOption) (*apigateway.CreateApiOperation, error)
	DeleteApi(ctx context.Context, req *apigatewaypb.DeleteApiRequest, opts ...gax.CallOption) (*apigateway.DeleteApiOperation, error)
	GetApi(ctx context.Context, req *apigatewaypb.GetApiRequest, opts ...gax.CallOption) (*apigatewaypb.Api, error)
	UpdateApi(ctx context.Context, req *apigatewaypb.UpdateApiRequest, opts ...gax.CallOption) (*apigateway.UpdateApiOperation, error)
}

func (c *client) GetApi(ctx context.Context, id identifier.ApiIdentifier) (api.Api, error) {
	res, err := c.GCP.GetApi(ctx, &apigatewaypb.GetApiRequest{
		Name: fmt.Sprintf("projects/%s/locations/global/apis/%s", id.Project, id.ApiId),
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return api.Api{}, sdkerrors.NewErrorNotFound()
		}

		return api.Api{}, err
	}

	return api.Api{
		Identifier: id,
		Config: api.Config{
			DisplayName: res.GetDisplayName(),
			Labels:      res.GetLabels(),
		},
		Attrs: api.Attrs{
			Create: res.GetCreateTime().String(),
			Update: res.GetUpdateTime().String(),
			State:  res.GetState().String(),
		},
	}, nil
}

func (c *client) CreateApi(ctx context.Context, id identifier.ApiIdentifier, config api.Config) (api.Api, error) {
	// TODO: make idempotent by getting active operation if exists.
	op, err := c.GCP.CreateApi(ctx, &apigatewaypb.CreateApiRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", id.Project),
	})
	if err != nil {
		return api.Api{}, err
	}

	res, err := op.Wait(ctx)
	if err != nil {
		return api.Api{}, err
	}

	return api.Api{
		Identifier: id,
		Config: api.Config{
			DisplayName: res.GetDisplayName(),
			Labels:      res.GetLabels(),
		},
		Attrs: api.Attrs{
			Create: res.GetCreateTime().String(),
			Update: res.GetUpdateTime().String(),
			State:  res.GetState().String(),
		},
	}, nil
}

func (c *client) UpdateApi(ctx context.Context, id identifier.ApiIdentifier, config api.Config, mask []value.UpdateMaskField) (api.Api, error) {
	updateMask := &fieldmaskpb.FieldMask{}
	object := &apigatewaypb.Api{
		Name: fmt.Sprintf("projects/%s/locations/global/apis/%s", id.Project, id.ApiId),
	}
	for _, m := range mask {
		switch m.Name {
		case "display_name":
			object.DisplayName = config.DisplayName
			updateMask.Paths = append(updateMask.Paths, "display_name")
		case "labels":
			object.Labels = config.Labels
			updateMask.Paths = append(updateMask.Paths, "labels")
		}
	}

	op, err := c.GCP.UpdateApi(ctx, &apigatewaypb.UpdateApiRequest{
		UpdateMask: updateMask,
		Api:        object,
	})
	if err != nil {
		return api.Api{}, err
	}

	res, err := op.Wait(ctx)
	if err != nil {
		return api.Api{}, err
	}

	return api.Api{
		Identifier: id,
		Config: api.Config{
			DisplayName: res.GetDisplayName(),
			Labels:      res.GetLabels(),
		},
		Attrs: api.Attrs{
			Create: res.GetCreateTime().String(),
			Update: res.GetUpdateTime().String(),
			State:  res.GetState().String(),
		},
	}, nil
}

func (c *client) DeleteApi(ctx context.Context, id identifier.ApiIdentifier) error {
	_, err := c.GCP.DeleteApi(ctx, &apigatewaypb.DeleteApiRequest{
		Name: fmt.Sprintf("projects/%s/locations/global/apis/%s", id.Project, id.ApiId),
	})
	return err
}
