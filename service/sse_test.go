package service

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"
)

type Broker struct {
	Notifier chan []byte

	newClients chan chan []byte

	closingClients chan chan []byte

	clients map[chan []byte]bool
}

func (broker *Broker) listen() {
	for {
		select {
		case s := <-broker.newClients:
			//新连接进行注册
			broker.clients[s] = true
		case s := <-broker.closingClients:
			//断开删除链接
			delete(broker.clients, s)
		case event := <-broker.Notifier:
			//广播消息
			for clientMessageChan := range broker.clients {
				clientMessageChan <- event
			}
		}
	}
}

func (broker *Broker) ServeHTTP(rw http.ResponseWriter, resp *http.Request) {
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

	notify := rw.(http.CloseNotifier).CloseNotify()

	go func() {
		<-notify
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

func NewSSE() (broker *Broker) {
	broker = &Broker{
		Notifier:       make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}

	go broker.listen()
	return
}

//func  ServeHTTP(writer http.ResponseWriter, request *http.Request) {
//	//TODO implement me
//	panic("implement me")
//}

func TestSse(*testing.T) {
	broker := NewSSE()
	go func() {
		for {
			time.Sleep(time.Second * 2)
			data := fmt.Sprintf("===>%s", time.Now().Format("2006-01-02 15:01:05"))
			log.Println("Sending event data")
			broker.Notifier <- []byte(data)
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second * 1)
			if time.Now().Second()%2 == 0 {
				data := fmt.Sprintf("-->%s", time.Now().Format("2006-01-02 15:01:05"))
				log.Println("Sending event data2")
				broker.Notifier <- []byte(data)
			}
		}
	}()
	http.Handle("/sse", broker)

	log.Fatalln("error", http.ListenAndServe(":8080", nil))
}
