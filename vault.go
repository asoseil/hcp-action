package main

import (
	"fmt"
	httptransport "github.com/go-openapi/runtime/client"
	resourcemanager "github.com/hashicorp/hcp-sdk-go/clients/cloud-resource-manager/stable/2019-12-10/client"
	"github.com/hashicorp/hcp-sdk-go/clients/cloud-resource-manager/stable/2019-12-10/client/project_service"
	hcpvaultsecrets "github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-secrets/stable/2023-11-28/client"
	"github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-secrets/stable/2023-11-28/client/secret_service"
	hcpconfig "github.com/hashicorp/hcp-sdk-go/config"
	hcpclient "github.com/hashicorp/hcp-sdk-go/httpclient"
	"os"
)

func hcpClient() (*httptransport.Runtime, error) {

	id := os.Getenv("HCP_CLIENT_ID")
	secret := os.Getenv("HCP_CLIENT_SECRET")

	hcpConfig, err := hcpconfig.NewHCPConfig(
		hcpconfig.WithClientCredentials(
			id,
			secret,
		),
		hcpconfig.WithoutBrowserLogin(),
	)
	if err != nil {
		return nil, err
	}

	cl, err := hcpclient.New(hcpclient.Config{
		HCPConfig: hcpConfig,
	})
	return cl, err
}

func vaultSecretsClient() (*hcpvaultsecrets.CloudVaultSecrets, error) {
	cl, err := hcpClient()
	if err != nil {
		return nil, err
	}

	return hcpvaultsecrets.New(cl, nil), nil
}

func resourceManagerClient() (*resourcemanager.CloudResourceManager, error) {
	cl, err := hcpClient()
	if err != nil {
		return nil, err
	}

	return resourcemanager.New(cl, nil), nil
}

func getOrganizationID(rmClient *resourcemanager.CloudResourceManager) (string, error) {
	orgID := os.Getenv("HCP_ORGANIZATION_ID")
	if len(orgID) > 0 {
		return orgID, nil
	}

	organizations, err := rmClient.OrganizationService.OrganizationServiceList(nil, nil)
	if err != nil {
		return "", err
	}

	orgName := os.Getenv("HCP_ORGANIZATION_NAME")
	if len(orgName) > 0 {
		for _, org := range organizations.Payload.Organizations {
			if org.Name == orgName {
				return org.ID, nil
			}
		}
	}

	if len(organizations.Payload.Organizations) == 0 {
		return "", fmt.Errorf("no organizations found")
	}
	return organizations.Payload.Organizations[0].ID, nil
}

func getProjectID(rmClient *resourcemanager.CloudResourceManager, organizationID string) (string, error) {
	projectID := os.Getenv("HCP_PROJECT_ID")
	if len(projectID) > 0 {
		return projectID, nil
	}

	scopeType := "ORGANIZATION"

	projects, err := rmClient.ProjectService.ProjectServiceList(
		project_service.NewProjectServiceListParams().WithScopeType(&scopeType).WithScopeID(&organizationID),
		nil,
	)
	if err != nil {
		return "", err
	}

	projectName := os.Getenv("HCP_PROJECT_NAME")
	if len(projectName) > 0 {
		for _, project := range projects.Payload.Projects {
			if project.Name == projectName {
				return project.ID, nil
			}
		}
	}
	if len(projects.Payload.Projects) == 0 {
		return "", fmt.Errorf("no projects found")
	}
	return projects.Payload.Projects[0].ID, nil
}

func getSecrets(vsClient *hcpvaultsecrets.CloudVaultSecrets, organization string, project string, app string, keys []string) (map[string]string, error) {
	secrets, err := vsClient.SecretService.OpenAppSecrets(
		secret_service.NewOpenAppSecretsParams().
			WithAppName(app).
			WithOrganizationID(organization).
			WithProjectID(project),
		nil,
	)
	if err != nil {
		return nil, err
	}

	secretsMap := make(map[string]string)

	keysMap := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		keysMap[key] = struct{}{}
	}

	// Filter the secrets by the keys provided
	for _, secret := range secrets.Payload.Secrets {
		if _, exists := keysMap[secret.Name]; exists {
			secretsMap[secret.Name] = secret.StaticVersion.Value
		}
	}

	return secretsMap, nil
}
