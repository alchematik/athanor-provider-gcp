package main

import (
	"github.com/alchematik/athanor-provider-gcp/gen/sdk/go/api"
	apiconfig "github.com/alchematik/athanor-provider-gcp/gen/sdk/go/api_config"
	apigateway "github.com/alchematik/athanor-provider-gcp/gen/sdk/go/api_gateway"
	"github.com/alchematik/athanor-provider-gcp/gen/sdk/go/bucket"
	bucketobject "github.com/alchematik/athanor-provider-gcp/gen/sdk/go/bucket_object"
	"github.com/alchematik/athanor-provider-gcp/gen/sdk/go/function"
	iampolicy "github.com/alchematik/athanor-provider-gcp/gen/sdk/go/iam_policy"
	iamcustomrole "github.com/alchematik/athanor-provider-gcp/gen/sdk/go/iam_role_custom_project"
	serviceaccount "github.com/alchematik/athanor-provider-gcp/gen/sdk/go/service_account"

	athanor "github.com/alchematik/athanor-go/sdk/consumer"
)

func main() {
	athanor.Build(func(_ any) (athanor.Blueprint, error) {
		bp := athanor.Blueprint{}

		provider := athanor.Provider{
			Name:    "gcp",
			Version: "v0.0.1",
			Repo: athanor.RepoLocal{
				Path: "build/provider",
			},
		}

		myBucket := athanor.Resource{
			Exists:   true,
			Provider: provider,
			Identifier: bucket.Identifier{
				Alias:    "my-bucket",
				Project:  "textapp-389501",
				Location: "us-east4",
				Name:     "athanor-test-bucket",
			},
			Config: bucket.Config{
				Labels: map[string]any{
					"test":    "hello_world",
					"meow":    "is_me",
					"another": "hey_hey_hey",
					"foo":     "bar",
				},
			},
		}

		bp = bp.WithResource(myBucket)

		bp = bp.WithBuild(
			"sub-blueprint",
			athanor.RepoLocal{
				Path: "./example/sub",
			},
			athanor.Translator{
				Name:    "go",
				Version: "v0.0.1",
				Repo: athanor.RepoLocal{
					Path: "build/translator",
				},
			},
			map[string]any{
				"bucket_exists": true,
				"bucket_name":   "athanor-test-sub-bucket",
			},
			nil,
		)

		bucketObjectConfig := bucketobject.Config{
			Contents: athanor.File{
				Path: "../test_cloud_func.zip",
			},
		}

		bucketObject := athanor.Resource{
			Exists:   true,
			Provider: provider,
			Identifier: bucketobject.Identifier{
				Alias:  "my-bucket-object",
				Bucket: myBucket.Identifier,
				Name:   "my-bucket-object",
			},
			Config: bucketObjectConfig,
		}

		bp = bp.WithResource(bucketObject)

		anotherBucketObject := athanor.Resource{
			Exists:   true,
			Provider: provider,
			Identifier: bucketobject.Identifier{
				Alias:  "my-other-bucket-object",
				Bucket: myBucket.Identifier,
				Name:   "my-other-bucket-object",
			},
			Config: bucketObjectConfig,
		}

		bp = bp.WithResource(anotherBucketObject)

		funcResource := athanor.Resource{
			Exists:   true,
			Provider: provider,
			Identifier: function.Identifier{
				Alias:    "my-function",
				Project:  "textapp-389501",
				Location: "us-east4",
				Name:     "athanor-test-function",
			},
			Config: function.Config{
				Description: "test function managed by athanor",
				Labels: map[string]any{
					"test":          "true",
					"another_label": "hi",
				},
				BuildConfig: function.BuildConfig{
					Runtime:    "go121",
					Entrypoint: "HelloHTTP",
					Source: athanor.File{
						Path: "../test_cloud_func.zip",
					},
				},
			},
		}

		bp = bp.WithResource(funcResource)

		serviceAccountResource := athanor.Resource{
			Exists:   true,
			Provider: provider,
			Identifier: serviceaccount.Identifier{
				Alias:     "my-service-account",
				Project:   "textapp-389501",
				AccountId: "athanor-test",
			},
			Config: serviceaccount.Config{
				Description: "Test service account",
				DisplayName: "Athanor Test",
			},
		}

		bp = bp.WithResource(serviceAccountResource)

		apiResource := athanor.Resource{
			Exists:   true,
			Provider: provider,
			Identifier: api.Identifier{
				Alias:   "my-api",
				ApiId:   "athanor-test",
				Project: "textapp-389501",
			},
			Config: api.Config{
				DisplayName: "test API for Athanor",
				Labels: map[string]any{
					"hello": "world",
				},
			},
		}

		bp = bp.WithResource(apiResource)

		apiConfigResource := athanor.Resource{
			Exists:   true,
			Provider: provider,
			Identifier: apiconfig.Identifier{
				Alias:       "my-api-config",
				Api:         apiResource.Identifier,
				ApiConfigId: "athanor-test-config",
			},
			Config: apiconfig.Config{
				DisplayName:    "Athanor test API config!",
				ServiceAccount: serviceAccountResource.Identifier,
				OpenApiDocuments: []any{
					athanor.File{
						Path: "example/openapi.yml",
					},
				},
			},
		}

		bp = bp.WithResource(apiConfigResource)

		apiGatewayResource := athanor.Resource{
			Exists:   true,
			Provider: provider,
			Identifier: apigateway.Identifier{
				Alias:     "my-api-gateway",
				Project:   "textapp-389501",
				Location:  "us-east4",
				GatewayId: "athanor-test-gateway",
			},
			Config: apigateway.Config{
				ApiConfig:   apiConfigResource.Identifier,
				DisplayName: "Athanor test gateway!",
				Labels: map[string]any{
					"test": "yes",
				},
			},
		}
		bp = bp.WithResource(apiGatewayResource)

		testRole := athanor.Resource{
			Exists:   true,
			Provider: provider,
			Identifier: iamcustomrole.Identifier{
				Alias:   "test-role",
				Project: "textapp-389501",
				Name:    "testrole",
			},
			Config: iamcustomrole.Config{
				Title:       "Test role",
				Description: "Test role for invoking cloud functions.",
				Stage:       "ALPHA",
				Permissions: []any{
					"cloudfunctions.functions.invoke",
					"run.jobs.run",
					"run.routes.invoke",
				},
			},
		}
		bp = bp.WithResource(testRole)

		functionPolicy := athanor.Resource{
			Exists:   true,
			Provider: provider,
			Identifier: iampolicy.Identifier{
				Alias:    "my-function-policy",
				Resource: funcResource.Identifier,
			},
			Config: iampolicy.Config{
				Bindings: []any{
					iampolicy.Binding{
						Role: testRole.Identifier,
						Members: []any{
							serviceAccountResource.Identifier,
						},
					},
				},
			},
		}
		bp = bp.WithResource(functionPolicy)

		return bp, nil
	})
}
