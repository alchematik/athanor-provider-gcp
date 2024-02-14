package main

import (
	"github.com/alchematik/athanor-go/sdk/provider/schema"
)

var apiGateway = schema.ResourceSchema{
	Type: "api_gateway",
	Identifier: schema.Struct(
		"identifier",
		map[string]schema.FieldSchema{
			"project":    schema.String(),
			"location":   schema.String(),
			"gateway_id": schema.String(),
		},
	),
	Config: schema.Struct(
		"config",
		map[string]schema.FieldSchema{
			"labels":       schema.Map(schema.String()),
			"display_name": schema.String(),
			"api_config":   schema.Identifier(),
		},
	),
	Attrs: schema.Struct("attrs", map[string]schema.FieldSchema{
		"create":           schema.String(),
		"update":           schema.String(),
		"state":            schema.String(),
		"default_hostname": schema.String(),
	}),
}
