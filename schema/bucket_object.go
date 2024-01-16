package main

import (
	"github.com/alchematik/athanor-go/sdk/provider/schema"
)

var bucketObject = schema.ResourceSchema{
	Type: "bucket_object",
	Identifier: schema.FieldSchema{
		IsIdentifier: true,
		Type:         schema.FieldTypeStruct,
		Fields: []schema.FieldSchema{
			{
				Name: "bucket",
				Type: schema.FieldTypeIdentifier,
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
				Name: "contents",
				Type: schema.FieldTypeFile,
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
