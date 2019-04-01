package event

var eventChannels []chan interface{}

func Send(event interface{}) {
	for _, value := range eventChannels {
		go func(channel chan interface{}) {
			channel <- event
		}(value)
	}
}

func RegisterChannel(channel chan interface{}) {
	eventChannels = append(eventChannels, channel)
}
