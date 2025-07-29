package handler

import (
	"notification-service/service"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"
)

type NotificationHandler struct {
	emailService *service.EmailService
	smsService *service.SMSService
	pushService *service.PushService
}

type PushSocket struct {

}

func NewNotificationHandler(e *service.EmailService, s *service.SMSService, p *service.PushService) *NotificationHandler {
	return &NotificationHandler{
		emailService: e,
		smsService: s,
		pushService: p,
	}
}

type NotificationRequest struct {
	Type    string `json:"type"`    // type of notification
	To      string `json:"to"`
	Subject string `json:"subject"` 
	Message string `json:"message"`
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
    if req.Message == "" {
        missing = append(missing, "message")
    }
	if req.Type == "" {
		missing = append(missing, "type")
	}

    if len(missing) > 0 {
        return nil, http.ErrorMissingParam{Params: missing}
    }

	// send the notification
	var err error
	switch req.Type {
	case "email":
		err = h.emailService.Send(req.To, req.Subject, req.Message)

	case "sms":        
		_, err = h.smsService.Send(req.To, req.Message)

	case "push":
		err = h.pushService.Send(req.To, req.Subject, req.Message)
	default:
		return nil, http.ErrorInvalidParam{Params: []string{"type"}}
	}

	if err != nil {
		return nil, err
	}

	return map[string]string{"status": "notification sent succesfully"}, nil
}

func (h *NotificationHandler) PushSocket(ctx *gofr.Context) (any, error) {
	userID := ctx.PathParam("userID")
	ctx.Logger.Infof("Registered WebSocket client for clientID: %s", userID)
	// ctx.Logger.Infof("🔥 Raw URL: %s", ctx.Request)
	// ctx.Logger.Infof("🔥 Extracted userID: %s", userID)

	h.pushService.RegisterClient(userID, ctx)
	for {
		var msg string
		if err := ctx.Bind(&msg); err != nil {
			break
		}
		ctx.Logger.Infof("Received from %s: %s", userID, msg)
	}

	return nil, nil

}