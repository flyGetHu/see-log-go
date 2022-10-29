package router

import (
	"fmt"
	"net/http"
	"see-log-go/enitty"
	"see-log-go/service"
)

func seeLog(rw http.ResponseWriter, req *http.Request) {

	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "不支持Stream", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	broker := enitty.NewBroker()
	messageChan := make(chan []byte)

	broker.NewClients <- messageChan

	defer func() {
		broker.ClosingClients <- messageChan
	}()

	ctx := req.Context()

	go func() {
		<-ctx.Done()
		broker.ClosingClients <- messageChan
	}()
	go func() {
		service.MonitorFile(broker)
	}()
	for {
		bytes := <-messageChan
		_, err := fmt.Fprintf(rw, "%s", bytes)
		if err != nil {
			fmt.Println("输出http失败:", err)
		}
		flusher.Flush()
	}
}
