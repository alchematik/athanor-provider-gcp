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

func NewHandler(ctx context.Context) (*bucket.BucketHandler, error) {
	gcp, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating GCP storage client: %v", err)
	}

	c := &client{
		Storage: gcp,
	}
	return &bucket.BucketHandler{
		BucketGetter:  c,
		BucketCreator: c,
		BucketUpdator: c,
		BucketDeleter: c,
		CloseFunc:     gcp.Close,
	}, nil
}

type client struct {
	Storage Storage
}

type Storage interface {
	Bucket(string) *storage.BucketHandle
}

func (c *client) GetBucket(ctx context.Context, id identifier.BucketIdentifier) (bucket.Bucket, error) {
	b := c.Storage.Bucket(id.Name)

	attrs, err := b.Attrs(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrBucketNotExist) {
			return bucket.Bucket{}, sdkerrors.NewErrorNotFound()
		}

		return bucket.Bucket{}, fmt.Errorf("oh no: %v", err)
	}

	return bucket.Bucket{
		Identifier: id,
		Config: bucket.Config{
			Labels: attrs.Labels,
		},
		Attrs: bucket.Attrs{
			Create: attrs.Created.String(),
			Etag:   fmt.Sprintf("%x", attrs.Etag),
		},
	}, nil
}

func (c *client) CreateBucket(ctx context.Context, id identifier.BucketIdentifier, config bucket.Config) (bucket.Bucket, error) {
	b := c.Storage.Bucket(id.Name)
	if err := b.Create(ctx, id.Project, &storage.BucketAttrs{
		Labels:   config.Labels,
		Location: id.Location,
	}); err != nil {
		return bucket.Bucket{}, err
	}

	attrs, err := b.Attrs(ctx)
	if err != nil {
		return bucket.Bucket{}, err
	}

	return bucket.Bucket{
		Identifier: id,
		Config: bucket.Config{
			Labels: attrs.Labels,
		},
		Attrs: bucket.Attrs{
			Create: attrs.Created.String(),
			Etag:   fmt.Sprintf("%x", attrs.Etag),
		},
	}, nil
}

func (c *client) UpdateBucket(ctx context.Context, id identifier.BucketIdentifier, config bucket.Config, mask []value.UpdateMaskField) (bucket.Bucket, error) {
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

					toUpdate.SetLabel(label.Name, val)
				}
			}
		}
	}

	b := c.Storage.Bucket(string(id.Name))
	attrs, err := b.Update(ctx, toUpdate)
	if err != nil {
		return bucket.Bucket{}, err
	}

	return bucket.Bucket{
		Identifier: id,
		Config: bucket.Config{
			Labels: attrs.Labels,
		},
		Attrs: bucket.Attrs{
			Create: attrs.Created.String(),
			Etag:   fmt.Sprintf("%x", attrs.Etag),
		},
	}, nil
}

func (c *client) DeleteBucket(ctx context.Context, id identifier.BucketIdentifier) error {
	b := c.Storage.Bucket(id.Name)
	return b.Delete(ctx)
}
