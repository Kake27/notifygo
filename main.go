package main

import (
	"gofr.dev/pkg/gofr"
	"notification-service/handler"
	"notification-service/service"
)


func main() {
	app := gofr.New()

	app.GET("/greet", func(ctx *gofr.Context) (any, error) {
		return "hello world", nil
	})

	// get health
	healthHandler := handler.NewHealthHandler()
	app.GET("/health", healthHandler.Health)


	//email
	emailSvc := service.NewEmailService()

	//sms
	smsSvc := service.NewSMSService(	) 


	//push
	pushSvc := service.NewPushService()

	//notification handler
	notificationHandler := handler.NewNotificationHandler(emailSvc, smsSvc, pushSvc)

	app.WebSocket("/ws/{userID}", notificationHandler.PushSocket)	
	app.POST("/notify", notificationHandler.Notify)

	

	app.Run()
}