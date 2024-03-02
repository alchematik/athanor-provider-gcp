// Code generated by athanor-go.
// DO NOT EDIT.

package identifier

import (
	"fmt"

	sdk "github.com/alchematik/athanor-go/sdk/provider/value"
)

type IamRoleCustomProjectIdentifier struct {
	Name    string
	Project string
}

func (x IamRoleCustomProjectIdentifier) ToValue() sdk.Identifier {
	return sdk.Identifier{
		ResourceType: "iam_role_custom_project",
		Value: map[string]any{
			"name":    sdk.ToType[any](x.Name),
			"project": sdk.ToType[any](x.Project),
		},
	}
}

func (x IamRoleCustomProjectIdentifier) ResourceType() string {
	return "iam_role_custom_project"
}

func ParseIamRoleCustomProjectIdentifier(v sdk.Identifier) (IamRoleCustomProjectIdentifier, error) {

	m, err := sdk.Map[any](v.Value)
	if err != nil {
		return IamRoleCustomProjectIdentifier{}, fmt.Errorf("error parsing iam_role_custom_project_identifier: %v", err)
	}

	name, err := sdk.String(m["name"])
	if err != nil {
		return IamRoleCustomProjectIdentifier{}, fmt.Errorf("error parsing iam_role_custom_project_identifier: %v", err)
	}
	project, err := sdk.String(m["project"])
	if err != nil {
		return IamRoleCustomProjectIdentifier{}, fmt.Errorf("error parsing iam_role_custom_project_identifier: %v", err)
	}

	return IamRoleCustomProjectIdentifier{
		Name:    name,
		Project: project,
	}, nil
}