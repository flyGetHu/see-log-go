package router

import (
	"see-log-go/enitty"
)

func newSSE() (broker *enitty.Broker) {
	broker = enitty.NewBroker()
	go broker.Listen()
	return
}
