package main

import (
	"github.com/alchematik/athanor-go/sdk/provider/schema"
)

var serviceAccount = schema.ResourceSchema{
	Type: "service_account",
	Identifier: schema.Struct("identifier", map[string]schema.FieldSchema{
		"project":    schema.String(),
		"account_id": schema.String(),
	}),
	Config: schema.Struct("config", map[string]schema.FieldSchema{
		"display_name": schema.String(),
		"description":  schema.String(),
	}),
	Attrs: schema.Struct("attrs", map[string]schema.FieldSchema{
		"unique_id": schema.String(),
		"disabled":  schema.Bool(),
	}),
}
