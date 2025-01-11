package main

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/elogiclab/shelly_em_exporter/internal/shelly"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/promslog/flag"
	"github.com/prometheus/common/version"
	"net/http"
	"os"
)

func main() {
	var (
		listenAddress = kingpin.Flag("listen-address", "Address on which to expose metrics and web interface.").Default(":9456").String()
		metricsPath   = kingpin.Flag("telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
		remoteHost    = kingpin.Flag("shelly-host", "Shelly EM IP Address").Required().IP()
	)
	promlogConfig := &promslog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("node_exporter"))
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promslog.New(promlogConfig)
	c := shelly.Collector{Logger: logger, Host: remoteHost.String()}
	prometheus.MustRegister(&c)
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Shelly EM exporter</title></head>
			<body>
			<h1>Shelly EM exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})
	level.Info(logger).Log("msg", "Starting shelly EM exporter")
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		//level.Error(logger).Log("msg", "Cannot start!")
	}
}
