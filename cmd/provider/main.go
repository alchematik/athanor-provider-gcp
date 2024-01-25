package main

import (
	"context"
	"github.com/alchematik/athanor-provider-gcp/internal/bucket"
	"github.com/alchematik/athanor-provider-gcp/internal/bucket_object"
	"github.com/alchematik/athanor-provider-gcp/internal/function"
	"github.com/alchematik/athanor-provider-gcp/internal/service_account"

	"github.com/alchematik/athanor-go/sdk/provider/plugin"
)

func main() {
	plugin.Serve(map[string]plugin.ResoureceHandlerInitializer{
		"bucket": func(ctx context.Context) (plugin.ResourceHandler, error) {
			return bucket.NewHandler(ctx)
		},
		"bucket_object": func(ctx context.Context) (plugin.ResourceHandler, error) {
			return bucket_object.NewHandler(ctx)
		},
		"function": func(ctx context.Context) (plugin.ResourceHandler, error) {
			return function.NewHandler(ctx)
		},
		"service_account": func(ctx context.Context) (plugin.ResourceHandler, error) {
			return service_account.NewHandler(ctx)
		},
	})
}
