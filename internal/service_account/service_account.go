package service_account

import (
	"context"
	"fmt"

	"github.com/alchematik/athanor-provider-gcp/gen/provider/identifier"
	serviceaccount "github.com/alchematik/athanor-provider-gcp/gen/provider/service_account"

	iamadmin "cloud.google.com/go/iam/admin/apiv1"
	"cloud.google.com/go/iam/admin/apiv1/adminpb"
	sdkerrors "github.com/alchematik/athanor-go/sdk/errors"
	"github.com/alchematik/athanor-go/sdk/provider/value"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewHandler(ctx context.Context) (*serviceaccount.ServiceAccountHandler, error) {
	gcp, err := iamadmin.NewIamClient(ctx)
	if err != nil {
		return nil, err
	}
	c := &client{
		GCP: gcp,
	}
	return &serviceaccount.ServiceAccountHandler{
		ServiceAccountGetter:  c,
		ServiceAccountCreator: c,
		ServiceAccountUpdator: c,
		ServiceAccountDeleter: c,
		CloseFunc:             gcp.Close,
	}, nil
}

type client struct {
	GCP GCP
}

type GCP interface {
	GetServiceAccount(context.Context, *adminpb.GetServiceAccountRequest, ...gax.CallOption) (*adminpb.ServiceAccount, error)
	CreateServiceAccount(context.Context, *adminpb.CreateServiceAccountRequest, ...gax.CallOption) (*adminpb.ServiceAccount, error)
	UpdateServiceAccount(context.Context, *adminpb.ServiceAccount, ...gax.CallOption) (*adminpb.ServiceAccount, error)
	DeleteServiceAccount(context.Context, *adminpb.DeleteServiceAccountRequest, ...gax.CallOption) error
}

func (c *client) GetServiceAccount(ctx context.Context, id identifier.ServiceAccountIdentifier) (serviceaccount.ServiceAccount, error) {
	req, err := c.GCP.GetServiceAccount(ctx, &adminpb.GetServiceAccountRequest{
		Name: fmt.Sprintf("projects/%s/serviceAccounts/%s@%s.iam.gserviceaccount.com", id.Project, id.AccountId, id.Project),
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return serviceaccount.ServiceAccount{}, sdkerrors.NewErrorNotFound()
		}

		return serviceaccount.ServiceAccount{}, err
	}

	return serviceaccount.ServiceAccount{
		Identifier: id,
		Config: serviceaccount.Config{
			DisplayName: req.GetDisplayName(),
			Description: req.GetDescription(),
		},
		Attrs: serviceaccount.Attrs{
			UniqueId: req.GetUniqueId(),
			Disabled: req.GetDisabled(),
		},
	}, nil
}

func (c *client) CreateServiceAccount(ctx context.Context, id identifier.ServiceAccountIdentifier, config serviceaccount.Config) (serviceaccount.ServiceAccount, error) {
	res, err := c.GCP.CreateServiceAccount(ctx, &adminpb.CreateServiceAccountRequest{
		Name:      fmt.Sprintf("projects/%s", id.Project),
		AccountId: id.AccountId,
		ServiceAccount: &adminpb.ServiceAccount{
			DisplayName: config.DisplayName,
			Description: config.Description,
		},
	})
	if err != nil {
		return serviceaccount.ServiceAccount{}, err
	}

	return serviceaccount.ServiceAccount{
		Identifier: id,
		Config: serviceaccount.Config{
			DisplayName: res.GetDisplayName(),
			Description: res.GetDescription(),
		},
		Attrs: serviceaccount.Attrs{
			UniqueId: res.GetUniqueId(),
			Disabled: res.GetDisabled(),
		},
	}, nil
}

func (c *client) UpdateServiceAccount(ctx context.Context, id identifier.ServiceAccountIdentifier, config serviceaccount.Config, mask []value.UpdateMaskField) (serviceaccount.ServiceAccount, error) {
	res, err := c.GCP.UpdateServiceAccount(ctx, &adminpb.ServiceAccount{
		Name:        fmt.Sprintf("projects/%s/serviceAccounts/%s@%s.iam.gserviceaccount.com", id.Project, id.AccountId, id.Project),
		DisplayName: config.DisplayName,
		Description: config.Description,
	})
	if err != nil {
		return serviceaccount.ServiceAccount{}, err
	}

	return serviceaccount.ServiceAccount{
		Identifier: id,
		Config: serviceaccount.Config{
			DisplayName: res.GetDisplayName(),
			Description: res.GetDescription(),
		},
		Attrs: serviceaccount.Attrs{
			UniqueId: res.GetUniqueId(),
			Disabled: res.GetDisabled(),
		},
	}, nil
}

func (c *client) DeleteServiceAccount(ctx context.Context, id identifier.ServiceAccountIdentifier) error {
	return c.GCP.DeleteServiceAccount(ctx, &adminpb.DeleteServiceAccountRequest{
		Name: fmt.Sprintf("projects/%s/serviceAccounts/%s@%s.iam.gserviceaccount.com", id.Project, id.AccountId, id.Project),
	})
}
