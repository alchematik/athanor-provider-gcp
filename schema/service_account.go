package main

import (
	"github.com/alchematik/athanor-go/sdk/provider/schema"
)

var serviceAccount = schema.ResourceSchema{
	Type: "service_account",
	Identifier: schema.FieldSchema{
		IsIdentifier: true,
		Type:         schema.FieldTypeStruct,
		Fields: []schema.FieldSchema{
			{
				Name: "project",
				Type: schema.FieldTypeString,
			},
			{
				Name: "account_id",
				Type: schema.FieldTypeString,
			},
		},
	},
	Config: schema.FieldSchema{
		Type: schema.FieldTypeStruct,
		Fields: []schema.FieldSchema{
			{
				Name: "display_name",
				Type: schema.FieldTypeString,
			},
			{
				Name: "description",
				Type: schema.FieldTypeString,
			},
		},
	},
	Attrs: schema.FieldSchema{
		Type: schema.FieldTypeStruct,
		Fields: []schema.FieldSchema{
			{
				Name: "unique_id",
				Type: schema.FieldTypeString,
			},
			{
				Name: "disabled",
				Type: schema.FieldTypeBool,
			},
		},
	},
}
