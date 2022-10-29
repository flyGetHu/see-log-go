package enitty

import (
	"fmt"
	"net/http"
)

type Broker struct {
	Notifier chan []byte

	newClients chan chan []byte

	closingClients chan chan []byte

	clients map[chan []byte]bool
}

func NewBroker() *Broker {
	broker := &Broker{
		Notifier:       make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}
	go broker.listen()
	return broker
}

func (broker *Broker) listen() {
	for {
		select {
		case s := <-broker.newClients:
			//新连接进行注册
			fmt.Println("新用户注册")
			broker.clients[s] = true
		case s := <-broker.closingClients:
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

func (broker *Broker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "不支持Stream", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	messageChan := make(chan []byte)

	broker.newClients <- messageChan

	defer func() {
		broker.closingClients <- messageChan
	}()

	ctx := req.Context()

	go func() {
		<-ctx.Done()
		broker.closingClients <- messageChan
	}()

	for {
		_, err := fmt.Fprintf(rw, "data: %s\n\n", <-messageChan)
		if err != nil {
			fmt.Println("输出http失败:", err)
		}
		flusher.Flush()
	}
}
