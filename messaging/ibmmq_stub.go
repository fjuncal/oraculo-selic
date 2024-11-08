package messaging

import "fmt"

// IBMMQClientStub é uma implementação que retorna um erro para ambientes sem IBM MQ
type IBMMQClient struct{}

// NewIBMMQClient retorna um erro informando que o IBM MQ não está disponível
func NewIBMMQClient(queueManager, queueName, connectionName, channel, userID, password string) (*IBMMQClient, error) {
	return nil, fmt.Errorf("IBM MQ não está disponível neste ambiente")
}

func (m *IBMMQClient) SendMessage(queueName, content string) error {
	return fmt.Errorf("IBM MQ não está disponível neste ambiente")
}

func (m *IBMMQClient) Close() error {
	return nil
}
