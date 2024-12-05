package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	var target string

	// Args follow the order declared in action file
	_ = os.Setenv("HCP_CLIENT_ID", os.Args[1])
	_ = os.Setenv("HCP_CLIENT_SECRET", os.Args[2])

	if len(os.Args) > 4 {
		target = os.Args[4]
	}

	config, err := parseConfig(os.Args[3])
	if err != nil {
		log.Fatalf("failed to parse config: %s", err)
	}

	// Setup HCP client
	rmClient, err := resourceManagerClient()
	if err != nil {
		log.Fatalf("Error creating resource manager client: %v", err)
	}

	organizationID, err := getOrganizationID(rmClient)
	if err != nil {
		log.Fatalf("Error getting organization ID: %v", err)
	}

	projectID, err := getProjectID(rmClient, organizationID)
	if err != nil {
		log.Fatalf("Error getting project ID: %v", err)
	}

	// Print loaded HCP config
	fmt.Println("Organization ID:", organizationID)
	fmt.Println("Project ID:", projectID)

	// Setup Vault Secrets Client
	client, err := vaultSecretsClient()
	if err != nil {
		log.Fatalf("Error creating vault secrets client: %v", err)
	}

	// Loop app and secrets given as arg
	for app, secrets := range config {
		values, err := getSecrets(client, organizationID, projectID, app, secrets)
		if err != nil {
			log.Fatalf("Error getting secrets: %v", err)
		}

		for key, value := range values {
			// Check target and expose
			switch target {
			case TargetEnv:
				setEnvVar(key, value)
			case TargetOut:
				setOutVar(key, value)
			default:
				setOutVar(key, value)
				setEnvVar(key, value)
			}
			fmt.Printf("- %s: %s\n", key, target)
		}
	}
}
