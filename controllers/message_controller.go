package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"oraculo-selic/api"
	"oraculo-selic/db"
	"oraculo-selic/messaging"
	"runtime/debug"
)

type MessageController struct {
	Api *api.Api
}

func NewMessageController(dbConn *db.DatabaseConnections, msgService messaging.Messaging) *MessageController {
	return &MessageController{
		Api: api.NewApi(dbConn, msgService),
	}
}

func (mc *MessageController) CreateMessageHandler(w http.ResponseWriter, r *http.Request) {
	mc.Api.CreateMessageHandler(w, r)
}

func (mc *MessageController) GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	mc.Api.GetMessagesHandler(w, r)
}

func (mc *MessageController) StatusHandler(w http.ResponseWriter, r *http.Request) {
	correlationId := r.URL.Query().Get("correlationId")
	if correlationId == "" {
		log.Printf("Erro: Message ID é obrigatório\nStack Trace:\n%s", debug.Stack())
		http.Error(w, "Message ID é obrigatório", http.StatusBadRequest)
		return
	}

	sentStatus, arrivedStatus, processedStatus, err := mc.Api.CheckStatus(correlationId)
	if err != nil {
		http.Error(w, "Erro ao verificar status", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	response := map[string]interface{}{
		"sent":      api.APIResponse{Status: sentStatus, Detail: "Status de envio"},
		"arrived":   api.APIResponse{Status: arrivedStatus, Detail: "Status de chegada"},
		"processed": api.APIResponse{Status: processedStatus, Detail: "Status de processamento"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
