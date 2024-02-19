package iam_policy

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	iampolicy "github.com/alchematik/athanor-provider-gcp/gen/provider/iam_policy"
	"github.com/alchematik/athanor-provider-gcp/gen/provider/identifier"

	"cloud.google.com/go/iam/apiv1/iampb"
	// resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	// cloudrun "cloud.google.com/go/run/apiv2"

	cloudfunction "cloud.google.com/go/functions/apiv2"
	sdkerrors "github.com/alchematik/athanor-go/sdk/errors"
	"github.com/alchematik/athanor-go/sdk/provider/value"
	"github.com/googleapis/gax-go/v2"
)

var (
	customServiceAccountRe = regexp.MustCompile(`serviceAccount:(.+)@(.+).iam.gserviceaccount.com`)
)

func NewHandler(ctx context.Context) (*iampolicy.IamPolicyHandler, error) {
	// cloudRun, err := cloudrun.NewServicesClient(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	fc, err := cloudfunction.NewFunctionRESTClient(ctx)
	if err != nil {
		return nil, err
	}

	c := &client{
		CloudFunction: fc,
	}

	return &iampolicy.IamPolicyHandler{
		IamPolicyCreator: c,
		IamPolicyDeleter: c,
		IamPolicyGetter:  c,
		IamPolicyUpdator: c,
	}, nil
}

type GCP interface {
	GetIamPolicy(context.Context, *iampb.GetIamPolicyRequest, ...gax.CallOption) (*iampb.Policy, error)
	SetIamPolicy(context.Context, *iampb.SetIamPolicyRequest, ...gax.CallOption) (*iampb.Policy, error)
}

type client struct {
	CloudFunction GCP
}

func (c *client) GetIamPolicy(ctx context.Context, id identifier.IamPolicyIdentifier) (iampolicy.IamPolicy, error) {
	switch resourceID := id.Resource.(type) {
	case identifier.FunctionIdentifier:
		resource := fmt.Sprintf("projects/%s/locations/%s/functions/%s", resourceID.Project, resourceID.Location, resourceID.Name)
		res, err := c.CloudFunction.GetIamPolicy(ctx, &iampb.GetIamPolicyRequest{
			Resource: resource,
		})
		if err != nil {
			return iampolicy.IamPolicy{}, err
		}

		if len(res.Bindings) == 0 {
			return iampolicy.IamPolicy{}, sdkerrors.NewErrorNotFound()
		}

		bindings := make([]iampolicy.Binding, len(res.Bindings))
		for i, binding := range res.Bindings {
			members := make([]value.ResourceIdentifier, len(binding.Members))
			for j, member := range binding.Members {
				// TODO: Support "allUsers" and "allAuthenticatedUsers"
				parts := strings.Split(member, ":")
				if len(parts) == 0 {
					return iampolicy.IamPolicy{}, fmt.Errorf("unsupported member type: %s", member)
				}
				switch parts[0] {
				case "serviceAccount":
					// TODO: support default service account. Probably makes sense to create a different resource identifier type for those.
					customServiceAccountParts := customServiceAccountRe.FindStringSubmatch(member)
					if len(customServiceAccountParts) < 2 {
						return iampolicy.IamPolicy{}, fmt.Errorf("invalid custom service account email: %s", member)
					}

					log.Printf("MEMBER >>>>>> %v\n", member)

					members[j] = identifier.ServiceAccountIdentifier{
						AccountId: customServiceAccountParts[1],
						Project:   customServiceAccountParts[2],
					}
				default:
					return iampolicy.IamPolicy{}, fmt.Errorf("unsupported member type: %s", member)
				}
			}

			var r value.ResourceIdentifier
			switch {
			case strings.HasPrefix(binding.Role, "roles/"):
				parts := strings.Split(binding.Role, "/")
				if len(parts) < 2 {
					return iampolicy.IamPolicy{}, fmt.Errorf("invalid pre-defined role: %s", binding.Role)
				}

				r = identifier.IamRoleIdentifier{
					Name: parts[1],
				}
			case strings.HasPrefix(binding.Role, "projects/"):
				parts := strings.Split(binding.Role, "/")
				if len(parts) < 4 {
					return iampolicy.IamPolicy{}, fmt.Errorf("invalid custom project level role: %s", binding.Role)
				}
				r = identifier.IamRoleCustomProjectIdentifier{
					Project: parts[1],
					Name:    parts[3],
				}
			}

			bindings[i] = iampolicy.Binding{
				Role:    r,
				Members: members,
			}
		}

		return iampolicy.IamPolicy{
			Identifier: id,
			Config: iampolicy.Config{
				Bindings: bindings,
			},
			Attrs: iampolicy.Attrs{
				Etag: fmt.Sprintf("%x", res.GetEtag()),
			},
		}, nil
	default:
		return iampolicy.IamPolicy{}, fmt.Errorf("unsupported identifier type: %T", id)
	}
}

