package collectors

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"strings"
)

const(
	lockQuery = `SHOW GLOBAL STATUS where Variable_name in('Table_locks_immediate','Table_locks_waited')`
)
var (
	table_locks_waitDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","Table_locks_waited"),"",nil,nil)
	table_locks_immediateDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","Table_locks_immediate"),"",nil,nil)
	)

func ScrapeLockSchema(db *sql.DB,ch chan<- prometheus.Metric) error {
	listRows,err := db.Query(lockQuery)
	if err!=nil {
		log.Printf("query error:",err)
		return err
	}
	defer listRows.Close()
	var str string
	var value uint64
	for listRows.Next() {
		if err := listRows.Scan(&str,&value); err !=nil {
			log.Printf("error:",err)
			return err
		}
		str = strings.ToLower(str)
		if strings.Contains(str,"waited") {
			ch<-prometheus.MustNewConstMetric(table_locks_waitDesc,prometheus.GaugeValue,float64(value),)
		}
		if strings.Contains(str,"immediate") {
			ch<-prometheus.MustNewConstMetric(table_locks_immediateDesc,prometheus.GaugeValue,float64(value),)
		}
	}
	return nil
}
