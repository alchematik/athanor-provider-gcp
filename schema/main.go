package main

import (
	"log"

	"github.com/alchematik/athanor-go/sdk/provider/schema"
)

func main() {
	s := schema.Schema{
		Name:    "gcp",
		Version: "v0.0.1",
		Resources: []schema.ResourceSchema{
			bucket,
			bucketObject,
			function,
		},
	}

	if err := schema.Build(s); err != nil {
		log.Fatalf("error building provider schema: %v", err)
	}
}
