package auth

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/DillonEnge/jolt/internal/api"
	"github.com/alexedwards/scs/v2"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

//go:embed casdoor.cert
var cert string

var (
	ErrEnvarNotFound = errors.New("Envar not found")
)

type Client struct {
	*casdoorsdk.Client
}

func NewClient(config *api.Config) *Client {
	authConfig := &casdoorsdk.AuthConfig{
		Endpoint:         config.Casdoor.Endpoint,
		ClientId:         config.Casdoor.ClientID,
		ClientSecret:     config.Casdoor.ClientSecret,
		Certificate:      cert,
		OrganizationName: config.Casdoor.OrganizationName,
		ApplicationName:  config.Casdoor.ApplicationName,
	}

	casdoorClient := casdoorsdk.NewClientWithConf(authConfig)

	client := &Client{
		Client: casdoorClient,
	}

	return client
}

func (c *Client) GetClaims(ctx context.Context, sm *scs.SessionManager) (*casdoorsdk.Claims, error) {
	token := sm.GetString(ctx, "authToken")
	if token != "" {
		return c.ParseJwtToken(token)
	}

	return nil, fmt.Errorf("failed to find authToken in session")
}
