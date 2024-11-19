package main

import (
	"fmt"
	"log"
	"net/http"
	"oraculo-selic/config"
	"oraculo-selic/controllers"
	"oraculo-selic/db"
	"oraculo-selic/db/repositories"
	"oraculo-selic/messaging"
	"oraculo-selic/routes"
	"os"
)

func main() {
	// Carregar as variáveis de ambiente e exibir
	fmt.Println("DATABASE_URL_1:", os.Getenv("DATABASE_URL_1"))
	fmt.Println("DATABASE_URL_2:", os.Getenv("DATABASE_URL_2"))
	fmt.Println("DATABASE_URL_3:", os.Getenv("DATABASE_URL_3"))
	fmt.Println("QUEUE_URL:", os.Getenv("QUEUE_URL"))
	fmt.Println("MESSAGING_TYPE:", os.Getenv("MESSAGING_TYPE"))

	// Carregar configurações
	cfg := config.LoadConfig()

	// Configurar serviço de mensageria
	log.Println("Configurando o serviço de mensageria...")
	var msgService messaging.Messaging
	var err error

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
	defer msgService.Close()
	log.Println("Conexão com o serviço de mensageria estabelecida com sucesso.")

	// Conectar aos bancos de dados
	log.Println("Conectando aos bancos de dados...")
	dbConn, err := db.NewDatabaseConnections(os.Getenv("DATABASE_URL_1"), os.Getenv("DATABASE_URL_2"), os.Getenv("DATABASE_URL_3"))
	if err != nil {
		log.Fatalf("Erro ao conectar aos bancos de dados: %v", err)
	}
	defer dbConn.Close()
	log.Println("Conexões com os bancos de dados estabelecidas com sucesso.")

	messageController := controllers.NewMessageController(dbConn, msgService)

	// Aqui criamos uma instância de db.DB a partir de dbConn.DB1 e passamos para o PassoTesteController
	dbInstance := &db.DB{Conn: dbConn.DB1}
	passoTesteController := controllers.NewPassoTesteController(dbInstance)

	cenarioRepository := repositories.NewCenarioRepository(dbConn.DB1)
	cenarioController := controllers.NewCenarioController(cenarioRepository)

	handler := routes.SetupRoutes(messageController, passoTesteController, cenarioController)
	log.Println("Servidor iniciado na porta 8086")
	log.Fatal(http.ListenAndServe(":8086", handler))
}
