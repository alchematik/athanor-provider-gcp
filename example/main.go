package main

import (
	"log"

	gcp "github.com/alchematik/athanor-provider-gcp/gen/sdk/go"

	athanor "github.com/alchematik/athanor-go/sdk/consumer"
)

func main() {
	provider := athanor.Provider("gcp", athanor.String("gcp"), athanor.String("v0.0.1"))

	bucketID := gcp.BucketIdentifier{
		Alias:    "my-bucket",
		Project:  athanor.String("textapp-389501"),
		Location: athanor.String("us-east4"),
		Name:     athanor.String("athanor-test-bucket"),
	}
	bucketConfig := gcp.BucketConfig{
		Labels: athanor.Map(map[string]athanor.Type{
			"test": athanor.String("hello"),
		}),
	}
	bucket := athanor.Resource(
		athanor.Bool(true),
		provider,
		bucketID,
		bucketConfig,
	)

	bp := athanor.Blueprint{}
	bp = bp.WithResource(bucket)

	if err := athanor.Build(bp); err != nil {
		log.Fatalf("error building blueprint: %v", err)
	}
}
