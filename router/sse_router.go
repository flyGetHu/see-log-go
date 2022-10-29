package router

import (
	"fmt"
	"net/http"
	"os"
	"see-log-go/enitty"
	"see-log-go/service"
	"strconv"
	"strings"
)

// 查看指定日志文件
func seeLog(rw http.ResponseWriter, req *http.Request) {
	values := req.URL.Query()
	logFilePath := values.Get("path")
	logFileCount, err := strconv.Atoi(values.Get("count"))
	if err != nil {
		logFileCount = 100
	}
	stat, err := os.Stat(logFilePath)
	if err != nil {
		_, _ = fmt.Fprint(rw, "文件不存在")
		return
	}
	if os.IsNotExist(err) {
		_, _ = fmt.Fprint(rw, "文件不存在")
		return
	}
	if stat.IsDir() {
		_, _ = fmt.Fprint(rw, "文件错误")
		return
	}
	if !strings.HasSuffix(logFilePath, ".log") {
		_, _ = fmt.Fprint(rw, "文件类型错误,支支持log类型文件")
		return
	}
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
		service.MonitorFile(broker, logFilePath, int64(logFileCount))
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
