// Code generated by athanor-go.
// DO NOT EDIT.

package identifier

import (
	"fmt"

	sdk "github.com/alchematik/athanor-go/sdk/provider/value"
)

type ApiIdentifier struct {
	ApiId   string
	Project string
}

func (x ApiIdentifier) ToValue() sdk.Identifier {
	return sdk.Identifier{
		ResourceType: "api",
		Value: map[string]any{
			"api_id":  sdk.ToType[any](x.ApiId),
			"project": sdk.ToType[any](x.Project),
		},
	}
}

func (x ApiIdentifier) ResourceType() string {
	return "api"
}

func ParseApiIdentifier(v sdk.Identifier) (ApiIdentifier, error) {

	m, err := sdk.Map[any](v.Value)
	if err != nil {
		return ApiIdentifier{}, fmt.Errorf("error parsing api_identifier: %v", err)
	}

	api_id, err := sdk.String(m["api_id"])
	if err != nil {
		return ApiIdentifier{}, fmt.Errorf("error parsing api_identifier: %v", err)
	}
	project, err := sdk.String(m["project"])
	if err != nil {
		return ApiIdentifier{}, fmt.Errorf("error parsing api_identifier: %v", err)
	}

	return ApiIdentifier{
		ApiId:   api_id,
		Project: project,
	}, nil
}
