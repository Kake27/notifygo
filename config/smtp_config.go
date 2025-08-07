package config

import (
    "fmt"
    "os"
)

type SMTPConfig struct {
    Host     string
    Port     string
    Sender   string
    Username string
    Password string
}

func GetSMTPConfig() SMTPConfig {
    return SMTPConfig{
        Host:     getEnv("SMTP_HOST", "localhost"),  
        Port:     getEnv("SMTP_PORT", "586"),
        Sender:   getEnv("SMTP_SENDER", "test@example.com"),
        Username: getEnv("SMTP_USERNAME", "tester"),
        Password: getEnv("SMTP_PASSWORD", "123"),
    }
}


func (c SMTPConfig) Address() string {
    return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}
