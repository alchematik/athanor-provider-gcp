// Code generated by athanor-go.
// DO NOT EDIT.

package identifier

import (
	"fmt"

	sdk "github.com/alchematik/athanor-go/sdk/provider/value"
)

type IamRoleIdentifier struct {
	Name string
}

func (x IamRoleIdentifier) ToValue() sdk.Identifier {
	return sdk.Identifier{
		ResourceType: "iam_role",
		Value: map[string]any{
			"name": sdk.ToType[any](x.Name),
		},
	}
}

func (x IamRoleIdentifier) ResourceType() string {
	return "iam_role"
}

func ParseIamRoleIdentifier(v sdk.Identifier) (IamRoleIdentifier, error) {

	m, err := sdk.Map[any](v.Value)
	if err != nil {
		return IamRoleIdentifier{}, fmt.Errorf("error parsing iam_role_identifier: %v", err)
	}

	name, err := sdk.String(m["name"])
	if err != nil {
		return IamRoleIdentifier{}, fmt.Errorf("error parsing iam_role_identifier: %v", err)
	}

	return IamRoleIdentifier{
		Name: name,
	}, nil
}
