package main

import (
	"gofr.dev/pkg/gofr"
	"notification-service/handler"
	"notification-service/service"
)


func main() {
	app := gofr.New()

	app.GET("/greet", func(ctx *gofr.Context) (any, error) {
		return "Testing server!", nil
	})

	healthHandler := handler.NewHealthHandler()
	app.GET("/health", healthHandler.Health)


	emailSvc := service.NewEmailService()
	smsSvc := service.NewSMSService(
		"http://localhost:13013/cgi-bin/sendsms",
		"test",         // username
		"test123",      // password
		"KANNEL",       // sender
	) 
	pushSvc := service.NewPushService()

	notificationHandler := handler.NewNotificationHandler(emailSvc, smsSvc, pushSvc)
	app.POST("/notify", notificationHandler.Notify)

	

	app.Run()
}