package function

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/alchematik/athanor-provider-gcp/gen/provider/function"
	"github.com/alchematik/athanor-provider-gcp/gen/provider/identifier"

	cloudfunction "cloud.google.com/go/functions/apiv2"
	"cloud.google.com/go/functions/apiv2/functionspb"
	"cloud.google.com/go/storage"
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
		return function.Function{}, err
	}

	req := &functionspb.GetFunctionRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/functions/%s", id.Project, id.Location, id.Name),
	}
	res, err := gcp.GetFunction(ctx, req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return function.Function{}, sdkerrors.NewErrorNotFound()
		}

		return function.Function{}, err
	}

	labels := map[string]any{}
	for k, v := range res.GetLabels() {
		labels[k] = v
	}

	storageSource := res.GetBuildConfig().GetSource().GetStorageSource()
	bucketName := storageSource.GetBucket()
	objectName := storageSource.GetObject()

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return function.Function{}, err
	}

	objectAttrs, err := storageClient.Bucket(bucketName).Object(objectName).Attrs(ctx)
	if err != nil {
		return function.Function{}, err
	}

	return function.Function{
		Identifier: id,
		Config: function.FunctionConfig{
			Description: res.GetDescription(),
			Labels:      labels,
			BuildConfig: function.BuildConfig{
				Runtime:    res.GetBuildConfig().GetRuntime(),
				Entrypoint: res.GetBuildConfig().GetEntryPoint(),
				Source: value.File{
					Checksum: fmt.Sprintf("%x", objectAttrs.MD5),
				},
			},
		},
		Attrs: function.FunctionAttrs{
			Url: res.Url,
		},
	}, nil
}

