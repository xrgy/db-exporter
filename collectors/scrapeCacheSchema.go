package collectors

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"db-exporter/config"
	"log"
	"sync"
)

const (
	query_cache_min_res_unitQuery = `select @@global.query_cache_min_res_unit`
	query_cache_typeQuery = `select @@global.query_cache_type`
	query_cache_wlock_invalidateQuery = `select @@global.query_cache_wlock_invalidate`
	cacheSizeQuery = `select @@global.query_cache_size`
	cacheLimitQuery = `select @@global.query_cache_limit`
	cacheHitsQuery = `SHOW GLOBAL STATUS where Variable_name in('Qcache_hits')`
	qcachefree_blocksQuery = `SHOW GLOBAL STATUS where Variable_name in('_blocks')`
	qcachefree_memoryQuery = `SHOW GLOBAL STATUS where Variable_name in('Qcache_free_memory')`
	qcachetotal_blocksQuery = `SHOW GLOBAL STATUS where Variable_name in('Qcache_total_blocks')`
	qcache_insertsQuery = `SHOW GLOBAL STATUS where Variable_name in('Qcache_inserts')`
	qcachelowmem_prunesQuery = `SHOW GLOBAL STATUS where Variable_name in('Qcache_lowmem_prunes')`
	qcachenot_cachedQuery = `SHOW GLOBAL STATUS where Variable_name in('Qcache_not_cached')`
	qcachequeries_in_cacheQuery = `SHOW GLOBAL STATUS where Variable_name in('Qcache_queries_in_cache')`
)
var(
	query_cache_wlock_invalidateDesc=prometheus.NewDesc(prometheus.BuildFQName("mysql","","query_cache_wlock_invalidate"),"",nil,nil)
	query_cache_typeDesc=prometheus.NewDesc(prometheus.BuildFQName("mysql","","query_cache_type"),"",nil,nil)
	query_cache_min_res_unitDesc=prometheus.NewDesc(prometheus.BuildFQName("mysql","","query_cache_min_res_unit"),"",nil,nil)
	cacheSizeDesc=prometheus.NewDesc(prometheus.BuildFQName("mysql","","query_cache_size"),"",nil,nil)
	cacheLimitDesc=prometheus.NewDesc(prometheus.BuildFQName("mysql","","query_cache_limit"),"",nil,nil)
	cacheHitsDesc=prometheus.NewDesc(prometheus.BuildFQName("mysql","","cache_hits"),"",nil,nil)
	Qcache_free_blocksDesc=prometheus.NewDesc(prometheus.BuildFQName("mysql","","Qcache_free_blocks"),"",nil,nil)
	Qcache_free_memoryDesc=prometheus.NewDesc(prometheus.BuildFQName("mysql","","Qcache_free_memory"),"",nil,nil)
	Qcache_insertsDesc=prometheus.NewDesc(prometheus.BuildFQName("mysql","","Qcache_inserts"),"",nil,nil)
	Qcache_lowmmem_prunesDesc=prometheus.NewDesc(prometheus.BuildFQName("mysql","","Qcache_lowmem_prunes"),"",nil,nil)
	Qcache_not_cachedDesc=prometheus.NewDesc(prometheus.BuildFQName("mysql","","Qcache_not_cached"),"",nil,nil)
	Qcache_queries_in_cacheDesc=prometheus.NewDesc(prometheus.BuildFQName("mysql","","Qcache_queries_in_cache"),"",nil,nil)
	Qcache_total_blocksDesc=prometheus.NewDesc(prometheus.BuildFQName("mysql","","Qcache_total_blocks"),"",nil,nil)
	)
func ScrapeCacheSchema(db *sql.DB,ch chan<- prometheus.Metric) error {
	config.DoQueryWithOneResult(query_cache_wlock_invalidateDesc,db,ch,query_cache_wlock_invalidateQuery)
	listRows,err := db.Query(query_cache_typeQuery)
	if err!=nil {
		log.Printf("query error:",err)
	}else {
		var str string
		for listRows.Next() {
			if err := listRows.Scan(&str); err !=nil {
				log.Printf("error:",err)
			}
		}
		value:=1
		if str == "OFF" {
			value = 0
		}
		ch <- prometheus.MustNewConstMetric(query_cache_typeDesc,prometheus.GaugeValue,float64(value),)
	}
	defer listRows.Close()
	config.DoQueryWithOneResult(query_cache_min_res_unitDesc,db,ch,query_cache_min_res_unitQuery)
	config.DoQueryWithTwoResult(Qcache_free_blocksDesc,db,ch,qcachefree_blocksQuery)
	config.DoQueryWithTwoResult(Qcache_free_memoryDesc,db,ch,qcachefree_memoryQuery)
	config.DoQueryWithTwoResult(Qcache_insertsDesc,db,ch,qcache_insertsQuery)
	config.DoQueryWithTwoResult(Qcache_lowmmem_prunesDesc,db,ch,qcachelowmem_prunesQuery)
	config.DoQueryWithTwoResult(Qcache_not_cachedDesc,db,ch,qcachenot_cachedQuery)
	config.DoQueryWithTwoResult(Qcache_queries_in_cacheDesc,db,ch,qcachequeries_in_cacheQuery)
	config.DoQueryWithTwoResult(Qcache_total_blocksDesc,db,ch,qcachetotal_blocksQuery)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		config.DoQueryWithOneResult(cacheSizeDesc,db,ch,cacheSizeQuery)
	}()
	go func() {
		defer wg.Done()
		config.DoQueryWithOneResult(cacheLimitDesc,db,ch,cacheLimitQuery)
	}()
	go func() {
		defer wg.Done()
		listRows,err := db.Query(cacheHitsQuery)
		if err!=nil {
			log.Printf("query error:",err)
			return
		}
		defer listRows.Close()
		var str string
		var value uint64
		for listRows.Next() {
			if err := listRows.Scan(&str,&value); err !=nil {
				log.Printf("error:",err)
				return
			}
			ch <- prometheus.MustNewConstMetric(cacheHitsDesc,prometheus.GaugeValue,float64(value),)
		}
	}()
	wg.Wait()
	return nil
}
