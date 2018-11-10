package collectors

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"db-exporter/config"
)

const (
	handler_read_firstQuery = `SHOW GLOBAL STATUS like 'handler_read_first'`
	handler_read_keyQuery = `SHOW GLOBAL STATUS like 'handler_read_key'`
	handler_read_nextQuery = `SHOW GLOBAL STATUS like 'handler_read_next'`
	handler_read_prevQuery = `SHOW GLOBAL STATUS like 'handler_read_prev'`
	handler_read_rndQuery = `SHOW GLOBAL STATUS like 'handler_read_rnd'`
	handler_read_rnd_nextQuery = `SHOW GLOBAL STATUS like 'handler_read_rnd_next'`
	)
var (
	handler_read_firstDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","handler_read_first"),"",nil,nil)
	handler_read_keyDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","handler_read_key"),"",nil,nil)
	handler_read_nextDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","handler_read_next"),"",nil,nil)
	handler_read_prevDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","handler_read_prev"),"",nil,nil)
	handler_read_rndDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","handler_read_rnd"),"",nil,nil)
	handler_read_rnd_nextDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","handler_read_rnd_next"),"",nil,nil)
	)
func ScrapeHandlerSchema(db *sql.DB,ch chan<- prometheus.Metric) error {
	config.DoQueryWithTwoResult(handler_read_firstDesc,db,ch,handler_read_firstQuery)
	config.DoQueryWithTwoResult(handler_read_keyDesc,db,ch,handler_read_keyQuery)
	config.DoQueryWithTwoResult(handler_read_nextDesc,db,ch,handler_read_nextQuery)
	config.DoQueryWithTwoResult(handler_read_prevDesc,db,ch,handler_read_prevQuery)
	config.DoQueryWithTwoResult(handler_read_rndDesc,db,ch,handler_read_rndQuery)
	config.DoQueryWithTwoResult(handler_read_rnd_nextDesc,db,ch,handler_read_rnd_nextQuery)
	return nil
}
