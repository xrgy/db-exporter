package collectors

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
	"fmt"
	"log"
)

const (
	dbSize = `SELECT concat(sum(DATA_LENGTH)) as data from information_schema.TABLES where table_schema = '%s'`
	tableTotalNumber = `SELECT count(*) as data from  information_schema.TABLES where table_schema = '%s'`
	dbQuery = `SELECT SCHEMA_NAME FROM information_schema.schemata`
)
var(
	dbSizeDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","information_schema","db_size"),
		"the size of db (bytes)",[]string{"mysql_info_schema"},nil)
	dbTableNumberDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","information_schema","table_numbers"),
		"the number of tables",[]string{"mysql_info_schema"},nil)
	)
func ScrapeDBSchema(db *sql.DB,ch chan<- prometheus.Metric) error {
	var dbList []string
	dbListRows,err := db.Query(dbQuery)
	if err!=nil {
		return err
	}
	defer dbListRows.Close()
	var database string
	for dbListRows.Next()  {
		if err := dbListRows.Scan(&database,);err!=nil {
			return err
		}
		dbList = append(dbList,database)
	}
	var wg sync.WaitGroup
	for _,database := range dbList{
		wg.Add(1)
		go func(databasestring string) {
			defer wg.Done()
			dbDetailQuery(dbTableNumberDesc,db,ch,databasestring,tableTotalNumber)
		}(database)
		wg.Add(1)
		go func(databasestring string) {
			defer wg.Done()
			dbDetailQuery(dbSizeDesc,db,ch,databasestring,dbSize)
		}(database)
	}
	wg.Wait()
	return nil
}
func dbDetailQuery(desc *prometheus.Desc, db *sql.DB, ch chan<- prometheus.Metric, databasestring string, queryString string) {
	dbSchemaRows,err := db.Query(fmt.Sprintf(queryString,databasestring))
	if err!=nil {
		log.Printf("error:",err)
		return
	}
	defer dbSchemaRows.Close()
	var data uint64
	for dbSchemaRows.Next() {
		err := dbSchemaRows.Scan(&data,)
		if err!=nil {
			log.Printf("error:",err)
			return
		}
		ch<-prometheus.MustNewConstMetric(desc,prometheus.GaugeValue,float64(data),databasestring,)
	}
}
