package main

import (
	"github.com/alchematik/athanor-go/sdk/provider/schema"
)

var function = schema.ResourceSchema{
	Type: "function",
	Identifier: schema.Struct("identifier", map[string]schema.FieldSchema{
		"project":  schema.String(),
		"location": schema.String(),
		"name":     schema.String(),
	}),
	Config: schema.Struct("config", map[string]schema.FieldSchema{
		"description": schema.String(),
		"labels":      schema.Map(schema.String()),
		"build_config": schema.Struct("build_config", map[string]schema.FieldSchema{
			"runtime":    schema.String(),
			"entrypoint": schema.String(),
			"source":     schema.File(),
		}),
	}),
	Attrs: schema.Struct("attrs", map[string]schema.FieldSchema{
		"url": schema.String(),
	}),
}
