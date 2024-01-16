package bucket_object

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/alchematik/athanor-provider-gcp/gen/provider/bucket"
	bucketobject "github.com/alchematik/athanor-provider-gcp/gen/provider/bucket_object"

	"cloud.google.com/go/storage"
	sdkerrors "github.com/alchematik/athanor-go/sdk/errors"
	"github.com/alchematik/athanor-go/sdk/provider/value"
)

func NewHandler() bucketobject.BucketObjectHandler {
	c := &client{}
	return bucketobject.BucketObjectHandler{
		BucketObjectGetter:  c,
		BucketObjectCreator: c,
		BucketObjectUpdator: c,
		BucketObjectDeleter: c,
	}
}

type client struct{}

func (c *client) GetBucketObject(ctx context.Context, id bucketobject.BucketObjectIdentifier) (bucketobject.BucketObject, error) {
	gcp, err := storage.NewClient(ctx)
	if err != nil {
		return bucketobject.BucketObject{}, fmt.Errorf("error creating storage client: %v", err)
	}

	defer gcp.Close()

	bucketID, err := bucket.ParseBucketIdentifier(id.Bucket)
	if err != nil {
		return bucketobject.BucketObject{}, err
	}

	b := gcp.Bucket(bucketID.Name)
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
		Config: bucketobject.BucketObjectConfig{
			Contents: value.File{
				Checksum: fmt.Sprintf("%x", attrs.MD5),
			},
		},
		Attrs: bucketobject.BucketObjectAttrs{
			Created: attrs.Created.String(),
		},
	}, nil
}

func (c *client) CreateBucketObject(ctx context.Context, id bucketobject.BucketObjectIdentifier, config bucketobject.BucketObjectConfig) (bucketobject.BucketObject, error) {
	gcp, err := storage.NewClient(ctx)
	if err != nil {
		return bucketobject.BucketObject{}, fmt.Errorf("error creating storage client: %v", err)
	}

	defer gcp.Close()

	bucketID, err := bucket.ParseBucketIdentifier(id.Bucket)
	if err != nil {
		return bucketobject.BucketObject{}, err
	}

	b := gcp.Bucket(bucketID.Name)
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
		Config: bucketobject.BucketObjectConfig{
			Contents: value.File{
				Checksum: fmt.Sprintf("%x", attrs.MD5),
			},
		},
		Attrs: bucketobject.BucketObjectAttrs{
			Created: attrs.Created.String(),
		},
	}, nil
}

func (c *client) UpdateBucketObject(ctx context.Context, id bucketobject.BucketObjectIdentifier, config bucketobject.BucketObjectConfig, mask []value.UpdateMaskField) (bucketobject.BucketObject, error) {
	gcp, err := storage.NewClient(ctx)
	if err != nil {
		return bucketobject.BucketObject{}, fmt.Errorf("error creating storage client: %v", err)
	}

	defer gcp.Close()

	bucketID, err := bucket.ParseBucketIdentifier(id.Bucket)
	if err != nil {
		return bucketobject.BucketObject{}, err
	}

	b := gcp.Bucket(bucketID.Name)
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
		Config: bucketobject.BucketObjectConfig{
			Contents: value.File{
				Checksum: fmt.Sprintf("%x", attrs.MD5),
			},
		},
		Attrs: bucketobject.BucketObjectAttrs{
			Created: attrs.Created.String(),
		},
	}, nil
}

func (c *client) DeleteBucketObject(ctx context.Context, id bucketobject.BucketObjectIdentifier) error {
	gcp, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating storage client: %v", err)
	}

	defer gcp.Close()

	bucketID, err := bucket.ParseBucketIdentifier(id.Bucket)
	if err != nil {
		return err
	}

	b := gcp.Bucket(bucketID.Name)
	object := b.Object(id.Name)
	return object.Delete(ctx)
}
