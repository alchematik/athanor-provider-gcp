package bucket_object

import (
	"context"
	"errors"
	"fmt"
	"os"

	bucketobject "github.com/alchematik/athanor-provider-gcp/gen/provider/bucket_object"
	"github.com/alchematik/athanor-provider-gcp/gen/provider/identifier"

	"cloud.google.com/go/storage"
	sdkerrors "github.com/alchematik/athanor-go/sdk/errors"
	"github.com/alchematik/athanor-go/sdk/provider/value"
)

func NewHandler(ctx context.Context) (*bucketobject.BucketObjectHandler, error) {
	gcp, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating GCP storage client: %v", err)
	}

	c := &client{
		Storage: gcp,
	}
	return &bucketobject.BucketObjectHandler{
		BucketObjectGetter:  c,
		BucketObjectCreator: c,
		BucketObjectUpdator: c,
		BucketObjectDeleter: c,
		CloseFunc:           gcp.Close,
	}, nil
}

type client struct {
	Storage Storage
}

type Storage interface {
	Bucket(string) *storage.BucketHandle
}

func (c *client) GetBucketObject(ctx context.Context, id identifier.BucketObjectIdentifier) (bucketobject.BucketObject, error) {
	bucketID, ok := id.Bucket.(identifier.BucketIdentifier)
	if !ok {
		return bucketobject.BucketObject{}, fmt.Errorf("field bucket must be a bucket identifier")
	}

	b := c.Storage.Bucket(bucketID.Name)
	obj := b.Object(id.Name)
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return bucketobject.BucketObject{}, sdkerrors.NewErrorNotFound()
		}

		return bucketobject.BucketObject{}, err
	}

	return bucketobject.BucketObject{
		Identifier: id,
		Config: bucketobject.Config{
			Contents: value.File{
				Checksum: fmt.Sprintf("%d", attrs.CRC32C),
			},
		},
		Attrs: bucketobject.Attrs{
			Create: attrs.Created.String(),
		},
	}, nil
}

func (c *client) CreateBucketObject(ctx context.Context, id identifier.BucketObjectIdentifier, config bucketobject.Config) (bucketobject.BucketObject, error) {
	bucketID, ok := id.Bucket.(identifier.BucketIdentifier)
	if !ok {
		return bucketobject.BucketObject{}, fmt.Errorf("field bucket must be a bucket identifier")
	}

	b := c.Storage.Bucket(bucketID.Name)
	object := b.Object(id.Name)
	w := object.NewWriter(ctx)
	data, err := os.ReadFile(config.Contents.Path)
	if err != nil {
		return bucketobject.BucketObject{}, err
	}

	if _, err := w.Write(data); err != nil {
		return bucketobject.BucketObject{}, err
	}

	if err := w.Close(); err != nil {
		return bucketobject.BucketObject{}, err
	}

	attrs, err := object.Attrs(ctx)
	if err != nil {
		return bucketobject.BucketObject{}, err
	}

	return bucketobject.BucketObject{
		Identifier: id,
		Config: bucketobject.Config{
			Contents: value.File{
				Checksum: fmt.Sprintf("%d", attrs.CRC32C),
			},
		},
		Attrs: bucketobject.Attrs{
			Create: attrs.Created.String(),
		},
	}, nil
}

func (c *client) UpdateBucketObject(ctx context.Context, id identifier.BucketObjectIdentifier, config bucketobject.Config, mask []value.UpdateMaskField) (bucketobject.BucketObject, error) {
	bucketID, ok := id.Bucket.(identifier.BucketIdentifier)
	if !ok {
		return bucketobject.BucketObject{}, fmt.Errorf("field bucket must be a bucket identifier")
	}

	b := c.Storage.Bucket(bucketID.Name)
	object := b.Object(id.Name)
	w := object.NewWriter(ctx)
	data, err := os.ReadFile(config.Contents.Path)
	if err != nil {
		return bucketobject.BucketObject{}, err
	}

	if _, err := w.Write(data); err != nil {
		return bucketobject.BucketObject{}, err
	}

	if err := w.Close(); err != nil {
		return bucketobject.BucketObject{}, err
	}

	attrs, err := object.Attrs(ctx)
	if err != nil {
		return bucketobject.BucketObject{}, err
	}

	return bucketobject.BucketObject{
		Identifier: id,
		Config: bucketobject.Config{
			Contents: value.File{
				Checksum: fmt.Sprintf("%d", attrs.CRC32C),
			},
		},
		Attrs: bucketobject.Attrs{
			Create: attrs.Created.String(),
		},
	}, nil
}

func (c *client) DeleteBucketObject(ctx context.Context, id identifier.BucketObjectIdentifier) error {
	bucketID, ok := id.Bucket.(identifier.BucketIdentifier)
	if !ok {
		return fmt.Errorf("field bucket must be a bucket identifier")
	}

	b := c.Storage.Bucket(bucketID.Name)
	object := b.Object(id.Name)
	return object.Delete(ctx)
}
