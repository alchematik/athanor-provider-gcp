// Code generated by athanor-go.
// DO NOT EDIT.

package iam_policy

import (
	sdk "github.com/alchematik/athanor-go/sdk/consumer"
)

type Binding struct {
	Members any
	Role    any
}

func (x Binding) ToExpr() any {
	return map[string]any{
		"members": x.Members,
		"role":    x.Role,
	}
}

type Config struct {
	Bindings any
}

func (x Config) ToExpr() any {
	return map[string]any{
		"bindings": x.Bindings,
	}
}

type Identifier struct {
	Alias    string
	Resource any
}

func (x Identifier) ToExpr() any {
	return sdk.ResourceIdentifier{
		ResourceType: "iam_policy",
		Alias:        x.Alias,
		Value: map[string]any{
			"resource": x.Resource,
		},
	}
}
