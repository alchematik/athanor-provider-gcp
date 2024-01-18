// Code generated by athanor-go.
// DO NOT EDIT.

package bucket

import (
	sdk "github.com/alchematik/athanor-go/sdk/consumer"
)

type BucketConfig struct {
	Labels any
}

func (x BucketConfig) ToExpr() any {
	return map[string]any{
		"labels": x.Labels,
	}
}

type BucketIdentifier struct {
	Alias string

	Project  any
	Location any
	Name     any
}

func (x BucketIdentifier) ToExpr() any {
	return sdk.ResourceIdentifier{
		ResourceType: "bucket",
		Alias:        x.Alias,
		Value: map[string]any{
			"project":  x.Project,
			"location": x.Location,
			"name":     x.Name,
		},
	}
}