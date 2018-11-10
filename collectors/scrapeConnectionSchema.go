package collectors

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
	"db-exporter/config"
)

const (
	connectionQuery = `SHOW STATUS WHERE Variable_name in ('Threads_connected')`
	abortedConnectsQuery = `SHOW STATUS WHERE Variable_name in ('Aborted_connects')`
	abortedClientsQuery = `SHOW STATUS WHERE Variable_name in ('Aborted_clients')`
	maxConnectionsQuery = `SHOW variables like 'max_connections'`
	maxUsedConnectionsQuery = `SHOW GLOBAL STATUS like 'max_used_connections'`
	)

var (
	connectionDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","connections"),"",nil,nil)
	abortedConnectsDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","abortedConnects"),"",nil,nil)
	abortedClientsDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","abortedClients"),"",nil,nil)
	maxConnectionsDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","maxConnections"),"",nil,nil)
	maxUsedConnectionsDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","maxUsedConnections"),"",nil,nil)
	)
func ScrapeConnectionSchema(db *sql.DB,ch chan<- prometheus.Metric) error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		config.DoQueryWithTwoResult(connectionDesc,db,ch,connectionQuery)
	}()
	wg.Add(1)
	go func() {
		wg.Done()
		config.DoQueryWithTwoResult(abortedConnectsDesc,db,ch,abortedConnectsQuery)
	}()
	wg.Add(1)
	go func() {
		wg.Done()
		config.DoQueryWithTwoResult(abortedClientsDesc,db,ch,abortedClientsQuery)
	}()
	wg.Add(1)
	go func() {
		wg.Done()
		config.DoQueryWithTwoResult(maxConnectionsDesc,db,ch,maxConnectionsQuery)
	}()
	wg.Add(1)
	go func() {
		wg.Done()
		config.DoQueryWithTwoResult(maxUsedConnectionsDesc,db,ch,maxUsedConnectionsQuery)
	}()
	wg.Wait()
	return nil
}
