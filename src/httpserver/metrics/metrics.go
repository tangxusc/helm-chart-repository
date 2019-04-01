package metrics

import (
	prometheusMiddleware "github.com/iris-contrib/middleware/prometheus"
	"github.com/kataras/iris"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"repository/event"
	"repository/httpserver"
	"repository/repository/index"
)

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

	event.Subscribe(100, event.Handlers{
		"*index.ChartUpdated": handlerEvent,
	})
}

func handlerEvent(event interface{}) {
	updated := event.(*index.ChartUpdated)
	logrus.WithField("chart count", updated).Debug()
	chartGauge.Set(float64(updated.ChartTotal))
}
