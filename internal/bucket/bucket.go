package bucket

import (
	"context"
	"errors"
	"fmt"

	"github.com/alchematik/athanor-provider-gcp/gen/provider"

	"cloud.google.com/go/storage"
	sdkerrors "github.com/alchematik/athanor-go/sdk/errors"
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

func (c *client) UpdateBucket(ctx context.Context, id provider.BucketIdentifier, config provider.BucketConfig) (provider.Bucket, error) {
	return provider.Bucket{}, nil
}

func (c *client) DeleteBucket(ctx context.Context, id provider.BucketIdentifier) error {
	return nil
}
