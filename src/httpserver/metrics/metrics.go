package metrics

import (
	"fmt"
	prometheusMiddleware "github.com/iris-contrib/middleware/prometheus"
	"github.com/kataras/iris"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"repository/event"
	"repository/httpserver"
)

var eventChan = make(chan interface{}, 1000)
var counter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "chart_total",
	Help: "chart 总数",
})

/**
注册prometheus
 */
func init() {
	promInst := prometheusMiddleware.New("helm-chart-repository")
	httpserver.AddRegister(func(app *iris.Application) {
		app.Use(promInst.ServeHTTP)
		app.Any("/metrics", iris.FromStd(promhttp.Handler()))
	})
	prometheus.MustRegister(counter)
	event.RegisterChannel(eventChan)
}

func Listen() {
	for {
		ok := true
		var event interface{}
		select {
		case event, ok = <-eventChan:
			fmt.Println(event, ok)
		}
		if !ok {
			break
		}
		handlerEvent(event)
	}
}

func handlerEvent(event interface{}) {
	//TODO:处理事件
	fmt.Println(event)
}
