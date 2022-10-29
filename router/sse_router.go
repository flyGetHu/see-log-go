package router

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"see-log-go/enitty"
	"see-log-go/service"
	"strconv"
	"strings"
)

// 查看指定日志文件
func seeLog(rw http.ResponseWriter, req *http.Request) {
	values := req.URL.Query()
	logFilePath, logFileCount, err := dealSseParams(&values)
	if err != nil {
		_, _ = fmt.Fprint(rw, err.Error())
		return
	}
	if logFileCount == 0 {
		logFileCount = 100
	}
	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "不支持Stream", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/event-stream;charset:utf-8")
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
		service.MonitorFile(broker, logFilePath, logFileCount)
	}()
	for {
		bytes := <-messageChan
		_, err := fmt.Fprintf(rw, "%s\n", bytes)
		if err != nil {
			log.Println("输出http失败:", err)
		}
		flusher.Flush()
	}
}

// 处理请求参数
func dealSseParams(values *url.Values) (string, int64, error) {
	logFilePath := values.Get("file_path")
	logFileCount, err := strconv.Atoi(values.Get("count"))
	if err != nil {
		logFileCount = 100
	}
	stat, err := os.Stat(logFilePath)
	if err != nil || os.IsNotExist(err) || stat.IsDir() {
		return "", 0, fmt.Errorf("文件不存在或无效日志文件")

	}
	if !strings.HasSuffix(logFilePath, ".log") {
		return "", 0, fmt.Errorf("文件类型错误,支支持log类型文件")
	}
	return logFilePath, int64(logFileCount), nil
}
