package service

import (
	"fmt"
	"github.com/hpcloud/tail"
	"testing"
	"time"
)

func TestFileTail(t *testing.T) {
	fileName := "D:\\home\\work\\admin-app\\admin-app_2022-07-11.log"
	config := tail.Config{
		ReOpen:    true,                                   // 重新打开
		Follow:    true,                                   // 是否跟随
		Location:  &tail.SeekInfo{Offset: 100, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                  // 文件不存在不报错
		Poll:      true,
	}
	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		fmt.Println("tail file failed, err:", err)
		return
	}
	for {
		line, ok := <-tails.Lines
		if !ok {
			fmt.Printf("tail file close reopen, filename:%s\n", tails.Filename)
			time.Sleep(time.Second)
			continue
		}
		fmt.Println("line:", line.Text)
	}
}