package shelly

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type Collector struct {
	Logger log.Logger
	Host   string
}

var (
	scrapeDurationDesc = prometheus.NewDesc(
		"shellyem_scrape_duration_seconds",
		"shellyem_exporter: Duration of scraping shellyem.",
		nil,
		nil,
	)

	scrapeSuccessDesc = prometheus.NewDesc(
		"shellyem_scrape_success",
		"shellyem_exporter: Whether scraping shellyem succeeded.",
		nil,
		nil,
	)
	instantPowerDesc = prometheus.NewDesc(
		"shellyem_emeter_power",
		"shellyem_exporter: Measured instant power.",
		[]string{"host", "channel"},
		nil,
	)
	voltageDesc = prometheus.NewDesc(
		"shellyem_emeter_voltage",
		"shellyem_exporter: Measured voltage.",
		[]string{"host", "channel"},
		nil,
	)
	totalDesc = prometheus.NewDesc(
		"shellyem_emeter_total",
		"shellyem_exporter: Total energy.",
		[]string{"host", "channel"},
		nil,
	)
	totalReturnedDesc = prometheus.NewDesc(
		"shellyem_emeter_total_returned",
		"shellyem_exporter: Total returned.",
		[]string{"host", "channel"},
		nil,
	)
	reactiveDesc = prometheus.NewDesc(
		"shellyem_emeter_reactive",
		"shellyem_exporter: Measured reactive power.",
		[]string{"host", "channel"},
		nil,
	)
)

func (c *Collector) Describe(descChan chan<- *prometheus.Desc) {
	descChan <- scrapeDurationDesc
	descChan <- scrapeSuccessDesc
	descChan <- instantPowerDesc
	descChan <- voltageDesc
}

func (c *Collector) Collect(metricChan chan<- prometheus.Metric) {
	start := time.Now()
	shellyStatus, err := c.getShellyStatus("192.168.178.21")
	duration := time.Since(start)
	metricChan <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds())
	if err != nil {
		metricChan <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, 0)
		return
	}
	metricChan <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, 1)

	for i, e := range shellyStatus.Emeters {
		metricChan <- prometheus.MustNewConstMetric(instantPowerDesc,
			prometheus.GaugeValue, e.Power,
			c.Host, strconv.Itoa(i))

		metricChan <- prometheus.MustNewConstMetric(voltageDesc,
			prometheus.GaugeValue, e.Voltage,
			c.Host, strconv.Itoa(i))

		metricChan <- prometheus.MustNewConstMetric(totalDesc,
			prometheus.CounterValue, e.Total,
			c.Host, strconv.Itoa(i))

		metricChan <- prometheus.MustNewConstMetric(totalReturnedDesc,
			prometheus.CounterValue, e.TotalReturned,
			c.Host, strconv.Itoa(i))

		metricChan <- prometheus.MustNewConstMetric(reactiveDesc,
			prometheus.GaugeValue, e.Reactive,
			c.Host, strconv.Itoa(i))
	}

}
