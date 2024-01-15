package bucket

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/alchematik/athanor-provider-gcp/gen/provider"

	"cloud.google.com/go/storage"
	sdkerrors "github.com/alchematik/athanor-go/sdk/errors"
	value "github.com/alchematik/athanor-go/sdk/provider/value"
)

func NewHandler() provider.BucketHandler {
	c := &client{}
	return provider.BucketHandler{
		BucketGetter:  c,
		BucketCreator: c,
		BucketUpdator: c,
		BucketDeleter: c,
	}
}

type client struct {
}

func (c *client) GetBucket(ctx context.Context, id provider.BucketIdentifier) (provider.Bucket, error) {
	gcp, err := storage.NewClient(ctx)
	if err != nil {
		return provider.Bucket{}, fmt.Errorf("error creating storage client: %v", err)
	}

	defer gcp.Close()

	bucket := gcp.Bucket(string(id.Name))

	attrs, err := bucket.Attrs(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrBucketNotExist) {
			return provider.Bucket{}, sdkerrors.NewErrorNotFound()
		}

		return provider.Bucket{}, fmt.Errorf("oh no: %v", err)
	}

	labels := map[string]any{}
	for k, v := range attrs.Labels {
		labels[k] = v
	}

	return provider.Bucket{
		Identifier: id,
		Config: provider.BucketConfig{
			Labels: labels,
		},
		Attrs: provider.BucketAttrs{
			Created: attrs.Created.String(),
		},
	}, nil
}

func (c *client) CreateBucket(ctx context.Context, id provider.BucketIdentifier, config provider.BucketConfig) (provider.Bucket, error) {
	gcp, err := storage.NewClient(ctx)
	if err != nil {
		return provider.Bucket{}, err
	}

	defer gcp.Close()

	labels := map[string]string{}
	for k, v := range config.Labels {
		str, ok := v.(string)
		if !ok {
			return provider.Bucket{}, fmt.Errorf("expected label value for key %q to be a string, got %T\n", k, v)
		}
		labels[k] = str
	}

	bucket := gcp.Bucket(string(id.Name))
	if err := bucket.Create(ctx, id.Project, &storage.BucketAttrs{
		Labels:   labels,
		Location: id.Location,
	}); err != nil {
		return provider.Bucket{}, err
	}

	attrs, err := bucket.Attrs(ctx)
	if err != nil {
		return provider.Bucket{}, err
	}

	resLabels := map[string]any{}
	for k, v := range attrs.Labels {
		resLabels[k] = v
	}

	return provider.Bucket{
		Identifier: id,
		Config: provider.BucketConfig{
			Labels: resLabels,
		},
		Attrs: provider.BucketAttrs{
			Created: attrs.Created.String(),
		},
	}, nil
}

func (c *client) UpdateBucket(ctx context.Context, id provider.BucketIdentifier, config provider.BucketConfig, mask []value.UpdateMaskField) (provider.Bucket, error) {
	gcp, err := storage.NewClient(ctx)
	if err != nil {
		return provider.Bucket{}, err
	}

	defer gcp.Close()

	toUpdate := storage.BucketAttrsToUpdate{}
	for _, m := range mask {
		switch m.Name {
		case "labels":
			for _, label := range m.SubFields {
				if label.Operation == value.OperationDelete {
					toUpdate.DeleteLabel(label.Name)
				} else {
					val, ok := config.Labels[label.Name]
					if !ok {
						return provider.Bucket{}, fmt.Errorf("value for label %q is missing", label.Name)
					}

					str, ok := val.(string)
					if !ok {
						return provider.Bucket{}, fmt.Errorf("value for label %q is not a string", label.Name)
					}
					toUpdate.SetLabel(label.Name, str)
				}
			}
		}
	}

	bucket := gcp.Bucket(string(id.Name))
	attrs, err := bucket.Update(ctx, toUpdate)
	if err != nil {
		return provider.Bucket{}, err
	}

	labels := map[string]any{}
	for k, v := range attrs.Labels {
		labels[k] = v
	}

	return provider.Bucket{
		Identifier: id,
		Config: provider.BucketConfig{
			Labels: labels,
		},
		Attrs: provider.BucketAttrs{
			Created: attrs.Created.String(),
		},
	}, nil
}

func (c *client) DeleteBucket(ctx context.Context, id provider.BucketIdentifier) error {
	gcp, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	log.Printf("deleting bucket >>>>>>>>>>>>>>>> %v", id.Name)

	bucket := gcp.Bucket(id.Name)
	return bucket.Delete(ctx)
}
