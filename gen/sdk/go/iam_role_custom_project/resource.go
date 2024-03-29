// Code generated by athanor-go.
// DO NOT EDIT.

package iam_role_custom_project

import (
	sdk "github.com/alchematik/athanor-go/sdk/consumer"
)

type Config struct {
	Description any
	Permissions any
	Stage       any
	Title       any
}

func (x Config) ToExpr() any {
	return map[string]any{
		"description": x.Description,
		"permissions": x.Permissions,
		"stage":       x.Stage,
		"title":       x.Title,
	}
}

type Identifier struct {
	Alias   string
	Name    any
	Project any
}

func (x Identifier) ToExpr() any {
	return sdk.ResourceIdentifier{
		ResourceType: "iam_role_custom_project",
		Alias:        x.Alias,
		Value: map[string]any{
			"name":    x.Name,
			"project": x.Project,
		},
	}
}
