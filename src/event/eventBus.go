package event

var eventChannels []chan interface{}

func Send(event interface{}) {
	go func() {
		for _, value := range eventChannels {
			value <- &event
		}
	}()
}

func RegisterChannel(channel chan interface{}) {
	eventChannels = append(eventChannels, channel)
}
