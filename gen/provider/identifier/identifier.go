// Code generated by athanor-go.
// DO NOT EDIT.

package identifier

import (
	"fmt"

	sdk "github.com/alchematik/athanor-go/sdk/provider/value"
)

func ParseIdentifier(v any) (sdk.ResourceIdentifier, error) {
	id, ok := v.(sdk.Identifier)
	if !ok {
		return nil, fmt.Errorf("expected Identifier type, got %T", v)
	}

	switch id.ResourceType {
	case "bucket":
		return ParseBucketIdentifier(id)
	case "bucket_object":
		return ParseBucketObjectIdentifier(id)
	case "function":
		return ParseFunctionIdentifier(id)
	case "service_account":
		return ParseServiceAccountIdentifier(id)

	default:
		return nil, fmt.Errorf("invalid resource type: %s", id.ResourceType)
	}
}
