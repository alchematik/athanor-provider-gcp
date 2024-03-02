// Code generated by athanor-go.
// DO NOT EDIT.

package api_gateway

import (
	"context"
	"fmt"
	sdk "github.com/alchematik/athanor-go/sdk/provider/value"
	"github.com/alchematik/athanor-provider-gcp/gen/provider/identifier"
)

type ApiGateway struct {
	Identifier identifier.ApiGatewayIdentifier
	Config     Config
	Attrs      Attrs
}

func (x ApiGateway) ToResourceValue() (sdk.Resource, error) {
	id := x.Identifier.ToValue()

	config := x.Config.ToValue()

	attrs := x.Attrs.ToValue()

	return sdk.Resource{
		Identifier: id,
		Config:     config,
		Attrs:      attrs,
	}, nil
}

type ApiGatewayGetter interface {
	GetApiGateway(context.Context, identifier.ApiGatewayIdentifier) (ApiGateway, error)
}

type ApiGatewayCreator interface {
	CreateApiGateway(context.Context, identifier.ApiGatewayIdentifier, Config) (ApiGateway, error)
}

type ApiGatewayUpdator interface {
	UpdateApiGateway(context.Context, identifier.ApiGatewayIdentifier, Config, []sdk.UpdateMaskField) (ApiGateway, error)
}

type ApiGatewayDeleter interface {
	DeleteApiGateway(context.Context, identifier.ApiGatewayIdentifier) error
}

type ApiGatewayHandler struct {
	ApiGatewayGetter  ApiGatewayGetter
	ApiGatewayCreator ApiGatewayCreator
	ApiGatewayUpdator ApiGatewayUpdator
	ApiGatewayDeleter ApiGatewayDeleter

	CloseFunc func() error
}

func (h *ApiGatewayHandler) GetResource(ctx context.Context, id sdk.Identifier) (sdk.Resource, error) {
	if h.ApiGatewayGetter == nil {
		return sdk.Resource{}, fmt.Errorf("unimplemented")
	}

	idVal, err := identifier.ParseApiGatewayIdentifier(id)
	if err != nil {
		return sdk.Resource{}, err
	}

	r, err := h.ApiGatewayGetter.GetApiGateway(ctx, idVal)
	if err != nil {
		return sdk.Resource{}, err
	}

	return r.ToResourceValue()
}

func (h *ApiGatewayHandler) CreateResource(ctx context.Context, id sdk.Identifier, config any) (sdk.Resource, error) {
	if h.ApiGatewayCreator == nil {
		return sdk.Resource{}, fmt.Errorf("unimplemented")
	}

	idVal, err := identifier.ParseApiGatewayIdentifier(id)
	if err != nil {
		return sdk.Resource{}, err
	}

	configVal, err := ParseConfig(config)
	if err != nil {
		return sdk.Resource{}, err
	}

	r, err := h.ApiGatewayCreator.CreateApiGateway(ctx, idVal, configVal)
	if err != nil {
		return sdk.Resource{}, err
	}

	return r.ToResourceValue()
}

func (h *ApiGatewayHandler) UpdateResource(ctx context.Context, id sdk.Identifier, config any, mask []sdk.UpdateMaskField) (sdk.Resource, error) {
	if h.ApiGatewayUpdator == nil {
		return sdk.Resource{}, fmt.Errorf("unimplemented")
	}

	idVal, err := identifier.ParseApiGatewayIdentifier(id)
	if err != nil {
		return sdk.Resource{}, err
	}

	configVal, err := ParseConfig(config)
	if err != nil {
		return sdk.Resource{}, err
	}

	r, err := h.ApiGatewayUpdator.UpdateApiGateway(ctx, idVal, configVal, mask)
	if err != nil {
		return sdk.Resource{}, err
	}

	return r.ToResourceValue()
}

func (h *ApiGatewayHandler) DeleteResource(ctx context.Context, id sdk.Identifier) error {
	if h.ApiGatewayDeleter == nil {
		return fmt.Errorf("unimplemented")
	}

	idVal, err := identifier.ParseApiGatewayIdentifier(id)
	if err != nil {
		return err
	}

	return h.ApiGatewayDeleter.DeleteApiGateway(ctx, idVal)
}

func (h *ApiGatewayHandler) Close() error {
	if h.CloseFunc != nil {
		return h.CloseFunc()
	}

	return nil
}

type Attrs struct {
	Create          string
	DefaultHostname string
	State           string
	Update          string
}

func (x Attrs) ToValue() any {
	return map[string]any{
		"create":           sdk.ToType[any](x.Create),
		"default_hostname": sdk.ToType[any](x.DefaultHostname),
		"state":            sdk.ToType[any](x.State),
		"update":           sdk.ToType[any](x.Update),
	}
}

func ParseAttrs(v any) (Attrs, error) {
	m, err := sdk.Map[any](v)
	if err != nil {
		return Attrs{}, fmt.Errorf("error parsing attrs: %v", err)
	}

	create, err := sdk.String(m["create"])
	if err != nil {
		return Attrs{}, fmt.Errorf("error parsing attrs for api_gateway: %v", err)
	}
	default_hostname, err := sdk.String(m["default_hostname"])
	if err != nil {
		return Attrs{}, fmt.Errorf("error parsing attrs for api_gateway: %v", err)
	}
	state, err := sdk.String(m["state"])
	if err != nil {
		return Attrs{}, fmt.Errorf("error parsing attrs for api_gateway: %v", err)
	}
	update, err := sdk.String(m["update"])
	if err != nil {
		return Attrs{}, fmt.Errorf("error parsing attrs for api_gateway: %v", err)
	}

	return Attrs{
		Create:          create,
		DefaultHostname: default_hostname,
		State:           state,
		Update:          update,
	}, nil
}

func ParseAttrsList(v any) ([]Attrs, error) {
	list, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("invalid type for list: %T", v)
	}

	var vals []Attrs
	for _, val := range list {
		p, err := ParseAttrs(val)
		if err != nil {
			return nil, err
		}

		vals = append(vals, p)
	}

	return vals, nil
}

type Config struct {
	ApiConfig   sdk.ResourceIdentifier
	DisplayName string
	Labels      map[string]string
}

func (x Config) ToValue() any {
	return map[string]any{
		"api_config":   sdk.ToType[any](x.ApiConfig),
		"display_name": sdk.ToType[any](x.DisplayName),
		"labels":       sdk.ToType[string](x.Labels),
	}
}

func ParseConfig(v any) (Config, error) {
	m, err := sdk.Map[any](v)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing config: %v", err)
	}

	api_config, err := identifier.ParseIdentifier(m["api_config"])
	if err != nil {
		return Config{}, fmt.Errorf("error parsing config for api_gateway: %v", err)
	}
	display_name, err := sdk.String(m["display_name"])
	if err != nil {
		return Config{}, fmt.Errorf("error parsing config for api_gateway: %v", err)
	}
	labels, err := sdk.Map[string](m["labels"])
	if err != nil {
		return Config{}, fmt.Errorf("error parsing config for api_gateway: %v", err)
	}

	return Config{
		ApiConfig:   api_config,
		DisplayName: display_name,
		Labels:      labels,
	}, nil
}

func ParseConfigList(v any) ([]Config, error) {
	list, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("invalid type for list: %T", v)
	}

	var vals []Config
	for _, val := range list {
		p, err := ParseConfig(val)
		if err != nil {
			return nil, err
		}

		vals = append(vals, p)
	}

	return vals, nil
}