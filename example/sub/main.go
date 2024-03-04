package main

import (
	"github.com/alchematik/athanor-provider-gcp/gen/sdk/go/bucket"

	athanor "github.com/alchematik/athanor-go/sdk/consumer"
)

func main() {
	athanor.Build(func(args ...any) (athanor.Blueprint, error) {

		exists := args[0]
		name := args[1]

		bp := athanor.Blueprint{}

		provider := athanor.Provider{
			Name:    "gcp",
			Version: "v0.0.1",
			Repo: athanor.RepoLocal{
				Path: "build/provider",
			},
		}

		b := athanor.Resource{
			Exists:   exists,
			Provider: provider,
			Identifier: bucket.Identifier{
				Alias:    "sub-resource-bucket",
				Project:  "textapp-389501",
				Location: "us-east4",
				Name:     name,
			},
			Config: bucket.Config{
				Labels: map[string]any{
					"foo": "bar",
					"hi":  athanor.RuntimeConfig{},
				},
			},
		}

		bp = bp.WithResource(b)
		return bp, nil
	})
}
