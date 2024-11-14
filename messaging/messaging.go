package messaging

import (
	"github.com/go-stomp/stomp"
	"log"
)

type Messaging interface {
	SendMessage(queueName, content string) error
	Close() error
}

type ActiveMQClient struct {
	conn     *stomp.Conn
	queueURL string
}

func NewActiveMQClient(queueURL string) (*ActiveMQClient, error) {
	client := &ActiveMQClient{queueURL: queueURL}
	if err := client.connect(); err != nil {
		return nil, err
	}
	return client, nil
}

// SendMessage Função para enviar mensagem com reconexão automática
func (c *ActiveMQClient) SendMessage(queueName, content string) error {
	// Tenta enviar a mensagem
	err := c.conn.Send(queueName, "text/plain", []byte(content))
	if err != nil {
		// Trata o erro de conexão fechada e tenta reconectar
		log.Printf("Erro ao enviar mensagem para a fila %s: %v", queueName, err)
		log.Println("Tentando reconectar...")
		if err := c.connect(); err != nil {
			return err
		}

		// Após reconectar, tenta enviar a mensagem novamente
		err = c.conn.Send(queueName, "text/plain", []byte(content))
		if err != nil {
			log.Printf("Erro ao enviar mensagem após reconexão para a fila %s: %v", queueName, err)
			return err
		}
	}
	log.Println("Mensagem enviada com sucesso para a fila:", queueName)
	return nil
}

// Close Função para fechar a conexão
func (c *ActiveMQClient) Close() error {
	if c.conn != nil {
		return c.conn.Disconnect()
	}
	return nil
}

// connect Função para estabelecer ou restabelecer a conexão
func (c *ActiveMQClient) connect() error {
	conn, err := stomp.Dial("tcp", c.queueURL)
	if err != nil {
		log.Printf("Erro ao conectar ao ActiveMQ: %v", err)
		return err
	}
	c.conn = conn
	log.Println("Conexão estabelecida com sucesso com o ActiveMQ")
	return nil
}
