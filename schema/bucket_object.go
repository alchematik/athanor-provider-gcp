package main

import (
	"github.com/alchematik/athanor-go/sdk/provider/schema"
)

var bucketObject = schema.ResourceSchema{
	Type: "bucket_object",
	Identifier: schema.Struct("identifier", map[string]schema.FieldSchema{
		"bucket": schema.Identifier(),
		"name":   schema.String(),
	}),
	Config: schema.Struct("config", map[string]schema.FieldSchema{
		"contents": schema.File(),
	}),
	Attrs: schema.Struct("attrs", map[string]schema.FieldSchema{
		"create": schema.String(),
	}),
}
