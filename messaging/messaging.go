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
	conn *stomp.Conn
}

func NewActiveMQClient(queueURL string) (*ActiveMQClient, error) {
	conn, err := stomp.Dial("tcp", queueURL)
	if err != nil {
		return nil, err
	}
	return &ActiveMQClient{conn: conn}, nil
}

func (c *ActiveMQClient) SendMessage(queueName, content string) error {
	// Define o tipo de conteúdo como "text/plain"
	err := c.conn.Send(
		queueName,
		"text/plain",
		[]byte(content),
	)
	if err != nil {
		log.Printf("Erro ao enviar mensagem para a fila %s: %v", queueName, err)
		return err
	}
	log.Println("Mensagem enviada com sucesso para a fila:", queueName)
	return nil
}

func (c *ActiveMQClient) Close() error {
	return c.conn.Disconnect()
}
