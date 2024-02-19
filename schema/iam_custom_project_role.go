package main

import (
	"github.com/alchematik/athanor-go/sdk/provider/schema"
)

var iamRoleCustomProject = schema.ResourceSchema{
	Type: "iam_role_custom_project",
	Identifier: schema.Struct("identifier", map[string]schema.FieldSchema{
		"project": schema.String(),
		"name":    schema.String(),
	}),
	Config: schema.Struct("config", map[string]schema.FieldSchema{
		"title":       schema.String(),
		"description": schema.String(),
		"permissions": schema.List(schema.String()),
		"stage":       schema.String(),
	}),
	Attrs: schema.Struct("attrs", map[string]schema.FieldSchema{
		"deleted": schema.Bool(),
		"etag":    schema.String(),
	}),
}
