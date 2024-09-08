package auth

import (
	_ "embed"
	"errors"

	"github.com/DillonEnge/jolt/internal/api"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

//go:embed casdoor.cert
var cert string

var (
	ErrEnvarNotFound = errors.New("Envar not found")
)

func NewClient(config *api.Config) *casdoorsdk.Client {
	authConfig := &casdoorsdk.AuthConfig{
		Endpoint:         config.Casdoor.Endpoint,
		ClientId:         config.Casdoor.ClientID,
		ClientSecret:     config.Casdoor.ClientSecret,
		Certificate:      cert,
		OrganizationName: config.Casdoor.OrganizationName,
		ApplicationName:  config.Casdoor.ApplicationName,
	}

	casdoorClient := casdoorsdk.NewClientWithConf(authConfig)

	return casdoorClient
}
