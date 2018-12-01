package collectors

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"db-exporter/config"
)

const(
	openFilesQuery = `SHOW global status like 'open_files'`
	openFilesLimitQuery = `SHOW variables like 'open_files_limit'`
	)
var (
	openFilesDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","openFiles"),"",nil,nil)
	openFilesLimitDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","openFilesLimit"),"",nil,nil)
	)
func ScrapeOpenFiles(db *sql.DB,ch chan<- prometheus.Metric) error {
	config.DoQueryWithTwoResult(openFilesDesc,db,ch,openFilesQuery)
	config.DoQueryWithTwoResult(openFilesLimitDesc,db,ch,openFilesLimitQuery)
	return nil
}
