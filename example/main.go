package main

import (
	"log"

	"github.com/alchematik/athanor-provider-gcp/gen/sdk/go/api"
	apiconfig "github.com/alchematik/athanor-provider-gcp/gen/sdk/go/api_config"
	"github.com/alchematik/athanor-provider-gcp/gen/sdk/go/bucket"
	bucketobject "github.com/alchematik/athanor-provider-gcp/gen/sdk/go/bucket_object"
	"github.com/alchematik/athanor-provider-gcp/gen/sdk/go/function"
	serviceaccount "github.com/alchematik/athanor-provider-gcp/gen/sdk/go/service_account"

	athanor "github.com/alchematik/athanor-go/sdk/consumer"
)

func main() {
	provider := athanor.Provider{Name: "gcp", Version: "v0.0.1"}

	bucketID := bucket.Identifier{
		Alias:    "my-bucket",
		Project:  "textapp-389501",
		Location: "us-east4",
		Name:     "athanor-test-bucket",
	}
	bucketConfig := bucket.Config{
		Labels: map[string]any{
			"test":    "hello_world",
			"meow":    "is_me",
			"another": "hey",
		},
	}
	myBucket := athanor.Resource{
		Exists:     true,
		Provider:   provider,
		Identifier: bucketID,
		Config:     bucketConfig,
	}

	bp := athanor.Blueprint{}
	bp = bp.WithResource(myBucket)

	bucketObjectID := bucketobject.Identifier{
		Alias:  "my-bucket-object",
		Bucket: bucketID,
		Name:   "my-bucket-object",
	}

	bucketObjectConfig := bucketobject.Config{
		Contents: athanor.File{
			Path: "../test_cloud_func.zip",
		},
	}

	bucketObject := athanor.Resource{
		Exists:     true,
		Provider:   provider,
		Identifier: bucketObjectID,
		Config:     bucketObjectConfig,
	}

	bp = bp.WithResource(bucketObject)

	anotherBucketObject := athanor.Resource{
		Exists:   true,
		Provider: provider,
		Identifier: bucketobject.Identifier{
			Alias:  "my-other-bucket-object",
			Bucket: bucketID,
			Name:   "my-other-bucket-object",
		},
		Config: bucketObjectConfig,
	}

	bp = bp.WithResource(anotherBucketObject)

	funcID := function.Identifier{
		Alias:    "my-function",
		Project:  "textapp-389501",
		Location: "us-east4",
		Name:     "athanor-test-function",
	}
	funcConfig := function.Config{
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
	}
	funcResource := athanor.Resource{
		Exists:     true,
		Provider:   provider,
		Identifier: funcID,
		Config:     funcConfig,
	}

	bp = bp.WithResource(funcResource)

	serviceAccountID := serviceaccount.Identifier{
		Alias:     "my-service-account",
		Project:   "textapp-389501",
		AccountId: "athanor-test",
	}
	serviceAccountConfig := serviceaccount.Config{
		Description: "Test service account",
		DisplayName: "Athanor Test",
	}
	serviceAccountResource := athanor.Resource{
		Exists:     true,
		Provider:   provider,
		Identifier: serviceAccountID,
		Config:     serviceAccountConfig,
	}

	bp = bp.WithResource(serviceAccountResource)

	apiID := api.Identifier{
		Alias:   "my-api",
		ApiId:   "athanor-test",
		Project: "textapp-389501",
	}
	apiConfig := api.Config{
		DisplayName: "test API for Athanor",
		Labels: map[string]any{
			"hello": "world",
		},
	}
	apiResource := athanor.Resource{
		Exists:     true,
		Provider:   provider,
		Identifier: apiID,
		Config:     apiConfig,
	}

	bp = bp.WithResource(apiResource)

	apiConfigID := apiconfig.Identifier{
		Alias:       "my-api-config",
		Api:         apiID,
		ApiConfigId: "athanor-test-config",
	}
	apiConfigConfig := apiconfig.Config{
		DisplayName:    "Athanor test API config!",
		ServiceAccount: serviceAccountID,
		OpenApiDocuments: []any{
			athanor.File{
				Path: "example/openapi.yml",
			},
		},
	}
	apiConfigResource := athanor.Resource{
		Exists:     true,
		Provider:   provider,
		Identifier: apiConfigID,
		Config:     apiConfigConfig,
	}

	bp = bp.WithResource(apiConfigResource)

	if err := athanor.Build(bp); err != nil {
		log.Fatalf("error building blueprint: %v", err)
	}
}
