package main

import (
	"log"

	"github.com/alchematik/athanor-provider-gcp/gen/sdk/go/bucket"
	bucketobject "github.com/alchematik/athanor-provider-gcp/gen/sdk/go/bucket_object"

	athanor "github.com/alchematik/athanor-go/sdk/consumer"
)

func main() {
	provider := athanor.Provider{Name: "gcp", Version: "v0.0.1"}

	bucketID := bucket.BucketIdentifier{
		Alias:    "my-bucket",
		Project:  "textapp-389501",
		Location: "us-east4",
		Name:     "athanor-test-bucket",
	}
	bucketConfig := bucket.BucketConfig{
		Labels: map[string]any{
			"test": "hello_world",
			"meow": "is_me",
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

	bucketObjectID := bucketobject.BucketObjectIdentifier{
		Alias:  "my-bucket-object",
		Bucket: bucketID,
		Name:   "my-bucket-object",
	}

	bucketObjectConfig := bucketobject.BucketObjectConfig{
		Contents: athanor.File{Path: "example/config.json"},
	}

	bucketObject := athanor.Resource{
		Exists:     true,
		Provider:   provider,
		Identifier: bucketObjectID,
		Config:     bucketObjectConfig,
	}

	bp = bp.WithResource(bucketObject)

	if err := athanor.Build(bp); err != nil {
		log.Fatalf("error building blueprint: %v", err)
	}
}
