package api

import (
	"encoding/json"
	"log"
	"net/http"
	"oraculo-selic/db"
	"oraculo-selic/messaging"
	"oraculo-selic/models"
)

type API struct {
	dbConnections *db.DatabaseConnections
	messaging     messaging.Messaging
}

// NewApi criando nova instancia de API
func NewApi(dbConnections *db.DatabaseConnections, messaging messaging.Messaging) *API {
	return &API{dbConnections: dbConnections, messaging: messaging}
}

// CreateMessageHandler Handler para criar uma nova mensagem
func (api *API) CreateMessageHandler(w http.ResponseWriter, r *http.Request) {
	var message models.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Entrada inválida", http.StatusBadRequest)
		return
	}

	// Define o status inicial e salva no banco usando DB1
	message.Status = "RECEIVED"
	log.Println("Salvando mensagem no banco de dados...")
	if err := api.dbConnections.SaveMessage(&message); err != nil {
		log.Printf("Erro ao salvar mensagem no banco: %v", err)
		http.Error(w, "Erro ao salvar a mensagem", http.StatusInternalServerError)
		return
	}
	log.Println("Mensagem salva com sucesso no banco de dados.")

	log.Println("Enviando mensagem para a fila...")
	if err := api.messaging.SendMessage("queue.RECEIVE_QUEUE", message.Content); err != nil {
		log.Printf("Erro ao enviar mensagem para a fila: %v", err)
		http.Error(w, "Erro ao enviar a mensagem para a fila", http.StatusInternalServerError)
		return
	}
	log.Println("Mensagem enviada com sucesso para a fila.")

	// Retorna a mensagem como resposta com status 201 Created
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func (api *API) CheckStatus(messageID string) (string, string, string, error) {
	var sentStatus, arrivedStatus, processedStatus string

	// Usando DB1 para verificar o status de envio
	err := api.dbConnections.DB1.QueryRow("SELECT status FROM messages WHERE id = $1", messageID).Scan(&sentStatus)
	if err != nil {
		return "", "", "", err
	}

	// Usando DB2 para verificar o status de chegada
	err = api.dbConnections.DB2.QueryRow("SELECT status FROM message_entity WHERE id = $1", messageID).Scan(&processedStatus)
	if err != nil {
		return "", "", "", err
	}

	// Usando DB3 para verificar o status de processamento
	err = api.dbConnections.DB3.QueryRow("SELECT status FROM message_entity WHERE id = $1", messageID).Scan(&arrivedStatus)
	if err != nil {
		return "", "", "", err
	}

	return sentStatus, arrivedStatus, processedStatus, nil
}
func (api *API) GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := api.dbConnections.DB1.Query("SELECT id, content, status FROM messages")
	if err != nil {
		http.Error(w, "Erro ao buscar mensagens", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.ID, &message.Content, &message.Status); err != nil {
			http.Error(w, "Erro ao ler mensagens", http.StatusInternalServerError)
			return
		}
		messages = append(messages, message)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
