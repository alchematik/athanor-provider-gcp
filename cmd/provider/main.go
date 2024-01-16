package main

import (
	"github.com/alchematik/athanor-provider-gcp/internal/bucket"
	"github.com/alchematik/athanor-provider-gcp/internal/bucket_object"

	"github.com/alchematik/athanor-go/sdk/provider/plugin"
)

func main() {
	plugin.Serve(map[string]plugin.ResourceHandler{
		"bucket":        bucket.NewHandler(),
		"bucket_object": bucket_object.NewHandler(),
	})
}
