package main

import (
	"github.com/alchematik/athanor-provider-gcp/internal/bucket"
	"github.com/alchematik/athanor-provider-gcp/internal/bucket_object"
	"github.com/alchematik/athanor-provider-gcp/internal/function"
	"github.com/alchematik/athanor-provider-gcp/internal/service_account"

	"github.com/alchematik/athanor-go/sdk/provider/plugin"
)

func main() {
	plugin.Serve(map[string]plugin.ResourceHandler{
		"bucket":          bucket.NewHandler(),
		"bucket_object":   bucket_object.NewHandler(),
		"function":        function.NewHandler(),
		"service_account": service_account.NewHandler(),
	})
}
