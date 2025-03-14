package api

import (
	"os"
	"strconv"

	"github.com/nats-io/nats.go"
)

type Config struct {
	DBUrl     string
	Port      int
	NatsURL   string
	Casdoor   CasdoorConfig
	SeaweedFS SeaweedFSConfig
}

type CasdoorConfig struct {
	Endpoint         string
	ClientID         string
	ClientSecret     string
	RedirectURI      string
	OrganizationName string
	ApplicationName  string
}

type SeaweedFSConfig struct {
	MasterURL  string
	VolumesURL string
}

func NewConfig() *Config {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}

	natsURL, ok := os.LookupEnv("NATS_URL")
	if !ok {
		natsURL = nats.DefaultURL
	}

	return &Config{
		DBUrl:   os.Getenv("DATABASE_URL"),
		Port:    port,
		NatsURL: natsURL,
		Casdoor: CasdoorConfig{
			Endpoint:         os.Getenv("CASDOOR_ENDPOINT"),
			ClientID:         os.Getenv("CASDOOR_CLIENT_ID"),
			ClientSecret:     os.Getenv("CASDOOR_CLIENT_SECRET"),
			RedirectURI:      os.Getenv("CASDOOR_REDIRECT_URI"),
			OrganizationName: os.Getenv("CASDOOR_ORGANIZATION_NAME"),
			ApplicationName:  os.Getenv("CASDOOR_APPLICATION_NAME"),
		},
		SeaweedFS: SeaweedFSConfig{
			MasterURL:  os.Getenv("SEAWEEDFS_MASTER_URL"),
			VolumesURL: os.Getenv("SEAWEEDFS_VOLUMES_URL"),
		},
	}
}
