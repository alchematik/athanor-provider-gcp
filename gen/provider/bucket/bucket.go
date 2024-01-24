// Code generated by athanor-go.
// DO NOT EDIT.

package bucket

import (
	"context"
	"fmt"
	sdk "github.com/alchematik/athanor-go/sdk/provider/value"
	"github.com/alchematik/athanor-provider-gcp/gen/provider/identifier"
)

type Bucket struct {
	Identifier identifier.BucketIdentifier
	Config     Config
	Attrs      Attrs
}

func (x Bucket) ToResourceValue() (sdk.Resource, error) {
	id := x.Identifier.ToValue()

	config := x.Config.ToValue()

	attrs := x.Attrs.ToValue()

	return sdk.Resource{
		Identifier: id,
		Config:     config,
		Attrs:      attrs,
	}, nil
}

type BucketGetter interface {
	GetBucket(context.Context, identifier.BucketIdentifier) (Bucket, error)
}

type BucketCreator interface {
	CreateBucket(context.Context, identifier.BucketIdentifier, Config) (Bucket, error)
}

type BucketUpdator interface {
	UpdateBucket(context.Context, identifier.BucketIdentifier, Config, []sdk.UpdateMaskField) (Bucket, error)
}

type BucketDeleter interface {
	DeleteBucket(context.Context, identifier.BucketIdentifier) error
}

type BucketHandler struct {
	BucketGetter  BucketGetter
	BucketCreator BucketCreator
	BucketUpdator BucketUpdator
	BucketDeleter BucketDeleter
}

func (h BucketHandler) GetResource(ctx context.Context, id sdk.Identifier) (sdk.Resource, error) {
	if h.BucketGetter == nil {
		return sdk.Resource{}, fmt.Errorf("unimplemented")
	}

	idVal, err := identifier.ParseBucketIdentifier(id)
	if err != nil {
		return sdk.Resource{}, err
	}

	r, err := h.BucketGetter.GetBucket(ctx, idVal)
	if err != nil {
		return sdk.Resource{}, err
	}

	return r.ToResourceValue()
}

func (h BucketHandler) CreateResource(ctx context.Context, id sdk.Identifier, config any) (sdk.Resource, error) {
	if h.BucketCreator == nil {
		return sdk.Resource{}, fmt.Errorf("unimplemented")
	}

	idVal, err := identifier.ParseBucketIdentifier(id)
	if err != nil {
		return sdk.Resource{}, err
	}

	configVal, err := ParseConfig(config)
	if err != nil {
		return sdk.Resource{}, err
	}

	r, err := h.BucketCreator.CreateBucket(ctx, idVal, configVal)
	if err != nil {
		return sdk.Resource{}, err
	}

	return r.ToResourceValue()
}

func (h BucketHandler) UpdateResource(ctx context.Context, id sdk.Identifier, config any, mask []sdk.UpdateMaskField) (sdk.Resource, error) {
	if h.BucketUpdator == nil {
		return sdk.Resource{}, fmt.Errorf("unimplemented")
	}

	idVal, err := identifier.ParseBucketIdentifier(id)
	if err != nil {
		return sdk.Resource{}, err
	}

	configVal, err := ParseConfig(config)
	if err != nil {
		return sdk.Resource{}, err
	}

	r, err := h.BucketUpdator.UpdateBucket(ctx, idVal, configVal, mask)
	if err != nil {
		return sdk.Resource{}, err
	}

	return r.ToResourceValue()
}

func (h BucketHandler) DeleteResource(ctx context.Context, id sdk.Identifier) error {
	if h.BucketDeleter == nil {
		return fmt.Errorf("unimplemented")
	}

	idVal, err := identifier.ParseBucketIdentifier(id)
	if err != nil {
		return err
	}

	return h.BucketDeleter.DeleteBucket(ctx, idVal)
}

type Attrs struct {
	Create string
}

func (x Attrs) ToValue() any {
	return map[string]any{
		"create": sdk.ToType[any](x.Create),
	}
}

func ParseAttrs(v any) (Attrs, error) {

	m, err := sdk.Map[any](v)
	if err != nil {
		return Attrs{}, nil
	}

	create, err := sdk.String(m["create"])
	if err != nil {
		return Attrs{}, nil
	}

	return Attrs{
		Create: create,
	}, nil
}

type Config struct {
	Labels map[string]string
}

func (x Config) ToValue() any {
	return map[string]any{
		"labels": sdk.ToType[string](x.Labels),
	}
}

func ParseConfig(v any) (Config, error) {

	m, err := sdk.Map[any](v)
	if err != nil {
		return Config{}, nil
	}

	labels, err := sdk.Map[string](m["labels"])
	if err != nil {
		return Config{}, nil
	}

	return Config{
		Labels: labels,
	}, nil
}
