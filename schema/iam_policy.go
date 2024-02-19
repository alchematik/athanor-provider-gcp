package main

import (
	"github.com/alchematik/athanor-go/sdk/provider/schema"
)

var iamPolicy = schema.ResourceSchema{
	Type: "iam_policy",
	Identifier: schema.Struct("identifier", map[string]schema.FieldSchema{
		"resource": schema.Identifier(),
	}),
	Config: schema.Struct("config", map[string]schema.FieldSchema{
		"bindings": schema.List(schema.Struct("binding", map[string]schema.FieldSchema{
			"role":    schema.Identifier(),
			"members": schema.List(schema.Identifier()),
		})),
	}),
	Attrs: schema.Struct("attrs", map[string]schema.FieldSchema{
		"etag": schema.String(),
	}),
}
