package main

import (
	"github.com/alchematik/athanor-go/sdk/provider/schema"
)

var apiConfig = schema.ResourceSchema{
	Type: "api_config",
	Identifier: schema.Struct("identifier", map[string]schema.FieldSchema{
		"api":           schema.Identifier(),
		"api_config_id": schema.String(),
	}),
	Config: schema.Struct("config", map[string]schema.FieldSchema{
		"display_name":       schema.String(),
		"open_api_documents": schema.List(schema.File()),
	}),
	Attrs: schema.Struct("attrs", map[string]schema.FieldSchema{
		"create": schema.String(),
		"update": schema.String(),
		"state":  schema.String(),
	}),
}
