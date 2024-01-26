package api_config

import (
	"context"
	"fmt"
	"hash/crc32"
	"os"

	apiconfig "github.com/alchematik/athanor-provider-gcp/gen/provider/api_config"
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

func NewHandler(ctx context.Context) (*apiconfig.ApiConfigHandler, error) {
	gcp, err := apigateway.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	c := &client{
		GCP: gcp,
	}
	return &apiconfig.ApiConfigHandler{
		ApiConfigGetter:  c,
		ApiConfigCreator: c,
		ApiConfigUpdator: c,
		ApiConfigDeleter: c,
		CloseFunc:        gcp.Close,
	}, nil
}

type client struct {
	GCP GCP
}

type GCP interface {
	CreateApiConfig(ctx context.Context, req *apigatewaypb.CreateApiConfigRequest, opts ...gax.CallOption) (*apigateway.CreateApiConfigOperation, error)
	DeleteApiConfig(ctx context.Context, req *apigatewaypb.DeleteApiConfigRequest, opts ...gax.CallOption) (*apigateway.DeleteApiConfigOperation, error)
	GetApiConfig(ctx context.Context, req *apigatewaypb.GetApiConfigRequest, opts ...gax.CallOption) (*apigatewaypb.ApiConfig, error)
	UpdateApiConfig(ctx context.Context, req *apigatewaypb.UpdateApiConfigRequest, opts ...gax.CallOption) (*apigateway.UpdateApiConfigOperation, error)
}

func (c *client) GetApiConfig(ctx context.Context, id identifier.ApiConfigIdentifier) (apiconfig.ApiConfig, error) {
	apiID, ok := id.Api.(identifier.ApiIdentifier)
	if !ok {
		return apiconfig.ApiConfig{}, fmt.Errorf("field api must be an api identifier")
	}

	res, err := c.GCP.GetApiConfig(ctx, &apigatewaypb.GetApiConfigRequest{
		Name: fmt.Sprintf("projects/%s/locations/global/apis/%s/configs/%s", apiID.Project, apiID.ApiId, id.ApiConfigId),
		View: apigatewaypb.GetApiConfigRequest_FULL,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return apiconfig.ApiConfig{}, sdkerrors.NewErrorNotFound()
		}

		return apiconfig.ApiConfig{}, err
	}

	files := make([]value.File, len(res.GetOpenapiDocuments()))
	for i, doc := range res.GetOpenapiDocuments() {
		checksum := crc32.Checksum(doc.GetDocument().GetContents(), crc32.MakeTable(crc32.Castagnoli))
		files[i] = value.File{
			Path:     doc.GetDocument().GetPath(),
			Checksum: fmt.Sprintf("%d", checksum),
		}
	}

	return apiconfig.ApiConfig{
		Identifier: id,
		Config: apiconfig.Config{
			DisplayName:      res.GetDisplayName(),
			OpenApiDocuments: files,
		},
		Attrs: apiconfig.Attrs{
			Create: res.GetCreateTime().String(),
			Update: res.GetUpdateTime().String(),
			State:  res.GetState().String(),
		},
	}, nil
}

func (c *client) CreateApiConfig(ctx context.Context, id identifier.ApiConfigIdentifier, config apiconfig.Config) (apiconfig.ApiConfig, error) {
	apiID, ok := id.Api.(identifier.ApiIdentifier)
	if !ok {
		return apiconfig.ApiConfig{}, fmt.Errorf("field api must be an api identifier")
	}

	serviceAccountID, ok := id.ServiceAccount.(identifier.ServiceAccountIdentifier)
	if !ok {
		return apiconfig.ApiConfig{}, fmt.Errorf("field service_account must be an api identifier")
	}

	docs := make([]*apigatewaypb.ApiConfig_OpenApiDocument, len(config.OpenApiDocuments))
	for i, doc := range config.OpenApiDocuments {
		data, err := os.ReadFile(doc.Path)
		if err != nil {
			return apiconfig.ApiConfig{}, err
		}

		docs[i] = &apigatewaypb.ApiConfig_OpenApiDocument{
			Document: &apigatewaypb.ApiConfig_File{
				Path:     doc.Path,
				Contents: data,
			},
		}
	}

	apiConfig := &apigatewaypb.ApiConfig{
		DisplayName:           config.DisplayName,
		GatewayServiceAccount: fmt.Sprintf("%s@%s.iam.gserviceaccount.com", serviceAccountID.AccountId, serviceAccountID.Project),
		OpenapiDocuments:      docs,
	}
	op, err := c.GCP.CreateApiConfig(ctx, &apigatewaypb.CreateApiConfigRequest{
		Parent:      fmt.Sprintf("projects/%s/locations/global/apis/%s", apiID.Project, apiID.ApiId),
		ApiConfigId: id.ApiConfigId,
		ApiConfig:   apiConfig,
	})
	if err != nil {
		return apiconfig.ApiConfig{}, err
	}

	res, err := op.Wait(ctx)
	if err != nil {
		return apiconfig.ApiConfig{}, err
	}

	files := make([]value.File, len(res.GetOpenapiDocuments()))
	for i, doc := range res.GetOpenapiDocuments() {
		checksum := crc32.Checksum(doc.GetDocument().GetContents(), crc32.MakeTable(crc32.Castagnoli))
		files[i] = value.File{
			Path:     doc.GetDocument().GetPath(),
			Checksum: fmt.Sprintf("%d", checksum),
		}
	}

	return apiconfig.ApiConfig{
		Identifier: id,
		Config: apiconfig.Config{
			DisplayName:      res.GetDisplayName(),
			OpenApiDocuments: files,
		},
		Attrs: apiconfig.Attrs{
			Create: res.GetCreateTime().String(),
			Update: res.GetUpdateTime().String(),
			State:  res.GetState().String(),
		},
	}, nil
}

func (c *client) UpdateApiConfig(ctx context.Context, id identifier.ApiConfigIdentifier, config apiconfig.Config, mask []value.UpdateMaskField) (apiconfig.ApiConfig, error) {
	apiID, ok := id.Api.(identifier.ApiIdentifier)
	if !ok {
		return apiconfig.ApiConfig{}, fmt.Errorf("field api must be an api identifier")
	}

	apiConfig := &apigatewaypb.ApiConfig{
		Name: fmt.Sprintf("projects/%s/locations/global/apis/%s/configs/%s", apiID.Project, apiID.ApiId, id.ApiConfigId),
	}
	updateMask := &fieldmaskpb.FieldMask{}

	for _, m := range mask {
		switch m.Name {
		case "display_name":
			apiConfig.DisplayName = config.DisplayName
			updateMask.Paths = append(updateMask.Paths, "display_name")
		}
	}

	req := &apigatewaypb.UpdateApiConfigRequest{
		UpdateMask: updateMask,
		ApiConfig:  apiConfig,
	}

	op, err := c.GCP.UpdateApiConfig(ctx, req)
	if err != nil {
		return apiconfig.ApiConfig{}, err
	}
	res, err := op.Wait(ctx)
	if err != nil {
		return apiconfig.ApiConfig{}, err
	}

	files := make([]value.File, len(res.GetOpenapiDocuments()))
	for i, doc := range res.GetOpenapiDocuments() {
		checksum := crc32.Checksum(doc.GetDocument().GetContents(), crc32.MakeTable(crc32.Castagnoli))
		files[i] = value.File{
			Path:     doc.GetDocument().GetPath(),
			Checksum: fmt.Sprintf("%d", checksum),
		}
	}

	return apiconfig.ApiConfig{
		Identifier: id,
		Config: apiconfig.Config{
			DisplayName:      res.GetDisplayName(),
			OpenApiDocuments: files,
		},
		Attrs: apiconfig.Attrs{
			Create: res.GetCreateTime().String(),
			Update: res.GetUpdateTime().String(),
			State:  res.GetState().String(),
		},
	}, nil
}

func (c *client) DeleteApiConfig(ctx context.Context, id identifier.ApiConfigIdentifier) error {
	apiID, ok := id.Api.(identifier.ApiIdentifier)
	if !ok {
		return fmt.Errorf("field api must be an api identifier")
	}
	op, err := c.GCP.DeleteApiConfig(ctx, &apigatewaypb.DeleteApiConfigRequest{
		Name: fmt.Sprintf("projects/%s/locations/global/apis/%s/configs/%s", apiID.Project, apiID.ApiId, id.ApiConfigId),
	})
	if err != nil {
		return err
	}

	return op.Wait(ctx)
}
