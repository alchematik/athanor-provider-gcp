package main

import (
	"github.com/alchematik/athanor-go/sdk/provider/schema"
)

var bucket = schema.ResourceSchema{
	Type: "bucket",
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
				Name: "labels",
				Type: schema.FieldTypeMap,
			},
		},
	},
	Attrs: schema.FieldSchema{
		Type: schema.FieldTypeStruct,
		Fields: []schema.FieldSchema{
			{
				Name: "created",
				Type: schema.FieldTypeString,
			},
		},
	},
}
