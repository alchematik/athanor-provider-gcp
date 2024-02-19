package iam_role

import (
	"context"
	"fmt"

	iamrole "github.com/alchematik/athanor-provider-gcp/gen/provider/iam_role"
	"github.com/alchematik/athanor-provider-gcp/gen/provider/identifier"

	iamadmin "cloud.google.com/go/iam/admin/apiv1"
	"cloud.google.com/go/iam/admin/apiv1/adminpb"
	gax "github.com/googleapis/gax-go/v2"
)

func NewHandler(ctx context.Context) (*iamrole.IamRoleHandler, error) {
	gcp, err := iamadmin.NewIamClient(ctx)
	if err != nil {
		return nil, err
	}
	c := &client{
		GCP: gcp,
	}

	return &iamrole.IamRoleHandler{
		IamRoleGetter: c,
	}, nil
}

type client struct {
	GCP GCP
}

type GCP interface {
	GetRole(context.Context, *adminpb.GetRoleRequest, ...gax.CallOption) (*adminpb.Role, error)
}

func (c *client) GetIamRole(ctx context.Context, id identifier.IamRoleIdentifier) (iamrole.IamRole, error) {
	res, err := c.GCP.GetRole(ctx, &adminpb.GetRoleRequest{
		Name: fmt.Sprintf("roles/%s", id.Name),
	})
	if err != nil {
		return iamrole.IamRole{}, err
	}

	return iamrole.IamRole{
		Identifier: id,
		Config:     iamrole.Config{},
		Attrs: iamrole.Attrs{
			Title:       res.GetTitle(),
			Description: res.GetDescription(),
			Stage:       res.GetStage().String(),
			Etag:        fmt.Sprintf("%x", res.GetEtag()),
			Permissions: res.GetIncludedPermissions(),
		},
	}, nil
}
