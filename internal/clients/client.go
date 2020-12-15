package clients

import (
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/manicminer/hamilton/auth"

	"github.com/hashicorp/terraform-provider-azuread/internal/common"
	applications "github.com/hashicorp/terraform-provider-azuread/internal/services/applications/client"
	domains "github.com/hashicorp/terraform-provider-azuread/internal/services/domains/client"
	groups "github.com/hashicorp/terraform-provider-azuread/internal/services/groups/client"
	serviceprincipals "github.com/hashicorp/terraform-provider-azuread/internal/services/serviceprincipals/client"
	users "github.com/hashicorp/terraform-provider-azuread/internal/services/users/client"
)

// Client contains the handles to all the specific Azure AD resource classes' respective clients
type Client struct {
	ClientID string
	ObjectID string
	TenantID string
	Claims   auth.Claims

	TerraformVersion string
	Environment      azure.Environment

	AuthenticatedAsAServicePrincipal bool
	EnableMsGraphBeta                bool

	StopContext context.Context

	Applications      *applications.Client
	Domains           *domains.Client
	Groups            *groups.Client
	ServicePrincipals *serviceprincipals.Client
	Users             *users.Client
}

func (client *Client) build(ctx context.Context, o *common.ClientOptions) error { //nolint:unparam
	autorest.Count429AsRetry = false
	client.StopContext = ctx

	client.Applications = applications.NewClient(o)
	client.Domains = domains.NewClient(o)
	client.Groups = groups.NewClient(o)
	client.ServicePrincipals = serviceprincipals.NewClient(o)
	client.Users = users.NewClient(o)

	if client.EnableMsGraphBeta {
		// Acquire an access token upfront so we can decode and populate the JWT claims
		token, err := o.MsGraphAuthorizer.Token()
		if err != nil {
			return fmt.Errorf("unable to obtain access token: %v", err)
		}
		client.Claims, err = auth.ParseClaims(token)
		if err != nil {
			return fmt.Errorf("unable to parse claims in access token: %v", err)
		}
	}

	return nil
}
