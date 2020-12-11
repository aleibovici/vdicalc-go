package secretmanager

import (
	"context"
	"fmt"
	"log"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

// GetSecret Retrieve Citrix clientSecret from Google Secret Manager */
/* Secret string is split by ";" and fields by ":" */
/* customerID, clientID, clientSecret */
func GetSecret(projectID string, guserid string) map[string]string {

	// Create the client
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		fmt.Println("failed to create secretmanager client:", err)
	}

	// Build the request
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/" + projectID + "/secrets/" + guserid + "/versions/latest",
	}

	// Call the API
	data, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		fmt.Println("failed to get secret:", err)
		return nil
	}

	entries := strings.Split(string(data.Payload.Data), "; ")
	result := make(map[string]string)
	for _, e := range entries {
		parts := strings.Split(e, ":")
		result[parts[0]] = parts[1]
	}

	return result
}

// DeleteSecret Delete Citrix clientSecret from Google Secret Manager */
func DeleteSecret(projectID string, guserid string) error {

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		fmt.Println("failed to create secretmanager client:", err)
	}

	// Build the request.
	req := &secretmanagerpb.DeleteSecretRequest{
		Name: "projects/" + projectID + "/secrets/" + guserid,
	}

	// Call the API.
	if err := client.DeleteSecret(ctx, req); err != nil {
		return fmt.Errorf("failed to delete secret: %v", err)
	}

	return nil
}

// CreateSecret creates a new secret with the given name. A secret is a logical
// wrapper around a collection of secret versions. Secret versions hold the
// actual secret material.
func CreateSecret(projectID string, guserid string, customerID string, clientID string, clientSecret string) error {

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		fmt.Println("failed to create secretmanager client:", err)
	}

	// Build the request.
	createSecretReq := &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s", projectID),
		SecretId: guserid,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	}

	// Call the API.
	secret, err := client.CreateSecret(ctx, createSecretReq)
	if err != nil {
		log.Fatalf("failed to create secret: %v", err)
	}

	// Declare the payload to store.
	data := "customerID:" + customerID + "; clientID:" + clientID + "; clientSecret:" + clientSecret
	payload := []byte(data)

	// Build the request.
	addSecretVersionReq := &secretmanagerpb.AddSecretVersionRequest{
		Parent: secret.Name,
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	}

	// Call the API.
	_, err = client.AddSecretVersion(ctx, addSecretVersionReq)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
