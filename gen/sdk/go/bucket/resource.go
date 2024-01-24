// Code generated by athanor-go.
// DO NOT EDIT.

package bucket

import (
	sdk "github.com/alchematik/athanor-go/sdk/consumer"
)

type Config struct {
	Labels any
}

func (x Config) ToExpr() any {
	return map[string]any{
		"labels": x.Labels,
	}
}

type Identifier struct {
	Alias    string
	Location any
	Name     any
	Project  any
}

func (x Identifier) ToExpr() any {
	return sdk.ResourceIdentifier{
		ResourceType: "bucket",
		Alias:        x.Alias,
		Value: map[string]any{
			"location": x.Location,
			"name":     x.Name,
			"project":  x.Project,
		},
	}
}
