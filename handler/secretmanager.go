package handler

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

const (
	defaultSecretVersion = "latest"
)

func accessSecretVersion(projectID, secretID, versionID string) (string, error) {
	if projectID == "" {
		return "", fmt.Errorf("projectID is required")
	}
	if secretID == "" {
		return "", fmt.Errorf("secretID is required")
	}
	if versionID == "" {
		return "", fmt.Errorf("versionID is required")
	}

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/%s", projectID, secretID, versionID),
	}

	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}

	return string(result.Payload.Data), nil
}