func (c *client) CreateFunction(ctx context.Context, id identifier.FunctionIdentifier, config function.FunctionConfig) (function.Function, error) {
	gcp, err := cloudfunction.NewFunctionClient(ctx)
	if err != nil {
		return function.Function{}, err
	}

	// TODO: Make idempotent by checking if there's an exising create operation.

	uploadURLRes, err := gcp.GenerateUploadUrl(ctx, &functionspb.GenerateUploadUrlRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", id.Project, id.Location),
	})
	if err != nil {
		return function.Function{}, err
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", filepath.Base(config.BuildConfig.Source.Path))
	if err != nil {
		return function.Function{}, err
	}

	file, err := os.Open(config.BuildConfig.Source.Path)
	if err != nil {
		return function.Function{}, err
	}
	defer file.Close()
	if _, err := io.Copy(fw, file); err != nil {
		return function.Function{}, err
	}
	w.Close()

	req, err := http.NewRequest(http.MethodPut, uploadURLRes.UploadUrl, &buf)
	if err != nil {
		return function.Function{}, err
	}
	req.Header.Set("Content-Type", "application/zip")
	client := &http.Client{}
	if _, err = client.Do(req); err != nil {
		return function.Function{}, err
	}

	labels := map[string]string{}
	for k, v := range config.Labels {
		str, ok := v.(string)
		if !ok {
			return function.Function{}, fmt.Errorf("label values must be string, got %T", v)
		}
		labels[k] = str
	}

	// TODO: create the function.
	operation, err := gcp.CreateFunction(ctx, &functionspb.CreateFunctionRequest{
		Parent:     fmt.Sprintf("projects/%s/locations/%s", id.Project, id.Location),
		FunctionId: id.Name,
		Function: &functionspb.Function{
			Name:        fmt.Sprintf("projects/%s/locations/%s/functions/%s", id.Project, id.Location, id.Name),
			Environment: functionspb.Environment_GEN_2,
			Description: config.Description,
			Labels:      labels,
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

	log.Printf("FUNCTION CREATE OPERATION >>>>>>>>> %v\n", operation.Name())

	res, err := operation.Wait(ctx)
	if err != nil {
		return function.Function{}, err
	}

	outLabels := map[string]any{}
	for k, v := range res.Labels {
		outLabels[k] = v
	}

	endStorage := res.GetBuildConfig().GetSource().GetStorageSource()
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return function.Function{}, err
	}
	objectAttrs, err := storageClient.Bucket(endStorage.GetBucket()).Object(endStorage.GetObject()).Attrs(ctx)
	if err != nil {
		return function.Function{}, err
	}

	return function.Function{
		Identifier: id,
		Config: function.FunctionConfig{
			Description: res.Description,
			Labels:      outLabels,
			BuildConfig: function.BuildConfig{
				Runtime:    res.GetBuildConfig().GetRuntime(),
				Entrypoint: res.GetBuildConfig().GetEntryPoint(),
				Source: value.File{
					Checksum: fmt.Sprintf("%x", objectAttrs.MD5),
				},
			},
		},
		Attrs: function.FunctionAttrs{
			Url: res.Url,
		},
	}, nil
}

func (c *client) UpdateFunction(ctx context.Context, id identifier.FunctionIdentifier, config function.FunctionConfig, mask []value.UpdateMaskField) (function.Function, error) {
	gcp, err := cloudfunction.NewFunctionClient(ctx)
	if err != nil {
		return function.Function{}, err
	}

	// TODO: Make idempotent by checking if there's an exising update operation.

	labels := map[string]string{}
	for k, v := range config.Labels {
		str, ok := v.(string)
		if !ok {
			return function.Function{}, fmt.Errorf("label values must be string, got %T", v)
		}
		labels[k] = str
	}

	uploadURLRes, err := gcp.GenerateUploadUrl(ctx, &functionspb.GenerateUploadUrlRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", id.Project, id.Location),
	})
	if err != nil {
		return function.Function{}, err
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", filepath.Base(config.BuildConfig.Source.Path))
	if err != nil {
		return function.Function{}, err
	}

	file, err := os.Open(config.BuildConfig.Source.Path)
	if err != nil {
		return function.Function{}, err
	}
	defer file.Close()
	if _, err := io.Copy(fw, file); err != nil {
		return function.Function{}, err
	}
	w.Close()

	req, err := http.NewRequest(http.MethodPut, uploadURLRes.UploadUrl, &buf)
	if err != nil {
		return function.Function{}, err
	}
	req.Header.Set("Content-Type", "application/zip")
	client := &http.Client{}
	if _, err = client.Do(req); err != nil {
		return function.Function{}, err
	}

	operation, err := gcp.UpdateFunction(ctx, &functionspb.UpdateFunctionRequest{
		Function: &functionspb.Function{
			Name:        fmt.Sprintf("projects/%s/locations/%s/functions/%s", id.Project, id.Location, id.Name),
			Environment: functionspb.Environment_GEN_2,
			Description: config.Description,
			Labels:      labels,
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

	res, err := operation.Wait(ctx)
	if err != nil {
		return function.Function{}, err
	}

	outLabels := map[string]any{}
	for k, v := range res.Labels {
		outLabels[k] = v
	}

	endStorage := res.GetBuildConfig().GetSource().GetStorageSource()
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return function.Function{}, err
	}
	objectAttrs, err := storageClient.Bucket(endStorage.GetBucket()).Object(endStorage.GetObject()).Attrs(ctx)
	if err != nil {
		return function.Function{}, err
	}

	return function.Function{
		Identifier: id,
		Config: function.FunctionConfig{
			Description: res.Description,
			Labels:      outLabels,
			BuildConfig: function.BuildConfig{
				Runtime:    res.GetBuildConfig().GetRuntime(),
				Entrypoint: res.GetBuildConfig().GetEntryPoint(),
				Source: value.File{
					Checksum: fmt.Sprintf("%x", objectAttrs.MD5),
				},
			},
		},
		Attrs: function.FunctionAttrs{
			Url: res.Url,
		},
	}, nil
}

func (c *client) DeleteFunction(ctx context.Context, id identifier.FunctionIdentifier) error {
	return nil
}
