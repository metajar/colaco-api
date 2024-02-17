package cmd

import (
	v1 "colaco-api/internal/api/v1"
	"context"
	"fmt"
	"log"
	"net/http"
)

func authenticate(client *v1.ClientWithResponses) (string, error) {
	resp, err := client.AuthLoginWithResponse(context.Background(), v1.AuthLoginJSONRequestBody{
		Username: username,
		Password: password,
	})
	if err != nil {
		return "", err
	}
	if resp.JSON200 != nil {
		return *resp.JSON200.Token, nil
	}
	log.Println("Authentication failed or did not return a token")
	return "", fmt.Errorf("authentication failed")
}

func addAuthHeader(ctx context.Context, req *http.Request, token string) error {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return nil
}
