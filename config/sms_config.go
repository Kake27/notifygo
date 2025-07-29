package config

import (
	"os"
)


type SMSConfig struct {
    GatewayURL      string
    Username string
    Password string
	Sender string
}

func GetSMSConfig() SMSConfig {
    return SMSConfig{
        GatewayURL:      os.Getenv("KANNEL_URL"),
        Username: os.Getenv("KANNEL_USERNAME"),
        Password: os.Getenv("KANNEL_PASSWORD"),
		Sender: os.Getenv("KANNEL_SENDER"),
    }
}
