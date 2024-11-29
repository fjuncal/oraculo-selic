package routes

import (
	"github.com/rs/cors"
	"net/http"
	"oraculo-selic/controllers"
)

func SetupRoutes(messageController *controllers.MessageController, passoTesteController *controllers.PassoTesteController, cenarioController *controllers.CenarioController) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/messages", messageController.CreateMessageHandler)
	mux.HandleFunc("/api/messages/list", messageController.GetMessagesHandler)
	mux.HandleFunc("/status", messageController.StatusHandler)

	mux.HandleFunc("/api/passo-teste", passoTesteController.SavePassoTesteHandler)
	mux.HandleFunc("/api/passo-teste/list", passoTesteController.GetPassoTesteHandler)

	// Rotas de cenários
	mux.HandleFunc("/api/cenarios/save", cenarioController.SaveCenarioHandler)
	mux.HandleFunc("/api/cenarios/relacionar", cenarioController.SaveRelacionamentoHandler)
	mux.HandleFunc("/api/cenarios/list", cenarioController.GetCenariosHandler)
	mux.HandleFunc("/api/cenarios/upload", cenarioController.UploadPlanilhaHandler)

	//mux.HandleFunc("/api/cenarios/passo-teste", cenarioController.GetCenariosWithPassosTestesHandler)

	// Adiciona suporte a CORS
	handler := cors.Default().Handler(mux)
	return handler
}
