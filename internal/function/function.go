package function

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/alchematik/athanor-provider-gcp/gen/provider/function"
	"github.com/alchematik/athanor-provider-gcp/gen/provider/identifier"

	cloudfunction "cloud.google.com/go/functions/apiv2"
	"cloud.google.com/go/functions/apiv2/functionspb"
	"cloud.google.com/go/storage"
	sdkerrors "github.com/alchematik/athanor-go/sdk/errors"
	"github.com/alchematik/athanor-go/sdk/provider/value"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
)

func NewHandler(ctx context.Context) (*function.FunctionHandler, error) {
	gcp, err := cloudfunction.NewFunctionClient(ctx)
	if err != nil {
		return nil, err
	}

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	c := &client{
		GCP:     gcp,
		Storage: storageClient,
	}
	return &function.FunctionHandler{
		FunctionGetter:  c,
		FunctionCreator: c,
		FunctionUpdator: c,
		FunctionDeleter: c,
		CloseFunc: func() error {
			if err := gcp.Close(); err != nil {
				return err
			}

			return storageClient.Close()
		},
	}, nil
}

type client struct {
	GCP     GCP
	Storage Storage
}

type GCP interface {
	GetFunction(context.Context, *functionspb.GetFunctionRequest, ...gax.CallOption) (*functionspb.Function, error)
	GenerateUploadUrl(context.Context, *functionspb.GenerateUploadUrlRequest, ...gax.CallOption) (*functionspb.GenerateUploadUrlResponse, error)
	CreateFunction(context.Context, *functionspb.CreateFunctionRequest, ...gax.CallOption) (*cloudfunction.CreateFunctionOperation, error)
	UpdateFunction(context.Context, *functionspb.UpdateFunctionRequest, ...gax.CallOption) (*cloudfunction.UpdateFunctionOperation, error)
	DeleteFunction(context.Context, *functionspb.DeleteFunctionRequest, ...gax.CallOption) (*cloudfunction.DeleteFunctionOperation, error)
}

type Storage interface {
	Bucket(string) *storage.BucketHandle
}

func (c *client) GetFunction(ctx context.Context, id identifier.FunctionIdentifier) (function.Function, error) {
	req := &functionspb.GetFunctionRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/functions/%s", id.Project, id.Location, id.Name),
	}
	res, err := c.GCP.GetFunction(ctx, req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return function.Function{}, sdkerrors.NewErrorNotFound()
		}

		return function.Function{}, err
	}

	storageSource := res.GetBuildConfig().GetSource().GetStorageSource()
	bucketName := storageSource.GetBucket()
	objectName := storageSource.GetObject()

	objectAttrs, err := c.Storage.Bucket(bucketName).Object(objectName).Attrs(ctx)
	if err != nil {
		return function.Function{}, err
	}

	return function.Function{
		Identifier: id,
		Config: function.Config{
			Description: res.GetDescription(),
			Labels:      res.GetLabels(),
			BuildConfig: function.BuildConfig{
				Runtime:    res.GetBuildConfig().GetRuntime(),
				Entrypoint: res.GetBuildConfig().GetEntryPoint(),
				Source: value.File{
					Checksum: fmt.Sprintf("%d", objectAttrs.CRC32C),
				},
			},
		},
		Attrs: function.Attrs{
			Url: res.Url,
		},
	}, nil
}

func (c *client) CreateFunction(ctx context.Context, id identifier.FunctionIdentifier, config function.Config) (function.Function, error) {
	// TODO: Make idempotent by checking if there's an exising create operation.

	uploadURLRes, err := c.GCP.GenerateUploadUrl(ctx, &functionspb.GenerateUploadUrlRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", id.Project, id.Location),
	})
	if err != nil {
		return function.Function{}, err
	}

	objHandle := c.Storage.Bucket(uploadURLRes.GetStorageSource().GetBucket()).Object(uploadURLRes.GetStorageSource().GetObject())
	writer := objHandle.NewWriter(ctx)
	file, err := os.Open(config.BuildConfig.Source.Path)
	if err != nil {
		return function.Function{}, err
	}
	if _, err := io.Copy(writer, file); err != nil {
		return function.Function{}, err
	}
	if err := writer.Close(); err != nil {
		return function.Function{}, err
	}

	operation, err := c.GCP.CreateFunction(ctx, &functionspb.CreateFunctionRequest{
		Parent:     fmt.Sprintf("projects/%s/locations/%s", id.Project, id.Location),
		FunctionId: id.Name,
		Function: &functionspb.Function{
			Name:        fmt.Sprintf("projects/%s/locations/%s/functions/%s", id.Project, id.Location, id.Name),
			Environment: functionspb.Environment_GEN_2,
			Description: config.Description,
			Labels:      config.Labels,
			BuildConfig: &functionspb.BuildConfig{
				Runtime:    config.BuildConfig.Runtime,
				EntryPoint: config.BuildConfig.Entrypoint,
				Source: &functionspb.Source{
					Source: &functionspb.Source_StorageSource{
						StorageSource: uploadURLRes.StorageSource,
					},
				},
			},
		},
	})
	if err != nil {
		return function.Function{}, err
	}

	res, err := operation.Wait(ctx)
	if err != nil {
		return function.Function{}, err
	}

	endStorage := res.GetBuildConfig().GetSource().GetStorageSource()
	objectAttrs, err := c.Storage.Bucket(endStorage.GetBucket()).Object(endStorage.GetObject()).Attrs(ctx)
	if err != nil {
		return function.Function{}, err
	}

	return function.Function{
		Identifier: id,
		Config: function.Config{
			Description: res.Description,
			Labels:      res.Labels,
			BuildConfig: function.BuildConfig{
				Runtime:    res.GetBuildConfig().GetRuntime(),
				Entrypoint: res.GetBuildConfig().GetEntryPoint(),
				Source: value.File{
					Checksum: fmt.Sprintf("%d", objectAttrs.CRC32C),
				},
			},
		},
		Attrs: function.Attrs{
			Url: res.Url,
		},
	}, nil
}

