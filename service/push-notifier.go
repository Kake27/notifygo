package service

import (
	"fmt"
	"log"
	"sync"

	"gofr.dev/pkg/gofr"
)

type PushService struct {
	clients map[string]*gofr.Context
	mu      sync.Mutex
}

func NewPushService() *PushService {
	return &PushService{
		clients: make(map[string]*gofr.Context),
	}
}

func (p *PushService) RegisterClient(clientID string, ctx *gofr.Context) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.clients[clientID] = ctx
	log.Printf("Registed WebSocket client for clientID: %s", clientID)
}

func (p *PushService) Send(to, subject, message string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	clientCtx, ok := p.clients[to]
	if !ok {
		return fmt.Errorf("user %s is not connected to WebSocket", to)
	}

	payload := map[string]string{
		"type":    "push",
		"subject": subject,
		"message": message,
	}

	return clientCtx.WriteMessageToSocket(payload)
}