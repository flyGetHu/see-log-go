package router

import (
	"fmt"
	"log"
	"see-log-go/enitty"
	"time"
)

func broker() *enitty.Broker {
	broker := enitty.NewBroker()
	go func() {
		for {
			time.Sleep(time.Second * 2)
			data := fmt.Sprintf("===>%s", time.Now().Format("2006-01-02 15:01:05"))
			log.Println("Sending event data")
			broker.Notifier <- []byte(data)
		}
	}()
	return broker
}