func (c *client) UpdateFunction(ctx context.Context, id identifier.FunctionIdentifier, config function.Config, mask []value.UpdateMaskField) (function.Function, error) {

	// TODO: Make idempotent by checking if there's an exising update operation.

	updateMask := fieldmaskpb.FieldMask{}
	updateFunc := &functionspb.Function{
		Name: fmt.Sprintf("projects/%s/locations/%s/functions/%s", id.Project, id.Location, id.Name),
		// Environment: functionspb.Environment_GEN_2,
	}
	for _, m := range mask {
		switch m.Name {
		case "labels":
			updateFunc.Labels = config.Labels
			updateMask.Paths = append(updateMask.Paths, "labels")
		case "description":
			updateFunc.Description = config.Description
			updateMask.Paths = append(updateMask.Paths, "description")
		case "build_config":
			bc := functionspb.BuildConfig{}
			for _, f := range m.SubFields {
				switch f.Name {
				case "runtime":
					bc.Runtime = config.BuildConfig.Runtime
					updateMask.Paths = append(updateMask.Paths, "build_config.runtime")
				case "entrypoint":
					bc.EntryPoint = config.BuildConfig.Entrypoint
					updateMask.Paths = append(updateMask.Paths, "build_config.entrypoint")
				case "source":
					uploadURLRes, err := c.GCP.GenerateUploadUrl(ctx, &functionspb.GenerateUploadUrlRequest{
						Parent: fmt.Sprintf("projects/%s/locations/%s", id.Project, id.Location),
					})
					if err != nil {
						return function.Function{}, err
					}

					objHandle := c.Storage.Bucket(uploadURLRes.GetStorageSource().GetBucket()).Object(uploadURLRes.GetStorageSource().GetObject())
					writer := objHandle.NewWriter(ctx)
					file, err := os.Open(config.BuildConfig.Source.Path)
					if err != nil {
						return function.Function{}, err
					}
					if _, err := io.Copy(writer, file); err != nil {
						return function.Function{}, err
					}
					if err := writer.Close(); err != nil {
						return function.Function{}, err
					}

					bc.Source = &functionspb.Source{
						Source: &functionspb.Source_StorageSource{
							StorageSource: uploadURLRes.StorageSource,
						},
					}
					updateMask.Paths = append(updateMask.Paths, "build_config.source")
				}
			}
			updateFunc.BuildConfig = &bc
		}
	}

	operation, err := c.GCP.UpdateFunction(ctx, &functionspb.UpdateFunctionRequest{
		Function:   updateFunc,
		UpdateMask: &updateMask,
	})
	if err != nil {
		return function.Function{}, err
	}

	res, err := operation.Wait(ctx)
	if err != nil {
		return function.Function{}, err
	}

	endStorage := res.GetBuildConfig().GetSource().GetStorageSource()
	objectAttrs, err := c.Storage.Bucket(endStorage.GetBucket()).Object(endStorage.GetObject()).Attrs(ctx)
	if err != nil {
		return function.Function{}, err
	}

	return function.Function{
		Identifier: id,
		Config: function.Config{
			Description: res.Description,
			Labels:      res.Labels,
			BuildConfig: function.BuildConfig{
				Runtime:    res.GetBuildConfig().GetRuntime(),
				Entrypoint: res.GetBuildConfig().GetEntryPoint(),
				Source: value.File{
					Checksum: fmt.Sprintf("%d", objectAttrs.CRC32C),
				},
			},
		},
		Attrs: function.Attrs{
			Url: res.Url,
		},
	}, nil
}

func (c *client) DeleteFunction(ctx context.Context, id identifier.FunctionIdentifier) error {
	operation, err := c.GCP.DeleteFunction(ctx, &functionspb.DeleteFunctionRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/functions/%s", id.Project, id.Location, id.Name),
	})
	if err != nil {
		return err
	}

	return operation.Wait(ctx)
}
