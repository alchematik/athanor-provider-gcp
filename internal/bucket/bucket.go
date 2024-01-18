package bucket

import (
	"context"
	"errors"
	"fmt"

	"github.com/alchematik/athanor-provider-gcp/gen/provider/bucket"
	"github.com/alchematik/athanor-provider-gcp/gen/provider/identifier"

	"cloud.google.com/go/storage"
	sdkerrors "github.com/alchematik/athanor-go/sdk/errors"
	value "github.com/alchematik/athanor-go/sdk/provider/value"
)

func NewHandler() bucket.BucketHandler {
	c := &client{}
	return bucket.BucketHandler{
		BucketGetter:  c,
		BucketCreator: c,
		BucketUpdator: c,
		BucketDeleter: c,
	}
}

type client struct {
}

func (c *client) GetBucket(ctx context.Context, id identifier.BucketIdentifier) (bucket.Bucket, error) {
	gcp, err := storage.NewClient(ctx)
	if err != nil {
		return bucket.Bucket{}, fmt.Errorf("error creating storage client: %v", err)
	}

	defer gcp.Close()

	b := gcp.Bucket(string(id.Name))

	attrs, err := b.Attrs(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrBucketNotExist) {
			return bucket.Bucket{}, sdkerrors.NewErrorNotFound()
		}

		return bucket.Bucket{}, fmt.Errorf("oh no: %v", err)
	}

	labels := map[string]any{}
	for k, v := range attrs.Labels {
		labels[k] = v
	}

	return bucket.Bucket{
		Identifier: id,
		Config: bucket.BucketConfig{
			Labels: labels,
		},
		Attrs: bucket.BucketAttrs{
			Created: attrs.Created.String(),
		},
	}, nil
}

func (c *client) CreateBucket(ctx context.Context, id identifier.BucketIdentifier, config bucket.BucketConfig) (bucket.Bucket, error) {
	gcp, err := storage.NewClient(ctx)
	if err != nil {
		return bucket.Bucket{}, err
	}

	defer gcp.Close()

	labels := map[string]string{}
	for k, v := range config.Labels {
		str, ok := v.(string)
		if !ok {
			return bucket.Bucket{}, fmt.Errorf("expected label value for key %q to be a string, got %T\n", k, v)
		}
		labels[k] = str
	}

	b := gcp.Bucket(string(id.Name))
	if err := b.Create(ctx, id.Project, &storage.BucketAttrs{
		Labels:   labels,
		Location: id.Location,
	}); err != nil {
		return bucket.Bucket{}, err
	}

	attrs, err := b.Attrs(ctx)
	if err != nil {
		return bucket.Bucket{}, err
	}

	resLabels := map[string]any{}
	for k, v := range attrs.Labels {
		resLabels[k] = v
	}

	return bucket.Bucket{
		Identifier: id,
		Config: bucket.BucketConfig{
			Labels: resLabels,
		},
		Attrs: bucket.BucketAttrs{
			Created: attrs.Created.String(),
		},
	}, nil
}

func (c *client) UpdateBucket(ctx context.Context, id identifier.BucketIdentifier, config bucket.BucketConfig, mask []value.UpdateMaskField) (bucket.Bucket, error) {
	gcp, err := storage.NewClient(ctx)
	if err != nil {
		return bucket.Bucket{}, err
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
						return bucket.Bucket{}, fmt.Errorf("value for label %q is missing", label.Name)
					}

					str, ok := val.(string)
					if !ok {
						return bucket.Bucket{}, fmt.Errorf("value for label %q is not a string", label.Name)
					}
					toUpdate.SetLabel(label.Name, str)
				}
			}
		}
	}

	b := gcp.Bucket(string(id.Name))
	attrs, err := b.Update(ctx, toUpdate)
	if err != nil {
		return bucket.Bucket{}, err
	}

	labels := map[string]any{}
	for k, v := range attrs.Labels {
		labels[k] = v
	}

	return bucket.Bucket{
		Identifier: id,
		Config: bucket.BucketConfig{
			Labels: labels,
		},
		Attrs: bucket.BucketAttrs{
			Created: attrs.Created.String(),
		},
	}, nil
}

func (c *client) DeleteBucket(ctx context.Context, id identifier.BucketIdentifier) error {
	gcp, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	b := gcp.Bucket(id.Name)
	return b.Delete(ctx)
}
