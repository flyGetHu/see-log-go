package service

import (
	"fmt"
	"github.com/hpcloud/tail"
	"io"
	"log"
	"see-log-go/enitty"
	"time"
)

func MonitorFile(broker *enitty.Broker, path string, count int64) {
	config := tail.Config{
		ReOpen:    true,                                                  // 重新打开
		Follow:    true,                                                  // 是否跟随
		Location:  &tail.SeekInfo{Offset: count, Whence: io.SeekCurrent}, // 从文件的哪个地方开始读
		MustExist: false,                                                 // 文件不存在不报错
		Poll:      true,
	}
	tails, err := tail.TailFile(path, config)
	if err != nil {
		log.Println("tail file failed, err:", err)

	}
	for {
		line, ok := <-tails.Lines
		if broker.HasClient() {
			if !ok {
				fmt.Printf("tail file close reopen, filename:%s\n", tails.Filename)
				time.Sleep(time.Second)
				continue
			}
			text := line.Text
			broker.Notifier <- []byte(text)
		}
	}
}
