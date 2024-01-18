package main

import (
	"github.com/alchematik/athanor-go/sdk/provider/schema"
)

var function = schema.ResourceSchema{
	Type: "function",
	Identifier: schema.FieldSchema{
		IsIdentifier: true,
		Type:         schema.FieldTypeStruct,
		Fields: []schema.FieldSchema{
			{
				Name: "project",
				Type: schema.FieldTypeString,
			},
			{
				Name: "location",
				Type: schema.FieldTypeString,
			},
			{
				Name: "name",
				Type: schema.FieldTypeString,
			},
		},
	},
	Config: schema.FieldSchema{
		Type: schema.FieldTypeStruct,
		Fields: []schema.FieldSchema{
			{
				Name: "description",
				Type: schema.FieldTypeString,
			},
			{
				Name: "labels",
				Type: schema.FieldTypeMap,
			},
			{
				Name: "build_config",
				Type: schema.FieldTypeStruct,
				Fields: []schema.FieldSchema{
					{
						Name: "runtime",
						Type: schema.FieldTypeString,
					},
					{
						Name: "entrypoint",
						Type: schema.FieldTypeString,
					},
					{
						Name: "source",
						Type: schema.FieldTypeFile,
					},
				},
			},
		},
	},
	Attrs: schema.FieldSchema{
		Type: schema.FieldTypeStruct,
		Fields: []schema.FieldSchema{
			{
				Name: "url",
				Type: schema.FieldTypeString,
			},
		},
	},
}
