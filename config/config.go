package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	appEnv         = "APP_ENV"
	servicePort    = "APP_SERVICE_PORT"
	databaseName   = "DATABASE_NAME"
	databaseURL    = "DATABASE_URL"
	serviceName    = "SERVICE_NAME"
	kafkaBrokers   = "KAFKA_BROKERS"
	googleClientID = "GOOGLE_CLIENT_ID"
	googleSecret   = "GOOGLE_SECRET"
	googleApiKey   = "API_KEY"
	redirectUrl    = "REDIRECT_URL"
)

type source interface {
	GetEnv(key string, fallback string) string
	GetEnvBool(key string, fallback bool) bool
	GetEnvInt(key string, fallback int) int
}

type OSSource struct {
	source //nolint
}

func (o OSSource) GetEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func (o OSSource) GetEnvBool(key string, fallback bool) bool {
	b := o.GetEnv(key, "")
	if len(b) == 0 {
		return fallback
	}
	v, err := strconv.ParseBool(b)
	if err != nil {
		return fallback
	}
	return v
}

func (o OSSource) GetEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		result, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}
		return result
	}
	return fallback
}

type Config struct {
	AppEnv             string
	ServicePort        string
	DatabaseName       string
	DatabaseURL        string
	AutoReload         bool
	ServiceName        string
	KafkaBrokers       []string
	GoogleApiKey       string
	GoogleClientID     string
	GoogleClientSecret string
	RedirectUrl        string
}

func ImportConfig(source source) Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	appEnv := source.GetEnv(appEnv, "")
	port := source.GetEnv(servicePort, "8001")
	databaseName := source.GetEnv(databaseName, "appservices")
	databaseURL := source.GetEnv(databaseURL, "mongodb://127.0.0.1:27017")
	serviceName := source.GetEnv(serviceName, "rescue-support.sv")
	googleApiKey := source.GetEnv(googleApiKey, "")
	googleClientID := source.GetEnv(googleClientID, "")
	googleClientSecret := source.GetEnv(googleSecret, "")
	redirectUrl := source.GetEnv(redirectUrl, "")

	return Config{
		AppEnv:             appEnv,
		ServicePort:        port,
		DatabaseName:       databaseName,
		DatabaseURL:        databaseURL,
		ServiceName:        serviceName,
		KafkaBrokers:       []string{"localhost:9092"},
		GoogleApiKey:       googleApiKey,
		GoogleClientID:     googleClientID,
		GoogleClientSecret: googleClientSecret,
		RedirectUrl:        redirectUrl,
	}
}
