package main

import (
	"context"
	"github.com/alchematik/athanor-provider-gcp/internal/api"
	"github.com/alchematik/athanor-provider-gcp/internal/api_config"
	"github.com/alchematik/athanor-provider-gcp/internal/api_gateway"
	"github.com/alchematik/athanor-provider-gcp/internal/bucket"
	"github.com/alchematik/athanor-provider-gcp/internal/bucket_object"
	"github.com/alchematik/athanor-provider-gcp/internal/function"
	"github.com/alchematik/athanor-provider-gcp/internal/iam_policy"
	"github.com/alchematik/athanor-provider-gcp/internal/iam_role"
	"github.com/alchematik/athanor-provider-gcp/internal/iam_role_custom_project"
	"github.com/alchematik/athanor-provider-gcp/internal/service_account"

	"github.com/alchematik/athanor-go/sdk/provider/plugin"
)

func main() {
	plugin.Serve(map[string]plugin.ResoureceHandlerInitializer{
		"api": func(ctx context.Context) (plugin.ResourceHandler, error) {
			return api.NewHandler(ctx)
		},
		"api_config": func(ctx context.Context) (plugin.ResourceHandler, error) {
			return api_config.NewHandler(ctx)
		},
		"api_gateway": func(ctx context.Context) (plugin.ResourceHandler, error) {
			return api_gateway.NewHandler(ctx)
		},
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
		"iam_role": func(ctx context.Context) (plugin.ResourceHandler, error) {
			return iam_role.NewHandler(ctx)
		},
		"iam_role_custom_project": func(ctx context.Context) (plugin.ResourceHandler, error) {
			return iam_role_custom_project.NewHandler(ctx)
		},
		"iam_policy": func(ctx context.Context) (plugin.ResourceHandler, error) {
			return iam_policy.NewHandler(ctx)
		},
	})
}