func (c *client) CreateIamPolicy(ctx context.Context, id identifier.IamPolicyIdentifier, config iampolicy.Config) (iampolicy.IamPolicy, error) {
	log.Printf("CREATING IAM POLICY>>>>>>>>>>>")
	switch resourceID := id.Resource.(type) {
	case identifier.FunctionIdentifier:
		bindings := make([]*iampb.Binding, len(config.Bindings))
		for i, b := range config.Bindings {
			var roleID string
			switch role := b.Role.(type) {
			case identifier.IamRoleCustomProjectIdentifier:
				roleID = fmt.Sprintf("projects/%s/roles/%s", role.Project, role.Name)
			default:
				return iampolicy.IamPolicy{}, fmt.Errorf("unsupported role type: %T", b.Role)
			}

			members := make([]string, len(b.Members))
			for i, m := range b.Members {
				switch m := m.(type) {
				case identifier.ServiceAccountIdentifier:
					members[i] = fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", m.AccountId, m.Project)
				default:
					return iampolicy.IamPolicy{}, fmt.Errorf("unsupported member type: %T", m)
				}
			}

			bindings[i] = &iampb.Binding{
				Role:    roleID,
				Members: members,
			}
		}

		res, err := c.CloudFunction.SetIamPolicy(ctx, &iampb.SetIamPolicyRequest{
			Resource: fmt.Sprintf("projects/%s/locations/%s/functions/%s", resourceID.Project, resourceID.Location, resourceID.Name),
			Policy: &iampb.Policy{
				Version:  3,
				Bindings: bindings,
			},
		})
		if err != nil {
			return iampolicy.IamPolicy{}, err
		}

		binds := make([]iampolicy.Binding, len(res.Bindings))
		for i, binding := range res.Bindings {
			members := make([]value.ResourceIdentifier, len(binding.Members))
			for j, member := range binding.Members {
				// TODO: Support "allUsers" and "allAuthenticatedUsers"
				parts := strings.Split(member, ":")
				if len(parts) == 0 {
					return iampolicy.IamPolicy{}, fmt.Errorf("unsupported member type: %s", member)
				}
				switch parts[0] {
				case "serviceAccount":
					// TODO: support default service account. Probably makes sense to create a different resource identifier type for those.
					customServiceAccountParts := customServiceAccountRe.FindStringSubmatch(member)
					if len(customServiceAccountParts) < 2 {
						return iampolicy.IamPolicy{}, fmt.Errorf("invalid custom service account email: %s", member)
					}

					members[j] = identifier.ServiceAccountIdentifier{
						AccountId: customServiceAccountParts[1],
						Project:   customServiceAccountParts[2],
					}
				default:
					return iampolicy.IamPolicy{}, fmt.Errorf("unsupported member type: %s", member)
				}
			}

			var r value.ResourceIdentifier
			switch {
			case strings.HasPrefix(binding.Role, "roles/"):
				parts := strings.Split(binding.Role, "/")
				if len(parts) < 2 {
					return iampolicy.IamPolicy{}, fmt.Errorf("invalid pre-defined role: %s", binding.Role)
				}

				r = identifier.IamRoleIdentifier{
					Name: parts[1],
				}
			case strings.HasPrefix(binding.Role, "projects/"):
				parts := strings.Split(binding.Role, "/")
				if len(parts) < 4 {
					return iampolicy.IamPolicy{}, fmt.Errorf("invalid custom project level role: %s", binding.Role)
				}
				r = identifier.IamRoleCustomProjectIdentifier{
					Project: parts[1],
					Name:    parts[3],
				}
			}

			binds[i] = iampolicy.Binding{
				Role:    r,
				Members: members,
			}
		}

		return iampolicy.IamPolicy{
			Identifier: id,
			Config: iampolicy.Config{
				Bindings: binds,
			},
			Attrs: iampolicy.Attrs{
				Etag: fmt.Sprintf("%x", res.GetEtag()),
			},
		}, nil
	default:
		return iampolicy.IamPolicy{}, fmt.Errorf("unsupported identifier type: %T", id)
	}
}

