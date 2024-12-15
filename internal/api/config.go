package api

import (
	"os"
	"strconv"
)

type Config struct {
	DBUrl   string
	Port    int
	Casdoor CasdoorConfig
	Minio   MinioConfig
}

type CasdoorConfig struct {
	Endpoint         string
	ClientID         string
	ClientSecret     string
	RedirectURI      string
	OrganizationName string
	ApplicationName  string
}

type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

func NewConfig() *Config {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}

	return &Config{
		DBUrl: os.Getenv("DATABASE_URL"),
		Port:  port,
		Casdoor: CasdoorConfig{
			Endpoint:         os.Getenv("CASDOOR_ENDPOINT"),
			ClientID:         os.Getenv("CASDOOR_CLIENT_ID"),
			ClientSecret:     os.Getenv("CASDOOR_CLIENT_SECRET"),
			RedirectURI:      os.Getenv("CASDOOR_REDIRECT_URI"),
			OrganizationName: os.Getenv("CASDOOR_ORGANIZATION_NAME"),
			ApplicationName:  os.Getenv("CASDOOR_APPLICATION_NAME"),
		},
		Minio: MinioConfig{
			Endpoint:        os.Getenv("MINIO_ENDPOINT"),
			AccessKeyID:     os.Getenv("MINIO_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("MINIO_SECRET_ACCESS_KEY"),
			UseSSL:          false,
		},
	}
}
