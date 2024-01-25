// Code generated by athanor-go.
// DO NOT EDIT.

package api

import (
	sdk "github.com/alchematik/athanor-go/sdk/consumer"
)

type Config struct {
	DisplayName any
	Labels      any
}

func (x Config) ToExpr() any {
	return map[string]any{
		"display_name": x.DisplayName,
		"labels":       x.Labels,
	}
}

type Identifier struct {
	Alias   string
	ApiId   any
	Project any
}

func (x Identifier) ToExpr() any {
	return sdk.ResourceIdentifier{
		ResourceType: "api",
		Alias:        x.Alias,
		Value: map[string]any{
			"api_id":  x.ApiId,
			"project": x.Project,
		},
	}
}