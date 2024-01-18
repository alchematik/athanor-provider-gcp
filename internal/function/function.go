package function

import (
	"context"

	"github.com/alchematik/athanor-provider-gcp/gen/provider/function"
	"github.com/alchematik/athanor-provider-gcp/gen/provider/identifier"

	cloudfunction "cloud.google.com/go/functions/apiv2"
	"cloud.google.com/go/functions/apiv2/functionspb"
	sdkerrors "github.com/alchematik/athanor-go/sdk/errors"
	"github.com/alchematik/athanor-go/sdk/provider/value"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewHandler() function.FunctionHandler {
	c := &client{}
	return function.FunctionHandler{
		FunctionGetter:  c,
		FunctionCreator: c,
		FunctionUpdator: c,
		FunctionDeleter: c,
	}
}

type client struct {
}

func (c *client) GetFunction(ctx context.Context, id identifier.FunctionIdentifier) (function.Function, error) {
	gcp, err := cloudfunction.NewFunctionClient(ctx)
	if err != nil {
		return function.Function{}, nil
	}

	req := &functionspb.GetFunctionRequest{}
	res, err := gcp.GetFunction(ctx, req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return function.Function{}, sdkerrors.NewErrorNotFound()
		}

		return function.Function{}, nil
	}

	labels := map[string]any{}
	for k, v := range res.GetLabels() {
		labels[k] = v
	}

	// storageSource := res.GetBuildConfig().GetSource().GetStorageSource()
	// bucketName := storageSource.GetBucket()
	// objectName := storageSource.GetObject()

	return function.Function{
		Identifier: id,
		Config: function.FunctionConfig{
			Description: res.GetDescription(),
			Labels:      labels,
			BuildConfig: function.BuildConfig{
				Runtime:    res.GetBuildConfig().GetRuntime(),
				Entrypoint: res.GetBuildConfig().GetEntryPoint(),
				Source:     identifier.BucketObjectIdentifier{},
			},
		},
		Attrs: function.FunctionAttrs{},
	}, nil
}

func (c *client) CreateFunction(ctx context.Context, id identifier.FunctionIdentifier, config function.FunctionConfig) (function.Function, error) {
	return function.Function{}, nil
}

func (c *client) UpdateFunction(ctx context.Context, id identifier.FunctionIdentifier, config function.FunctionConfig, mask []value.UpdateMaskField) (function.Function, error) {
	return function.Function{}, nil
}

func (c *client) DeleteFunction(ctx context.Context, id identifier.FunctionIdentifier) error {
	return nil
}
