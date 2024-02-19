package iam_role_custom_project

import (
	"context"
	"fmt"
	"log"

	iamrole "github.com/alchematik/athanor-provider-gcp/gen/provider/iam_role_custom_project"
	"github.com/alchematik/athanor-provider-gcp/gen/provider/identifier"

	iamadmin "cloud.google.com/go/iam/admin/apiv1"
	"cloud.google.com/go/iam/admin/apiv1/adminpb"
	"github.com/alchematik/athanor-go/sdk/provider/value"
	gax "github.com/googleapis/gax-go/v2"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
)

func NewHandler(ctx context.Context) (*iamrole.IamRoleCustomProjectHandler, error) {
	gcp, err := iamadmin.NewIamClient(ctx)
	if err != nil {
		return nil, err
	}

	c := &client{GCP: gcp}

	return &iamrole.IamRoleCustomProjectHandler{
		IamRoleCustomProjectGetter:  c,
		IamRoleCustomProjectCreator: c,
		IamRoleCustomProjectDeleter: c,
		IamRoleCustomProjectUpdator: c,
		CloseFunc:                   gcp.Close,
	}, nil
}

type client struct {
	GCP GCP
}

type GCP interface {
	CreateRole(context.Context, *adminpb.CreateRoleRequest, ...gax.CallOption) (*adminpb.Role, error)
	GetRole(context.Context, *adminpb.GetRoleRequest, ...gax.CallOption) (*adminpb.Role, error)
	UpdateRole(context.Context, *adminpb.UpdateRoleRequest, ...gax.CallOption) (*adminpb.Role, error)
	DeleteRole(context.Context, *adminpb.DeleteRoleRequest, ...gax.CallOption) (*adminpb.Role, error)
}

func (c *client) GetIamRoleCustomProject(ctx context.Context, id identifier.IamRoleCustomProjectIdentifier) (iamrole.IamRoleCustomProject, error) {
	res, err := c.GCP.GetRole(ctx, &adminpb.GetRoleRequest{
		Name: fmt.Sprintf("projects/%s/roles/%s", id.Project, id.Name),
	})
	if err != nil {
		return iamrole.IamRoleCustomProject{}, err
	}

	return iamrole.IamRoleCustomProject{
		Identifier: id,
		Config: iamrole.Config{
			Description: res.GetDescription(),
			Title:       res.GetTitle(),
			Permissions: res.GetIncludedPermissions(),
			Stage:       res.GetStage().String(),
		},
		Attrs: iamrole.Attrs{
			Deleted: res.GetDeleted(),
			Etag:    fmt.Sprintf("%x", res.GetEtag()),
		},
	}, nil
}

func (c *client) CreateIamRoleCustomProject(ctx context.Context, id identifier.IamRoleCustomProjectIdentifier, config iamrole.Config) (iamrole.IamRoleCustomProject, error) {
	stage, err := convertStage(config.Stage)
	if err != nil {
		return iamrole.IamRoleCustomProject{}, err
	}

	res, err := c.GCP.CreateRole(ctx, &adminpb.CreateRoleRequest{
		Parent: fmt.Sprintf("projects/%s", id.Project),
		RoleId: id.Name,
		Role: &adminpb.Role{
			Title:               config.Title,
			Description:         config.Description,
			IncludedPermissions: config.Permissions,
			Stage:               stage,
		},
	})
	if err != nil {
		return iamrole.IamRoleCustomProject{}, err
	}

	return iamrole.IamRoleCustomProject{
		Identifier: id,
		Config: iamrole.Config{
			Description: res.GetDescription(),
			Title:       res.GetTitle(),
			Permissions: res.GetIncludedPermissions(),
			Stage:       res.GetStage().String(),
		},
		Attrs: iamrole.Attrs{
			Deleted: res.GetDeleted(),
			Etag:    fmt.Sprintf("%x", res.GetEtag()),
		},
	}, nil
}

func (c *client) UpdateIamRoleCustomProject(ctx context.Context, id identifier.IamRoleCustomProjectIdentifier, config iamrole.Config, mask []value.UpdateMaskField) (iamrole.IamRoleCustomProject, error) {
	log.Printf("UPDATING ROLE>>>>>>>>>>>>>>>")
	updateMask := &fieldmaskpb.FieldMask{}
	var r adminpb.Role

	for _, m := range mask {
		switch m.Name {
		case "title":
			r.Title = config.Title
			updateMask.Paths = append(updateMask.Paths, "title")
		case "description":
			r.Description = config.Description
			updateMask.Paths = append(updateMask.Paths, "description")
		case "stage":
			stage, err := convertStage(config.Stage)
			if err != nil {
				return iamrole.IamRoleCustomProject{}, err
			}
			r.Stage = stage
			updateMask.Paths = append(updateMask.Paths, "stage")
		case "permissions":
			r.IncludedPermissions = config.Permissions
			updateMask.Paths = append(updateMask.Paths, "included_permissions")
		}
	}

	log.Printf("MASK: %+v\n", updateMask.Paths)
	log.Printf("PERMISSIONS: %v\n", r.IncludedPermissions)

	res, err := c.GCP.UpdateRole(ctx, &adminpb.UpdateRoleRequest{
		Name:       fmt.Sprintf("projects/%s/roles/%s", id.Project, id.Name),
		Role:       &r,
		UpdateMask: updateMask,
	})
	if err != nil {
		return iamrole.IamRoleCustomProject{}, err
	}

	log.Printf("RES >>> %+v\n", res.IncludedPermissions)

	return iamrole.IamRoleCustomProject{
		Identifier: id,
		Config: iamrole.Config{
			Description: res.GetDescription(),
			Title:       res.GetTitle(),
			Permissions: res.GetIncludedPermissions(),
			Stage:       res.GetStage().String(),
		},
		Attrs: iamrole.Attrs{
			Deleted: res.GetDeleted(),
			Etag:    fmt.Sprintf("%x", res.GetEtag()),
		},
	}, nil
}

func (c *client) DeleteIamRoleCustomProject(ctx context.Context, id identifier.IamRoleCustomProjectIdentifier) error {
	_, err := c.GCP.DeleteRole(ctx, &adminpb.DeleteRoleRequest{
		Name: fmt.Sprintf("projects/%s/roles/%s", id.Project, id.Name),
	})
	return err
}

func convertStage(str string) (adminpb.Role_RoleLaunchStage, error) {
	switch str {
	case "ALPHA":
		return adminpb.Role_ALPHA, nil
	case "BETA":
		return adminpb.Role_BETA, nil
	case "GA":
		return adminpb.Role_GA, nil
	case "DEPRECATED":
		return adminpb.Role_DEPRECATED, nil
	case "DISABLED":
		return adminpb.Role_DISABLED, nil
	case "EAP":
		return adminpb.Role_EAP, nil
	default:
		return 0, fmt.Errorf("invalid role launch stage: %s", str)
	}
}
