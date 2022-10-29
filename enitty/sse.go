package enitty

import (
	"log"
)

type Broker struct {
	Notifier chan []byte

	NewClients chan chan []byte

	ClosingClients chan chan []byte

	//客户端建立链接 链接通道->是否停用
	clients map[chan []byte]bool
}

func NewBroker() *Broker {
	broker := &Broker{
		Notifier:       make(chan []byte, 1),
		NewClients:     make(chan chan []byte),
		ClosingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}
	go broker.listen()
	return broker
}

func (broker *Broker) HasClient() bool {
	return len(broker.clients) != 0
}

func (broker *Broker) listen() {
	for {
		select {
		case socketChan := <-broker.NewClients:
			//新连接进行注册
			log.Println("新用户注册", &socketChan)
			broker.clients[socketChan] = false
		case socketChan := <-broker.ClosingClients:
			//断开删除链接
			log.Println("用户断开", &socketChan)
			delete(broker.clients, socketChan)
		case event := <-broker.Notifier:
			//广播消息
			for clientMessageChan, isDeactivate := range broker.clients {
				if isDeactivate {
					clientMessageChan <- event
				}
			}
		}
	}
}
