package handler

import (
	"fmt"
	"notification-service/service"
	"notification-service/store"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"
)

// NotificationHandler handles notification requests
type NotificationHandler struct {
	emailService     *service.EmailService
	smsService       *service.SMSService
	pushService      *service.PushService
	templateRenderer *service.TemplateRenderer
	dbTemplateStore  *store.DBTemplateStore
}


func NewNotificationHandler(e *service.EmailService, s *service.SMSService, p *service.PushService, tr *service.TemplateRenderer, dbStore *store.DBTemplateStore) *NotificationHandler {
	return &NotificationHandler{
		emailService:     e,
		smsService:       s,
		pushService:      p,
		templateRenderer: tr,
		dbTemplateStore:  dbStore,
	}
}

type NotificationRequest struct {
	Type      string                 `json:"type"`      // type of notification
	To        string                 `json:"to"`        // recipient
	Subject   string                 `json:"subject"`   // subject/title
	Message   string                 `json:"message"`   // message content
	Template  string                 `json:"template"`  // template name
	Vars      map[string]string      `json:"vars"`      // template variables

}

func (h *NotificationHandler) Notify(ctx *gofr.Context) (interface{}, error) {
	var req NotificationRequest

	if err := ctx.Bind(&req); err != nil {
		return nil, http.ErrorInvalidParam{Params: []string{"request body"}}
	}

	// checks if any missing params
	missing := []string{}
	if req.To == "" {
		missing = append(missing, "to")
	}
	if req.Message == "" && req.Template == "" {
		missing = append(missing, "message or template")
	}
	if req.Type == "" {
		missing = append(missing, "type")
	}

	if len(missing) > 0 {
		return nil, http.ErrorMissingParam{Params: missing}
	}

	// templating
	var finalMessage string
	if req.Template != "" {
		rendered, err := h.templateRenderer.Render(req.Template, req.Vars)
		if err != nil {
			return nil, fmt.Errorf("failed to render template %w", err)
		}

		finalMessage = rendered
	} else {
		finalMessage = req.Message
	}

	// send the notification
	var err error
	switch req.Type {
	case "email":
		err = h.emailService.Send(req.To, req.Subject, finalMessage)
		if err != nil {
			return nil, fmt.Errorf("failed to send email: %w", err)
		}
		return map[string]string{"status": "email sent successfully"}, nil

	case "sms":
		_, err = h.smsService.Send(req.To, finalMessage)
		if err != nil {
			return nil, fmt.Errorf("failed to send SMS: %w", err)
		}

		return map[string]string{"status": "SMS sent successfully"}, nil

	case "push":

		err = h.pushService.Send(req.To, req.Subject, finalMessage)
		if err != nil {
			return nil, fmt.Errorf("failed to send push notification: %w", err)
		}
		return map[string]string{"status": "push notification sent successfully"}, nil

	default:
		return nil, http.ErrorInvalidParam{Params: []string{"type"}}
	}
}

// creating and deleting templates

func (h *NotificationHandler) CreateTemplate(ctx *gofr.Context) (interface{}, error) {
	var req struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	}

	if err := ctx.Bind(&req); err != nil {
		return nil, http.ErrorInvalidParam{Params: []string{"name", "content"}}
	}

	err := h.dbTemplateStore.Create(req.Name, req.Content)
	if err != nil {
		return nil, fmt.Errorf("template creation failed: %w", err)
	}

	return map[string]string{"status": "template created"}, nil
}

func (h *NotificationHandler) DeleteTemplate(ctx *gofr.Context) (interface{}, error) {
	name := ctx.PathParam("name")
	if name == "" {
		return nil, http.ErrorMissingParam{Params: []string{"name"}}
	}

	err := h.dbTemplateStore.Delete(name)
	if err != nil {
		return nil, fmt.Errorf("template deletion failed: %w", err)
	}

	return map[string]string{"status": "template deleted"}, nil
}
