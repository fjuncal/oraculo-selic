package routes

import (
	"github.com/rs/cors"
	"net/http"
	"oraculo-selic/controllers"
)

func SetupRoutes(messageController *controllers.MessageController, cenarioController *controllers.CenarioController) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/messages", messageController.CreateMessageHandler)
	mux.HandleFunc("/api/messages/list", messageController.GetMessagesHandler)
	mux.HandleFunc("/status", messageController.StatusHandler)

	mux.HandleFunc("/api/cenarios", cenarioController.SaveCenarioHandler)
	mux.HandleFunc("/api/cenarios/list", cenarioController.GetCenariosHandler) // Nova rota para buscar cenários

	// Adiciona suporte a CORS
	handler := cors.Default().Handler(mux)
	return handler
}