func (c *client) UpdateIamPolicy(ctx context.Context, id identifier.IamPolicyIdentifier, config iampolicy.Config, mask []value.UpdateMaskField) (iampolicy.IamPolicy, error) {
	log.Printf("UPDATING IAM POLICY>>>>>>>>>>>")
	switch resourceID := id.Resource.(type) {
	case identifier.FunctionIdentifier:
		resource := fmt.Sprintf("projects/%s/locations/%s/functions/%s", resourceID.Project, resourceID.Location, resourceID.Name)
		getRes, err := c.CloudFunction.GetIamPolicy(ctx, &iampb.GetIamPolicyRequest{
			Resource: resource,
		})
		if err != nil {
			return iampolicy.IamPolicy{}, err
		}

		bindings := make([]*iampb.Binding, len(config.Bindings))
		for i, b := range config.Bindings {
			var roleID string
			switch role := b.Role.(type) {
			case identifier.IamRoleCustomProjectIdentifier:
				roleID = fmt.Sprintf("projects/%s/roles/%s", role.Project, role.Name)
			default:
				return iampolicy.IamPolicy{}, fmt.Errorf("unsupported role type: %T", b.Role)
			}

			members := make([]string, len(b.Members))
			for i, m := range b.Members {
				switch m := m.(type) {
				case identifier.ServiceAccountIdentifier:
					members[i] = fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", m.AccountId, m.Project)
				default:
					return iampolicy.IamPolicy{}, fmt.Errorf("unsupported member type: %T", m)
				}
			}

			bindings[i] = &iampb.Binding{
				Role:    roleID,
				Members: members,
			}
		}

		res, err := c.CloudFunction.SetIamPolicy(ctx, &iampb.SetIamPolicyRequest{
			Resource: resource,
			Policy: &iampb.Policy{
				Version:  3,
				Bindings: bindings,
				Etag:     getRes.GetEtag(),
			},
		})
		if err != nil {
			return iampolicy.IamPolicy{}, err
		}

		binds := make([]iampolicy.Binding, len(res.Bindings))
		for i, binding := range res.Bindings {
			members := make([]value.ResourceIdentifier, len(binding.Members))
			for j, member := range binding.Members {
				// TODO: Support "allUsers" and "allAuthenticatedUsers"
				parts := strings.Split(member, ":")
				if len(parts) == 0 {
					return iampolicy.IamPolicy{}, fmt.Errorf("unsupported member type: %s", member)
				}
				switch parts[0] {
				case "serviceAccount":
					// TODO: support default service account. Probably makes sense to create a different resource identifier type for those.
					customServiceAccountParts := customServiceAccountRe.FindStringSubmatch(member)
					if len(customServiceAccountParts) < 2 {
						return iampolicy.IamPolicy{}, fmt.Errorf("invalid custom service account email: %s", member)
					}

					members[j] = identifier.ServiceAccountIdentifier{
						AccountId: customServiceAccountParts[1],
						Project:   customServiceAccountParts[2],
					}
				default:
					return iampolicy.IamPolicy{}, fmt.Errorf("unsupported member type: %s", member)
				}
			}

			var r value.ResourceIdentifier
			switch {
			case strings.HasPrefix(binding.Role, "roles/"):
				parts := strings.Split(binding.Role, "/")
				if len(parts) < 2 {
					return iampolicy.IamPolicy{}, fmt.Errorf("invalid pre-defined role: %s", binding.Role)
				}

				r = identifier.IamRoleIdentifier{
					Name: parts[1],
				}
			case strings.HasPrefix(binding.Role, "projects/"):
				parts := strings.Split(binding.Role, "/")
				if len(parts) < 4 {
					return iampolicy.IamPolicy{}, fmt.Errorf("invalid custom project level role: %s", binding.Role)
				}
				r = identifier.IamRoleCustomProjectIdentifier{
					Project: parts[1],
					Name:    parts[3],
				}
			}

			binds[i] = iampolicy.Binding{
				Role:    r,
				Members: members,
			}
		}

		return iampolicy.IamPolicy{
			Identifier: id,
			Config: iampolicy.Config{
				Bindings: binds,
			},
			Attrs: iampolicy.Attrs{
				Etag: fmt.Sprintf("%x", res.GetEtag()),
			},
		}, nil
	default:
		return iampolicy.IamPolicy{}, fmt.Errorf("unsupported identifier type: %T", id)
	}
}

func (c *client) DeleteIamPolicy(ctx context.Context, id identifier.IamPolicyIdentifier) error {
	switch resourceID := id.Resource.(type) {
	case identifier.FunctionIdentifier:
		resource := fmt.Sprintf("projects/%s/locations/%s/functions/%s", resourceID.Project, resourceID.Location, resourceID.Name)
		getRes, err := c.CloudFunction.GetIamPolicy(ctx, &iampb.GetIamPolicyRequest{
			Resource: resource,
		})
		if err != nil {
			return err
		}

		_, err = c.CloudFunction.SetIamPolicy(ctx, &iampb.SetIamPolicyRequest{
			Resource: resource,
			Policy: &iampb.Policy{
				Version:  3,
				Bindings: nil,
				Etag:     getRes.GetEtag(),
			},
		})
		return err
	default:
		return fmt.Errorf("invalid identifier type: %T", resourceID)
	}
}
