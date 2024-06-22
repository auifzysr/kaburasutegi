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

type SecretManagerConfig struct {
	projectID string
	secretID  string
	versionID string
}

type SecretManagerOption func(*SecretManagerConfig)

func WithProjectID(projectID string) SecretManagerOption {
	return func(config *SecretManagerConfig) {
		config.projectID = projectID
	}
}

func WithSecretID(secretID string) SecretManagerOption {
	return func(config *SecretManagerConfig) {
		config.secretID = secretID
	}
}

func WithVersionID(versionID string) SecretManagerOption {
	return func(config *SecretManagerConfig) {
		config.versionID = versionID
	}
}

func accessSecretVersion(s *SecretManagerConfig) (string, error) {
	if s.projectID == "" {
		return "", fmt.Errorf("projectID is required")
	}
	if s.secretID == "" {
		return "", fmt.Errorf("secretID is required")
	}
	if s.versionID == "" {
		return "", fmt.Errorf("versionID is required")
	}

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/%s", s.projectID, s.secretID, s.versionID),
	}

	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}

	return string(result.Payload.Data), nil
}
