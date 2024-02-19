package main

import (
	"github.com/alchematik/athanor-go/sdk/provider/schema"
)

var iamRole = schema.ResourceSchema{
	Type: "iam_role",
	Identifier: schema.Struct("identifier", map[string]schema.FieldSchema{
		"name": schema.String(),
	}),
	Config: schema.Struct("config", map[string]schema.FieldSchema{}),
	Attrs: schema.Struct("attrs", map[string]schema.FieldSchema{
		"title":       schema.String(),
		"description": schema.String(),
		"stage":       schema.String(),
		"etag":        schema.String(),
		"permissions": schema.List(schema.String()),
	}),
}
