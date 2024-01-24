// Code generated by athanor-go.
// DO NOT EDIT.

package service_account

import (
	"context"
	"fmt"
	sdk "github.com/alchematik/athanor-go/sdk/provider/value"
	"github.com/alchematik/athanor-provider-gcp/gen/provider/identifier"
)

type ServiceAccount struct {
	Identifier identifier.ServiceAccountIdentifier
	Config     ServiceAccountConfig
	Attrs      ServiceAccountAttrs
}

func (x ServiceAccount) ToResourceValue() (sdk.Resource, error) {
	id := x.Identifier.ToValue()

	config := x.Config.ToValue()

	attrs := x.Attrs.ToValue()

	return sdk.Resource{
		Identifier: id,
		Config:     config,
		Attrs:      attrs,
	}, nil
}

type ServiceAccountGetter interface {
	GetServiceAccount(context.Context, identifier.ServiceAccountIdentifier) (ServiceAccount, error)
}

type ServiceAccountCreator interface {
	CreateServiceAccount(context.Context, identifier.ServiceAccountIdentifier, ServiceAccountConfig) (ServiceAccount, error)
}

type ServiceAccountUpdator interface {
	UpdateServiceAccount(context.Context, identifier.ServiceAccountIdentifier, ServiceAccountConfig, []sdk.UpdateMaskField) (ServiceAccount, error)
}

type ServiceAccountDeleter interface {
	DeleteServiceAccount(context.Context, identifier.ServiceAccountIdentifier) error
}

type ServiceAccountHandler struct {
	ServiceAccountGetter  ServiceAccountGetter
	ServiceAccountCreator ServiceAccountCreator
	ServiceAccountUpdator ServiceAccountUpdator
	ServiceAccountDeleter ServiceAccountDeleter
}

func (h ServiceAccountHandler) GetResource(ctx context.Context, id sdk.Identifier) (sdk.Resource, error) {
	if h.ServiceAccountGetter == nil {
		return sdk.Resource{}, fmt.Errorf("unimplemented")
	}

	idVal, err := identifier.ParseServiceAccountIdentifier(id)
	if err != nil {
		return sdk.Resource{}, err
	}

	r, err := h.ServiceAccountGetter.GetServiceAccount(ctx, idVal)
	if err != nil {
		return sdk.Resource{}, err
	}

	return r.ToResourceValue()
}

func (h ServiceAccountHandler) CreateResource(ctx context.Context, id sdk.Identifier, config any) (sdk.Resource, error) {
	if h.ServiceAccountCreator == nil {
		return sdk.Resource{}, fmt.Errorf("unimplemented")
	}

	idVal, err := identifier.ParseServiceAccountIdentifier(id)
	if err != nil {
		return sdk.Resource{}, err
	}

	configVal, err := ParseServiceAccountConfig(config)
	if err != nil {
		return sdk.Resource{}, err
	}

	r, err := h.ServiceAccountCreator.CreateServiceAccount(ctx, idVal, configVal)
	if err != nil {
		return sdk.Resource{}, err
	}

	return r.ToResourceValue()
}

func (h ServiceAccountHandler) UpdateResource(ctx context.Context, id sdk.Identifier, config any, mask []sdk.UpdateMaskField) (sdk.Resource, error) {
	if h.ServiceAccountUpdator == nil {
		return sdk.Resource{}, fmt.Errorf("unimplemented")
	}

	idVal, err := identifier.ParseServiceAccountIdentifier(id)
	if err != nil {
		return sdk.Resource{}, err
	}

	configVal, err := ParseServiceAccountConfig(config)
	if err != nil {
		return sdk.Resource{}, err
	}

	r, err := h.ServiceAccountUpdator.UpdateServiceAccount(ctx, idVal, configVal, mask)
	if err != nil {
		return sdk.Resource{}, err
	}

	return r.ToResourceValue()
}

func (h ServiceAccountHandler) DeleteResource(ctx context.Context, id sdk.Identifier) error {
	if h.ServiceAccountDeleter == nil {
		return fmt.Errorf("unimplemented")
	}

	idVal, err := identifier.ParseServiceAccountIdentifier(id)
	if err != nil {
		return err
	}

	return h.ServiceAccountDeleter.DeleteServiceAccount(ctx, idVal)
}

type ServiceAccountAttrs struct {
	UniqueId string
	Disabled bool
}

func (x ServiceAccountAttrs) ToValue() any {
	return map[string]any{
		"unique_id": sdk.ToType(x.UniqueId),
		"disabled":  sdk.ToType(x.Disabled),
	}
}

func ParseServiceAccountAttrs(v any) (ServiceAccountAttrs, error) {

	m, err := sdk.Map(v)
	if err != nil {
		return ServiceAccountAttrs{}, nil
	}

	unique_id, err := sdk.String(m["unique_id"])
	if err != nil {
		return ServiceAccountAttrs{}, nil
	}
	disabled, err := sdk.Bool(m["disabled"])
	if err != nil {
		return ServiceAccountAttrs{}, nil
	}

	return ServiceAccountAttrs{
		UniqueId: unique_id,
		Disabled: disabled,
	}, nil
}

type ServiceAccountConfig struct {
	DisplayName string
	Description string
}

func (x ServiceAccountConfig) ToValue() any {
	return map[string]any{
		"display_name": sdk.ToType(x.DisplayName),
		"description":  sdk.ToType(x.Description),
	}
}

func ParseServiceAccountConfig(v any) (ServiceAccountConfig, error) {

	m, err := sdk.Map(v)
	if err != nil {
		return ServiceAccountConfig{}, nil
	}

	display_name, err := sdk.String(m["display_name"])
	if err != nil {
		return ServiceAccountConfig{}, nil
	}
	description, err := sdk.String(m["description"])
	if err != nil {
		return ServiceAccountConfig{}, nil
	}

	return ServiceAccountConfig{
		DisplayName: display_name,
		Description: description,
	}, nil
}