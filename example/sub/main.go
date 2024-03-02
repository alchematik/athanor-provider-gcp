package main

import (
	"fmt"

	"github.com/alchematik/athanor-provider-gcp/gen/sdk/go/bucket"

	athanor "github.com/alchematik/athanor-go/sdk/consumer"
)

func main() {
	athanor.Build(func(args any) (athanor.Blueprint, error) {
		m, ok := args.(map[string]any)
		if !ok {
			return athanor.Blueprint{}, fmt.Errorf("expected map, got %T", args)
		}

		bp := athanor.Blueprint{}

		provider := athanor.Provider{
			Name:    "gcp",
			Version: "v0.0.1",
			Repo: athanor.RepoLocal{
				Path: "build/provider",
			},
		}

		b := athanor.Resource{
			Exists:   m["bucket_exists"],
			Provider: provider,
			Identifier: bucket.Identifier{
				Alias:    "sub-resource-bucket",
				Project:  "textapp-389501",
				Location: "us-east4",
				Name:     m["bucket_name"],
			},
			Config: bucket.Config{
				Labels: map[string]any{
					"foo": "bar",
				},
			},
		}

		bp = bp.WithResource(b)
		return bp, nil
	})
}
