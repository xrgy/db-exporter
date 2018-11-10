package collectors

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"db-exporter/config"
	"sync"
)

const (
	keyBufferSizeQuery = `SHOW GLOBAL VARIABLES where variable_name in('key_buffer_size')`
	keyBlockUsedQuery = `SHOW GLOBAL STATUS where variable_name in('key_blocks_used')`
	keyBlockUnusedQuery = `SHOW GLOBAL STATUS where variable_name in('key_blocks_unused')`
	keyBlockSizeQuery = `SHOW GLOBAL VARIABLES where variable_name in('key_cache_block_size')`
	keyReadQuery = `SHOW GLOBAL STATUS where variable_name in('key_reads')`
	keyReadRequestQuery = `SHOW GLOBAL STATUS where variable_name in('key_read_requests')`
	keyWriteQuery = `SHOW GLOBAL STATUS where variable_name in('key_writes')`
	keyWriteRequestQuery = `SHOW GLOBAL STATUS where variable_name in('key_write_requests')`
	)
var (
	keyBufferSizeDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","keyBufferSize"),"",nil,nil)
	keyBlockUsedDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","keyBlockUsed"),"",nil,nil)
	keyBlockUnUsedDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","keyBlocksUnused"),"",nil,nil)
	keyBlockSizeDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","keyBlockSize"),"",nil,nil)
	keyReadDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","keyRead"),"",nil,nil)
	keyReadRequestDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","keyReadRequest"),"",nil,nil)
	keyWriteDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","keyWrite"),"",nil,nil)
	keyWriteRequestDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","keyWriteRequest"),"",nil,nil)
	)
func ScrapeKeySchema(db *sql.DB,ch chan<- prometheus.Metric) error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		config.DoQueryWithTwoResult(keyBufferSizeDesc,db,ch,keyBufferSizeQuery)
	}()
	wg.Add(1)
	go func() {
		wg.Done()
		config.DoQueryWithTwoResult(keyBlockUsedDesc,db,ch,keyBlockUsedQuery)
	}()
	wg.Add(1)
	go func() {
		wg.Done()
		config.DoQueryWithTwoResult(keyBlockUnUsedDesc,db,ch,keyBlockUnusedQuery)
	}()
	wg.Add(1)
	go func() {
		wg.Done()
		config.DoQueryWithTwoResult(keyBlockSizeDesc,db,ch,keyBlockSizeQuery)
	}()
	wg.Add(1)
	go func() {
		wg.Done()
		config.DoQueryWithTwoResult(keyReadDesc,db,ch,keyReadQuery)
	}()
	wg.Add(1)
	go func() {
		wg.Done()
		config.DoQueryWithTwoResult(keyReadRequestDesc,db,ch,keyReadRequestQuery)
	}()
	wg.Add(1)
	go func() {
		wg.Done()
		config.DoQueryWithTwoResult(keyWriteDesc,db,ch,keyWriteQuery)
	}()
	wg.Add(1)
	go func() {
		wg.Done()
		config.DoQueryWithTwoResult(keyWriteRequestDesc,db,ch,keyWriteRequestQuery)
	}()
	wg.Wait()
	return nil
}
