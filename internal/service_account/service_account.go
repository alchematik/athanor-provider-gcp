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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewHandler() serviceaccount.ServiceAccountHandler {
	c := &client{}
	return serviceaccount.ServiceAccountHandler{
		ServiceAccountGetter:  c,
		ServiceAccountCreator: c,
		ServiceAccountUpdator: c,
		ServiceAccountDeleter: c,
	}
}

type client struct{}

func (c *client) GetServiceAccount(ctx context.Context, id identifier.ServiceAccountIdentifier) (serviceaccount.ServiceAccount, error) {
	iamClient, err := iamadmin.NewIamClient(ctx)
	if err != nil {
		return serviceaccount.ServiceAccount{}, err
	}

	req, err := iamClient.GetServiceAccount(ctx, &adminpb.GetServiceAccountRequest{
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
		Config: serviceaccount.ServiceAccountConfig{
			DisplayName: req.GetDisplayName(),
			Description: req.GetDescription(),
		},
		Attrs: serviceaccount.ServiceAccountAttrs{
			UniqueId: req.GetUniqueId(),
			Disabled: req.GetDisabled(),
		},
	}, nil
}

func (c *client) CreateServiceAccount(ctx context.Context, id identifier.ServiceAccountIdentifier, config serviceaccount.ServiceAccountConfig) (serviceaccount.ServiceAccount, error) {
	iamClient, err := iamadmin.NewIamClient(ctx)
	if err != nil {
		return serviceaccount.ServiceAccount{}, err
	}

	res, err := iamClient.CreateServiceAccount(ctx, &adminpb.CreateServiceAccountRequest{
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
		Config: serviceaccount.ServiceAccountConfig{
			DisplayName: res.GetDisplayName(),
			Description: res.GetDescription(),
		},
		Attrs: serviceaccount.ServiceAccountAttrs{
			UniqueId: res.GetUniqueId(),
			Disabled: res.GetDisabled(),
		},
	}, nil
}

func (c *client) UpdateServiceAccount(ctx context.Context, id identifier.ServiceAccountIdentifier, config serviceaccount.ServiceAccountConfig, mask []value.UpdateMaskField) (serviceaccount.ServiceAccount, error) {
	iamClient, err := iamadmin.NewIamClient(ctx)
	if err != nil {
		return serviceaccount.ServiceAccount{}, err
	}

	res, err := iamClient.UpdateServiceAccount(ctx, &adminpb.ServiceAccount{
		Name:        fmt.Sprintf("projects/%s/serviceAccounts/%s@%s.iam.gserviceaccount.com", id.Project, id.AccountId, id.Project),
		DisplayName: config.DisplayName,
		Description: config.Description,
	})
	if err != nil {
		return serviceaccount.ServiceAccount{}, err
	}

	return serviceaccount.ServiceAccount{
		Identifier: id,
		Config: serviceaccount.ServiceAccountConfig{
			DisplayName: res.GetDisplayName(),
			Description: res.GetDescription(),
		},
		Attrs: serviceaccount.ServiceAccountAttrs{
			UniqueId: res.GetUniqueId(),
			Disabled: res.GetDisabled(),
		},
	}, nil
}

func (c *client) DeleteServiceAccount(ctx context.Context, id identifier.ServiceAccountIdentifier) error {
	iamClient, err := iamadmin.NewIamClient(ctx)
	if err != nil {
		return err
	}

	return iamClient.DeleteServiceAccount(ctx, &adminpb.DeleteServiceAccountRequest{
		Name: fmt.Sprintf("projects/%s/serviceAccounts/%s@%s.iam.gserviceaccount.com", id.Project, id.AccountId, id.Project),
	})
}
