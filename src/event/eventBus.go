package event

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

var subscribers []*subscriber

func Send(event interface{}) {
	for _, value := range subscribers {
		go func(subscriber *subscriber) {
			_, ok := subscriber.handlers[getEventKey(event)]
			if ok {
				subscriber.channel <- event
			}
		}(value)
	}
}

type Handlers map[string]func(interface{})

func Subscribe(cacheSize int, handlers Handlers) {
	i := make(chan interface{}, cacheSize)
	subscribers = append(subscribers, &subscriber{
		channel:  i,
		handlers: handlers,
	})
}

type subscriber struct {
	channel  chan interface{}
	handlers Handlers
}

func Listen() {
	for _, value := range subscribers {
		go func(subscriber *subscriber) {
			ok := true
			var evt interface{}
			for {
				select {
				case evt, ok = <-subscriber.channel:
					logrus.WithFields(logrus.Fields{
						"event": evt,
						"ok":    ok,
					}).Debug("index handler Event")
				}
				if !ok {
					break
				}
				handler, ok := subscriber.handlers[getEventKey(evt)]
				if ok {
					handler(evt)
				}
			}
		}(value)
	}
}

func getEventKey(evt interface{}) string {
	return fmt.Sprintf("%T", evt)
}
