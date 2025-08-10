package service

import (
	"fmt"
	"log"
	"sync"
	"time"
	"github.com/gorilla/websocket"
)

type PushService struct {
	clients map[string]*websocket.Conn
	mu      sync.RWMutex
}

func NewPushService() *PushService {
	return &PushService{
		clients: make(map[string]*websocket.Conn),
	}
}

func (p *PushService) RegisterClient(userID string, conn *websocket.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// close previous if exists
	if old, ok := p.clients[userID]; ok {
		_ = old.Close()
	}

	p.clients[userID] = conn
	log.Printf("Registered WS client for userID=%s", userID)
}

func (p *PushService) UnregisterClient(userID string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if conn, ok := p.clients[userID]; ok {
		_ = conn.Close()
		delete(p.clients, userID)
		log.Printf("Unregistered WS client for userID=%s", userID)
	}
}

func (p *PushService) Send(userID, subject, message string) error {
	p.mu.RLock()
	conn, ok := p.clients[userID]
	p.mu.RUnlock()

	if !ok || conn == nil {
		return fmt.Errorf("user %s is not connected to WebSocket", userID)
	}

	payload := map[string]string{
		"type":    "push",
		"subject": subject,
		"message": message,
	}

	_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err := conn.WriteJSON(payload); err != nil {
		log.Printf("error writing to %s: %v — unregistering", userID, err)
		p.UnregisterClient(userID)
		return err
	}

	return nil
}