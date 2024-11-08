//go:build ibmmq
// +build ibmmq

package messaging

import (
	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
	"log"
)

type IBMMQClient struct {
	qMgr  *ibmmq.MQQueueManager
	queue *ibmmq.MQObject
}

func NewIBMMQClient(queueManager, queueName, connectionName, channel, userID, password string) (*IBMMQClient, error) {
	cno := ibmmq.NewMQCNO()
	cno.Options = ibmmq.MQCNO_CLIENT_BINDING
	cd := ibmmq.NewMQCD()
	cd.ChannelName = channel
	cd.ConnectionName = connectionName
	cno.ClientConn = cd

	csp := ibmmq.NewMQCSP()
	csp.UserId = userID
	csp.Password = password
	cno.SecurityParms = csp

	qMgr, err := ibmmq.Connx(queueManager, cno)
	if err != nil {
		return nil, err
	}

	mqod := ibmmq.NewMQOD()
	mqod.ObjectName = queueName
	openOptions := ibmmq.MQOO_OUTPUT
	queue, err := qMgr.Open(mqod, openOptions)
	if err != nil {
		qMgr.Disc()
		return nil, err
	}

	log.Println("Conectado ao IBM MQ com sucesso.")
	return &IBMMQClient{qMgr: &qMgr, queue: &queue}, nil
}

func (m *IBMMQClient) SendMessage(queueName, content string) error {
	msg := ibmmq.NewMQMD()
	pmo := ibmmq.NewMQPMO()
	pmo.Options = ibmmq.MQPMO_NO_SYNCPOINT
	data := []byte(content)
	return m.queue.Put(msg, pmo, data)
}

func (m *IBMMQClient) Close() error {
	if m.queue != nil {
		m.queue.Close(0)
	}
	if m.qMgr != nil {
		m.qMgr.Disc()
	}
	return nil
}
