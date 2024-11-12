package main

import (
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"log"
	"net/http"
	"oraculo-selic/api"
	"oraculo-selic/config"
	"oraculo-selic/db"
	"oraculo-selic/messaging"
	"os"
	"runtime/debug"
)

type MessageStatus struct {
	Sent      *api.APIResponse `json:"sent"`
	Arrived   *api.APIResponse `json:"arrived"`
	Processed *api.APIResponse `json:"processed"`
}

func main() {
	fmt.Println("DATABASE_URL_1:", os.Getenv("DATABASE_URL_1"))
	fmt.Println("DATABASE_URL_2:", os.Getenv("DATABASE_URL_2"))
	fmt.Println("DATABASE_URL_3:", os.Getenv("DATABASE_URL_3"))
	fmt.Println("QUEUE_URL:", os.Getenv("QUEUE_URL"))
	fmt.Println("MESSAGING_TYPE:", os.Getenv("MESSAGING_TYPE"))

	// Carregar configurações
	cfg := config.LoadConfig()
	var msgService messaging.Messaging
	var err error

	// Configurar o serviço de mensageria com logs detalhados
	log.Println("Configurando o serviço de mensageria...")
	switch os.Getenv("MESSAGING_TYPE") {
	case "activemq":
		msgService, err = messaging.NewActiveMQClient(cfg.QueueURL)
	case "ibmmq":
		msgService, err = messaging.NewIBMMQClient(cfg.QueueURL, cfg.QueueName, cfg.ConnectionName, cfg.Channel, cfg.UserID, cfg.Password)
	default:
		log.Fatal("Tipo de mensageria não suportado. Use 'activemq' ou 'ibmmq'.")
	}
	if err != nil {
		log.Fatalf("Erro ao conectar ao serviço de mensageria: %v", err)
	}
	defer func() {
		log.Println("Fechando a conexão com o serviço de mensageria...")
		msgService.Close()
	}()
	log.Println("Conexão com o serviço de mensageria estabelecida com sucesso.")

	// Conectar às três bases de dados com logs detalhados
	log.Println("Conectando aos bancos de dados...")
	dbConn, err := db.NewDatabaseConnections(os.Getenv("DATABASE_URL_1"), os.Getenv("DATABASE_URL_2"), os.Getenv("DATABASE_URL_3"))
	if err != nil {
		log.Fatalf("Erro ao conectar aos bancos de dados: %v", err)
	}
	defer func() {
		log.Println("Fechando as conexões com os bancos de dados...")
		dbConn.Close()
	}()
	log.Println("Conexões com os bancos de dados estabelecidas com sucesso.")

	// Configurar e iniciar a API
	newApi := api.NewApi(dbConn, msgService)
	mux := http.NewServeMux() // Usando um multiplexer para rotas
	mux.HandleFunc("/api/messages", newApi.CreateMessageHandler)
	mux.HandleFunc("/api/messages/list", newApi.GetMessagesHandler) // Novo endpoint de listagem
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		correlationId := r.URL.Query().Get("correlationId")
		if correlationId == "" {
			log.Printf("Erro: Message ID é obrigatório\nStack Trace:\n%s", debug.Stack())
			http.Error(w, "Message ID é obrigatório", http.StatusBadRequest)
			return
		}

		// Usar a função CheckStatus para obter os status de envio, chegada e processamento
		sentStatus, arrivedStatus, processedStatus, err := newApi.CheckStatus(correlationId)
		if err != nil {
			http.Error(w, "Erro ao verificar status", http.StatusInternalServerError)
			log.Print(err)
			return
		}

		response := MessageStatus{
			Sent:      &api.APIResponse{Status: sentStatus, Detail: "Status de envio"},
			Arrived:   &api.APIResponse{Status: arrivedStatus, Detail: "Status de chegada"},
			Processed: &api.APIResponse{Status: processedStatus, Detail: "Status de processamento"},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Adicionar o middleware CORS às rotas
	handler := cors.Default().Handler(mux)

	log.Println("Servidor iniciado na porta 8086")
	log.Fatal(http.ListenAndServe(":8086", handler))
}
