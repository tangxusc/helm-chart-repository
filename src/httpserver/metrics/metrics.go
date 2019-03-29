package metrics

import (
	prometheusMiddleware "github.com/iris-contrib/middleware/prometheus"
	"github.com/kataras/iris"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"repository/event"
	"repository/httpserver"
	"repository/repository"
)

var eventChan = make(chan interface{}, 1000)
var chartGauge = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "chart",
	Help: "chart 数量",
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
	prometheus.MustRegister(chartGauge)
	event.RegisterChannel(eventChan)
}

func Listen() {
	for {
		ok := true
		var evt interface{}
		select {
		case evt, ok = <-eventChan:
			logrus.WithFields(logrus.Fields{
				"event": evt,
				"ok":    ok,
			}).Debug("handler Event")
		}
		if !ok {
			break
		}
		handlerEvent(evt)
	}
}

func handlerEvent(event interface{}) {
	switch event.(type) {
	case *repository.ChartUpdated:
		updated := event.(*repository.ChartUpdated)
		logrus.WithField("chart count", updated).Debug()
		chartGauge.Set(float64(updated.ChartTotal))
	}
}
