package collectors

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
	"db-exporter/config"
)

const (
	cachedThreadQuery = `SHOW GLOBAL STATUS where Variable_name in ('Threads_cached')`
	connectedThreadQuery = `SHOW GLOBAL STATUS where Variable_name in ('Threads_connected')`
	createdThreadQuery = `SHOW GLOBAL STATUS where Variable_name in ('Threads_created')`
	runningThreadQuery = `SHOW GLOBAL STATUS where Variable_name in ('Threads_running')`
	cacheSizeThreadQuery = `SHOW GLOBAL STATUS where Variable_name in ('thread_cache_size')`
	)
var (
	cachedThreadDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","cachedThread"),"",nil,nil)
	connectedThreadDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","Threads_connected"),"",nil,nil)
	createdThreadDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","Threads_created"),"",nil,nil)
	runningThreadDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","runningThread"),"",nil,nil)
	cacheSizeThreadDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","thread_cache_size"),"",nil,nil)
	)
func ScrapeThreadSchema(db *sql.DB,ch chan<- prometheus.Metric) error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		config.DoQueryWithTwoResult(cachedThreadDesc,db,ch,cachedThreadQuery)
	}()
	wg.Add(1)
	go func() {
		wg.Done()
		config.DoQueryWithTwoResult(connectedThreadDesc,db,ch,connectedThreadQuery)
	}()
	wg.Add(1)
	go func() {
		wg.Done()
		config.DoQueryWithTwoResult(createdThreadDesc,db,ch,createdThreadQuery)
	}()
	wg.Add(1)
	go func() {
		wg.Done()
		config.DoQueryWithTwoResult(runningThreadDesc,db,ch,runningThreadQuery)
	}()
	wg.Add(1)
	go func() {
		wg.Done()
		config.DoQueryWithTwoResult(cacheSizeThreadDesc,db,ch,cacheSizeThreadQuery)
	}()
	wg.Wait()
	return nil
}
