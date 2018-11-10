package collectors

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"db-exporter/config"
	"log"
	"strconv"
)

const(
	slowLaunchTimeQuery=`SHOW variables like 'slow_launch_time'`
	slowLaunchThreadsQuery=`SHOW global status like 'Slow_launch_threads'`
	slowQuery=`SHOW global status like 'Slow_queries'`
	longTimeQuery=`SHOW variables like 'long_query_time'`
)
var(
	longQueryTimeDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","long_query_time"),"",nil,nil)
	slowLaunchTimeDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","slow_launch_time"),"",nil,nil)
	slowLaunchThreadsDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","slow_launch_threads"),"",nil,nil)
	slowDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","slow"),"",nil,nil)
	)
func ScrapeSlow(db *sql.DB,ch chan<- prometheus.Metric) error {
	config.DoQueryWithTwoResult(slowLaunchTimeDesc,db,ch,slowLaunchTimeQuery)
	config.DoQueryWithTwoResult(slowLaunchThreadsDesc,db,ch,slowLaunchThreadsQuery)
	config.DoQueryWithTwoResult(slowDesc,db,ch,slowQuery)
	keyRows,err := db.Query(longTimeQuery)
	if err!=nil {
		log.Printf("error:",err)
		return nil
	}
	defer keyRows.Close()
	var str string
	var value string
	for keyRows.Next(){
		if err:= keyRows.Scan(&str,&value);err!=nil {
			log.Printf("error:",err)
			return nil
		}
		vv,s := strconv.ParseFloat(value,64)
		if s!=nil {
			log.Printf("error:",err)
			return nil
		}
		ch<-prometheus.MustNewConstMetric(longQueryTimeDesc,prometheus.GaugeValue,vv,)
	}
	return nil
}
