package collectors

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"db-exporter/config"
)

const (
	connectTimeoutQuery = `SHOW GLOBAL variables where Variable_name in ('connect_timeout')`
)
var (
	connectResponseTimeDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","connectResponseTime"),"",nil,nil)
	connectTimeoutDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","connectTimeout"),"",nil,nil)
	)
func ScrapeResponseSchema(db *sql.DB,ch chan<- prometheus.Metric,seconds int64) error {
	config.DoQueryWithTwoResult(connectTimeoutDesc,db,ch,connectTimeoutQuery)
	ch<-prometheus.MustNewConstMetric(connectResponseTimeDesc,prometheus.GaugeValue,float64(seconds))
	return nil
}
