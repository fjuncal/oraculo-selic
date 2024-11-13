package routes

import (
	"github.com/rs/cors"
	"net/http"
	"oraculo-selic/controllers"
)

func SetupRoutes(messageController *controllers.MessageController) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/messages", messageController.CreateMessageHandler)
	mux.HandleFunc("/api/messages/list", messageController.GetMessagesHandler)
	mux.HandleFunc("/status", messageController.StatusHandler)

	// Adiciona suporte a CORS
	handler := cors.Default().Handler(mux)
	return handler
}
