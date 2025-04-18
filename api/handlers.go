package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"oraculo-selic/db"
	"oraculo-selic/messaging"
	"oraculo-selic/models"
	"time"
)

type Api struct {
	dbConnections *db.DatabaseConnections
	messaging     messaging.Messaging
}

// NewApi criando nova instancia de Api
func NewApi(dbConnections *db.DatabaseConnections, messaging messaging.Messaging) *Api {
	return &Api{dbConnections: dbConnections, messaging: messaging}
}

// CreateMessageHandler Handler para criar e processar uma lista de mensagens
func (api *Api) CreateMessageHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Descricao    string              `json:"descricao"`
		Tipo         string              `json:"tipo"`
		PassosTestes []models.PassoTeste `json:"passosTestes"`
	}

	// Decodifica o JSON recebido
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Erro ao decodificar request: %v", err)
		http.Error(w, "Entrada inválida", http.StatusBadRequest)
		return
	}

	// Verifica se existem passos testes no cenário
	if len(request.PassosTestes) == 0 {
		log.Printf("Nenhum passo teste fornecido no cenário")
		http.Error(w, "Cenário sem passos testes", http.StatusBadRequest)
		return
	}

	// Processa cada passo teste no cenário
	for _, passoTeste := range request.PassosTestes {
		message := models.Mensagem{
			CodigoMensagem: passoTeste.CodigoMsg,
			Canal:          passoTeste.Canal,
			XML:            passoTeste.MsgDocXML,
			StringSelic:    passoTeste.Msg,
			Status:         "ENVIANDO",
			DataInclusao:   NowInBrazil(),
		}

		// Salva a mensagem no banco de dados
		if err := api.dbConnections.SaveMessage(&message); err != nil {
			log.Printf("Erro ao salvar mensagem no banco: %v", err)
			http.Error(w, "Erro ao salvar a mensagem", http.StatusInternalServerError)
			return
		}

		// Serializa e envia para a fila
		messageJSON, err := json.Marshal(message)
		if err != nil {
			log.Printf("Erro ao serializar a mensagem para JSON: %v", err)
			http.Error(w, "Erro ao preparar a mensagem para a fila", http.StatusInternalServerError)
			return
		}

		if err := api.messaging.SendMessage("queue.RECEIVE_QUEUE", string(messageJSON)); err != nil {
			log.Printf("Erro ao enviar mensagem para a fila: %v", err)
			http.Error(w, "Erro ao enviar a mensagem para a fila", http.StatusInternalServerError)
			return
		}
	}

	// Responde com sucesso
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Cenário enviado com sucesso"}`))
}

func (api *Api) CheckStatus(correlationId string) (string, string, string, error) {
	var sentStatus, arrivedStatus, processedStatus string

	err := api.dbConnections.DB1.QueryRow("SELECT txt_status FROM mensagens WHERE txt_correl_id = $1", correlationId).Scan(&sentStatus)
	if errors.Is(err, sql.ErrNoRows) {
		sentStatus = "NÃO PROCESSADO"
	} else if err != nil {
		return "", "", "", err
	}

	err = api.dbConnections.DB2.QueryRow("SELECT txt_status FROM mensagens WHERE txt_correl_id = $1", correlationId).Scan(&processedStatus)
	if errors.Is(err, sql.ErrNoRows) {
		processedStatus = "NÃO PROCESSADO"
	} else if err != nil {
		return "", "", "", err
	}

	// Usando DB3 para verificar o status de processamento
	err = api.dbConnections.DB3.QueryRow("SELECT txt_status FROM mensagens WHERE txt_correl_id = $1", correlationId).Scan(&arrivedStatus)
	if errors.Is(err, sql.ErrNoRows) {
		arrivedStatus = "NÃO PROCESSADO"
	} else if err != nil {
		return "", "", "", err
	}

	return sentStatus, arrivedStatus, processedStatus, nil
}
func (api *Api) GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := api.dbConnections.DB1.Query(`
    SELECT id, txt_cod_msg, txt_canal, txt_msg_doc_xml, txt_msg, txt_status, dt_incl, txt_correl_id
    FROM mensagens
`)
	if err != nil {
		http.Error(w, "Erro ao buscar mensagens", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []map[string]interface{}
	for rows.Next() {
		var message models.Mensagem
		if err := rows.Scan(
			&message.ID,
			&message.CodigoMensagem,
			&message.Canal,
			&message.XML,
			&message.StringSelic,
			&message.Status,
			&message.DataInclusao,
			&message.CorrelationID,
		); err != nil {
			http.Error(w, "Erro ao ler mensagens", http.StatusInternalServerError)
			return
		}

		// Consultar o status final na base SELIC_OPE_POC (DB2)
		var finalStatus string
		err = api.dbConnections.DB2.QueryRow(`
			SELECT txt_status FROM mensagens WHERE txt_correl_id  = $1
		`, message.CorrelationID).Scan(&finalStatus)

		// Se o status não for encontrado, define um valor padrão
		if err == sql.ErrNoRows {
			finalStatus = "NÃO PROCESSADO" // Valor padrão para mensagens sem status final
		} else if err != nil {
			fmt.Println("Erro ao buscar status final na base SELIC_OPE_POC", err)
			http.Error(w, "Erro ao buscar status final na base SELIC_OPE_POC", http.StatusInternalServerError)
			return
		}
		messageMap := map[string]interface{}{
			"id":             message.ID,
			"codigoMensagem": message.CodigoMensagem,
			"canal":          message.Canal,
			"xml":            message.XML,
			"stringSelic":    message.StringSelic,
			"status":         message.Status,
			"statusFinal":    finalStatus, // Status final obtido da SELIC_OPE_POC
			"dataInclusao":   message.DataInclusao,
			"correlationId":  message.CorrelationID,
		}
		messages = append(messages, messageMap)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func NowInBrazil() string {
	location, _ := time.LoadLocation("America/Sao_Paulo")
	return time.Now().In(location).Format("2006-01-02T15:04:05")
}
