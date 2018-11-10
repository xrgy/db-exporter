package collectors

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
	"db-exporter/config"
)

const (
	bytesReceiveQuery = `SHOW STATUS where Variable_name in('Bytes_received')`
	bytesSentQuery = `SHOW STATUS where Variable_name in('Bytes_sent')`
	questionsQuery = `SHOW GLOBAL STATUS where Variable_name in('Questions')`
	)
var (
	bytesReceivedDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","bytesReceived"),"",nil,nil)
	bytesSentDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","bytesSent"),"",nil,nil)
	questionsDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","questions"),"",nil,nil)

	)
func ScrapeRequestSchema(db *sql.DB,ch chan<- prometheus.Metric) error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		config.DoQueryWithTwoResult(bytesReceivedDesc,db,ch,bytesReceiveQuery)
	}()
	wg.Add(1)
	go func() {
		wg.Done()
		config.DoQueryWithTwoResult(bytesSentDesc,db,ch,bytesSentQuery)
	}()
	wg.Add(1)
	go func() {
		wg.Done()
		config.DoQueryWithTwoResult(questionsDesc,db,ch,questionsQuery)
	}()
	wg.Wait()
	return nil
}
