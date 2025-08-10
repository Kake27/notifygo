package main

import (
	"log"
	"notification-service/handler"
	"notification-service/service"
	"notification-service/store"
	"notification-service/ws"
	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	// basic route to check app is running
	app.GET("/greet", func(ctx *gofr.Context) (any, error) {
		return "hello world", nil
	})

	// get health
	healthHandler := handler.NewHealthHandler()
	app.GET("/health", healthHandler.Health)

	// connect to db
	db, err := store.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to db: ", err)
	}

	store.CreateDB(db)

	// db and file template store
	dbTemplateStore := store.NewDBTemplateStore(db)
	fileTemplateStore, _ := store.NewTemplateStore()

	// render templates
	templateRenderer := service.NewTemplateRenderer(fileTemplateStore, dbTemplateStore)

	//email
	emailSvc := service.NewEmailService()

	//sms
	smsSvc := service.NewSMSService()

	// push
	pushSvc := service.NewPushService()

	// WebSocket server for push notifications
	ws.StartServer(":8081", pushSvc)


	//notification handler
	notificationHandler := handler.NewNotificationHandler(emailSvc, smsSvc, pushSvc, templateRenderer, dbTemplateStore)

	// notify route
	app.POST("/notify", notificationHandler.Notify)

	// create and delete templated messages
	app.POST("/template/create", notificationHandler.CreateTemplate)
	app.DELETE("/template/delete/{name}", notificationHandler.DeleteTemplate)



	app.Run()
}
