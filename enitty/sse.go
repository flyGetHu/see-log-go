package enitty

import (
	"fmt"
)

type Broker struct {
	Notifier chan []byte

	NewClients chan chan []byte

	ClosingClients chan chan []byte

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
		case s := <-broker.NewClients:
			//新连接进行注册
			fmt.Println("新用户注册")
			broker.clients[s] = true
		case s := <-broker.ClosingClients:
			//断开删除链接
			fmt.Println("用户断开")
			delete(broker.clients, s)
		case event := <-broker.Notifier:
			//广播消息
			for clientMessageChan := range broker.clients {
				clientMessageChan <- event
			}
		}
	}
}
