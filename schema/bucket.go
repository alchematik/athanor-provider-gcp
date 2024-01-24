package main

import (
	"github.com/alchematik/athanor-go/sdk/provider/schema"
)

var bucket = schema.ResourceSchema{
	Type: "bucket",
	Identifier: schema.Struct("identifier", map[string]schema.FieldSchema{
		"project":  schema.String(),
		"location": schema.String(),
		"name":     schema.String(),
	}),
	Config: schema.Struct("config", map[string]schema.FieldSchema{
		"labels": schema.Map(schema.String()),
	}),
	Attrs: schema.Struct("attrs", map[string]schema.FieldSchema{
		"create": schema.String(),
	}),
}
