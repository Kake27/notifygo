package ws

import (
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"notification-service/service"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// In dev allow any origin. In prod restrict this.
	CheckOrigin: func(r *http.Request) bool { return true },
}

func StartServer(wsAddr string, pushService *service.PushService) {
	r := mux.NewRouter()
	r.HandleFunc("/ws/{userID}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userID"]
		if userID == "" {
			http.Error(w, "missing userID", http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("ws upgrade failed for %s: %v", userID, err)
			return
		}

		// register client
		pushService.RegisterClient(userID, conn)


		go func() {
			defer pushService.UnregisterClient(userID)

			for {
				_, _, err := conn.ReadMessage()
				if err != nil {

					log.Printf("ws read closed for %s: %v", userID, err)
					return
				}
			}
		}()
	}).Methods("GET")

	srv := &http.Server{
		Addr:         wsAddr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("WebSocket server listening on %s", wsAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ws server failed: %v", err)
		}
	}()
}