package collectors

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"db-exporter/config"
)

const (
	tmpDiskTablesQuery=`SHOW global status like 'Created_tmp_disk_tables'`
	tmpFilesQuery=`SHOW global status like 'Created_tmp_files'`
	tmpTablesQuery=`SHOW global status like 'Created_tmp_tables'`
	tmpTableSizeQuery=`SHOW variables like 'tmp_table_size'`
	max_heap_table_sizeQuery=`SHOW variables like 'max_heap_table_size'`
	openTablesQuery=`SHOW global status like 'open_tables'`
	openedTablesQuery=`SHOW global status like 'opened_tables'`
	tableCacheQuery=`SHOW variables like 'table_open_cache'`
	)
var (
	tmpDiskTablesDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","Created_tmp_disk_tables"),"",nil,nil)
	tmpFilesDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","Created_tmp_files"),"",nil,nil)
	tmpTableSizeDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","Created_tmp_table_size"),"",nil,nil)
	max_heap_table_sizeDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","max_heap_table_size"),"",nil,nil)
	tmpTablesDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","Created_tmp_tables"),"",nil,nil)
	openTablesDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","openTables"),"",nil,nil)
	openedTablesDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","openedTables"),"",nil,nil)
	tableCacheDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","table_cache"),"",nil,nil)
	)
func ScrapeTable(db *sql.DB,ch chan<- prometheus.Metric) error {
	config.DoQueryWithTwoResult(tmpDiskTablesDesc,db,ch,tmpDiskTablesQuery)
	config.DoQueryWithTwoResult(tmpFilesDesc,db,ch,tmpFilesQuery)
	config.DoQueryWithTwoResult(tmpTablesDesc,db,ch,tmpTablesQuery)
	config.DoQueryWithTwoResult(tmpTableSizeDesc,db,ch,tmpTableSizeQuery)
	config.DoQueryWithTwoResult(max_heap_table_sizeDesc,db,ch,max_heap_table_sizeQuery)
	config.DoQueryWithTwoResult(openTablesDesc,db,ch,openTablesQuery)
	config.DoQueryWithTwoResult(openedTablesDesc,db,ch,openedTablesQuery)
	config.DoQueryWithTwoResult(tableCacheDesc,db,ch,tableCacheQuery)
	return  nil
}
