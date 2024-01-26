// Code generated by athanor-go.
// DO NOT EDIT.

package api_config

import (
	"context"
	"fmt"
	sdk "github.com/alchematik/athanor-go/sdk/provider/value"
	"github.com/alchematik/athanor-provider-gcp/gen/provider/identifier"
)

type ApiConfig struct {
	Identifier identifier.ApiConfigIdentifier
	Config     Config
	Attrs      Attrs
}

func (x ApiConfig) ToResourceValue() (sdk.Resource, error) {
	id := x.Identifier.ToValue()

	config := x.Config.ToValue()

	attrs := x.Attrs.ToValue()

	return sdk.Resource{
		Identifier: id,
		Config:     config,
		Attrs:      attrs,
	}, nil
}

type ApiConfigGetter interface {
	GetApiConfig(context.Context, identifier.ApiConfigIdentifier) (ApiConfig, error)
}

type ApiConfigCreator interface {
	CreateApiConfig(context.Context, identifier.ApiConfigIdentifier, Config) (ApiConfig, error)
}

type ApiConfigUpdator interface {
	UpdateApiConfig(context.Context, identifier.ApiConfigIdentifier, Config, []sdk.UpdateMaskField) (ApiConfig, error)
}

type ApiConfigDeleter interface {
	DeleteApiConfig(context.Context, identifier.ApiConfigIdentifier) error
}

type ApiConfigHandler struct {
	ApiConfigGetter  ApiConfigGetter
	ApiConfigCreator ApiConfigCreator
	ApiConfigUpdator ApiConfigUpdator
	ApiConfigDeleter ApiConfigDeleter

	CloseFunc func() error
}

func (h *ApiConfigHandler) GetResource(ctx context.Context, id sdk.Identifier) (sdk.Resource, error) {
	if h.ApiConfigGetter == nil {
		return sdk.Resource{}, fmt.Errorf("unimplemented")
	}

	idVal, err := identifier.ParseApiConfigIdentifier(id)
	if err != nil {
		return sdk.Resource{}, err
	}

	r, err := h.ApiConfigGetter.GetApiConfig(ctx, idVal)
	if err != nil {
		return sdk.Resource{}, err
	}

	return r.ToResourceValue()
}

func (h *ApiConfigHandler) CreateResource(ctx context.Context, id sdk.Identifier, config any) (sdk.Resource, error) {
	if h.ApiConfigCreator == nil {
		return sdk.Resource{}, fmt.Errorf("unimplemented")
	}

	idVal, err := identifier.ParseApiConfigIdentifier(id)
	if err != nil {
		return sdk.Resource{}, err
	}

	configVal, err := ParseConfig(config)
	if err != nil {
		return sdk.Resource{}, err
	}

	r, err := h.ApiConfigCreator.CreateApiConfig(ctx, idVal, configVal)
	if err != nil {
		return sdk.Resource{}, err
	}

	return r.ToResourceValue()
}

func (h *ApiConfigHandler) UpdateResource(ctx context.Context, id sdk.Identifier, config any, mask []sdk.UpdateMaskField) (sdk.Resource, error) {
	if h.ApiConfigUpdator == nil {
		return sdk.Resource{}, fmt.Errorf("unimplemented")
	}

	idVal, err := identifier.ParseApiConfigIdentifier(id)
	if err != nil {
		return sdk.Resource{}, err
	}

	configVal, err := ParseConfig(config)
	if err != nil {
		return sdk.Resource{}, err
	}

	r, err := h.ApiConfigUpdator.UpdateApiConfig(ctx, idVal, configVal, mask)
	if err != nil {
		return sdk.Resource{}, err
	}

	return r.ToResourceValue()
}

func (h *ApiConfigHandler) DeleteResource(ctx context.Context, id sdk.Identifier) error {
	if h.ApiConfigDeleter == nil {
		return fmt.Errorf("unimplemented")
	}

	idVal, err := identifier.ParseApiConfigIdentifier(id)
	if err != nil {
		return err
	}

	return h.ApiConfigDeleter.DeleteApiConfig(ctx, idVal)
}

func (h *ApiConfigHandler) Close() error {
	if h.CloseFunc != nil {
		return h.CloseFunc()
	}

	return nil
}

type Attrs struct {
	Create string
	State  string
	Update string
}

func (x Attrs) ToValue() any {
	return map[string]any{
		"create": sdk.ToType[any](x.Create),
		"state":  sdk.ToType[any](x.State),
		"update": sdk.ToType[any](x.Update),
	}
}

func ParseAttrs(v any) (Attrs, error) {
	m, err := sdk.Map[any](v)
	if err != nil {
		return Attrs{}, fmt.Errorf("error parsing attrs: %v", err)
	}

	create, err := sdk.String(m["create"])
	if err != nil {
		return Attrs{}, fmt.Errorf("error parsing attrs for api_config: %v", err)
	}
	state, err := sdk.String(m["state"])
	if err != nil {
		return Attrs{}, fmt.Errorf("error parsing attrs for api_config: %v", err)
	}
	update, err := sdk.String(m["update"])
	if err != nil {
		return Attrs{}, fmt.Errorf("error parsing attrs for api_config: %v", err)
	}

	return Attrs{
		Create: create,
		State:  state,
		Update: update,
	}, nil
}

type Config struct {
	DisplayName      string
	OpenApiDocuments []sdk.File
}

func (x Config) ToValue() any {
	return map[string]any{
		"display_name":       sdk.ToType[any](x.DisplayName),
		"open_api_documents": sdk.ToImmutableType(sdk.ToType[sdk.File])(x.OpenApiDocuments),
	}
}

func ParseConfig(v any) (Config, error) {
	m, err := sdk.Map[any](v)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing config: %v", err)
	}

	display_name, err := sdk.String(m["display_name"])
	if err != nil {
		return Config{}, fmt.Errorf("error parsing config for api_config: %v", err)
	}
	open_api_documents, err := sdk.List[sdk.File](m["open_api_documents"])
	if err != nil {
		return Config{}, fmt.Errorf("error parsing config for api_config: %v", err)
	}

	return Config{
		DisplayName:      display_name,
		OpenApiDocuments: open_api_documents,
	}, nil
}
