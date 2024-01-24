package main

import (
	"github.com/alchematik/athanor-go/sdk/provider/schema"
)

var api = schema.ResourceSchema{
	Type: "api",
	Identifier: schema.Struct("identifier", map[string]schema.FieldSchema{
		"project": schema.String(),
		"api_id":  schema.String(),
	}),
	Config: schema.Struct("config", map[string]schema.FieldSchema{
		"display_name": schema.String(),
		"labels":       schema.Map(schema.String()),
	}),
	Attrs: schema.Struct("attrs", map[string]schema.FieldSchema{
		"create": schema.String(),
		"update": schema.String(),
		"state":  schema.String(),
	}),
}
