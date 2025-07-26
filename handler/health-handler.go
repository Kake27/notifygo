package handler

import (
	"net"
	"time"
	"gofr.dev/pkg/gofr"
	"notification-service/config"
)

type HealthHandler struct {}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(ctx *gofr.Context) (interface{}, error){
	smtpCfg := config.GetSMTPConfig()

	conn, err := net.DialTimeout("tcp", smtpCfg.Address(), 2*time.Second)
	if err!=nil {
		return map[string]string {
			"status": "DOWN",
			"smtp_status": "unreachable",
		}, nil
	}
	_ = conn.Close()

	return map[string]string{
		"status": "UP",
		"smtp_status": "reachable",
	}, nil
}