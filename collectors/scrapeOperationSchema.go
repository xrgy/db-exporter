package collectors

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

const (
	operationQuery = `SHOW GLOBAL STATUS where variable_name in('com_select','com_delete','com_insert','com_update')`
)
var (
	selectTimesDesc= prometheus.NewDesc("mysql_select_total_times","mysql select total times",nil,nil)
	deleteTimesDesc= prometheus.NewDesc("mysql_delete_total_times","mysql delete total times",nil,nil)
	insertTimesDesc= prometheus.NewDesc("mysql_insert_total_times","mysql insert total times",nil,nil)
	updateTimesDesc= prometheus.NewDesc("mysql_update_total_times","mysql update total times",nil,nil)

)
func ScrapeOperationSchema(db *sql.DB,ch chan<- prometheus.Metric) error {
	operationListRows,err := db.Query(operationQuery)
	if err!=nil {
		return err
	}
	defer operationListRows.Close()
	var str string
	var value uint64
	for operationListRows.Next()  {
		if err := operationListRows.Scan(&str, &value); err !=nil {
			return err
		}
		str = strings.ToLower(str)
		if strings.Contains(str,"select") {
			ch <- prometheus.MustNewConstMetric(selectTimesDesc,prometheus.GaugeValue,float64(value),)
		}
		if strings.Contains(str,"insert") {
			ch <- prometheus.MustNewConstMetric(insertTimesDesc,prometheus.GaugeValue,float64(value),)
		}
		if strings.Contains(str,"delete") {
			ch <- prometheus.MustNewConstMetric(deleteTimesDesc,prometheus.GaugeValue,float64(value),)
		}
		if strings.Contains(str,"update") {
			ch <- prometheus.MustNewConstMetric(updateTimesDesc,prometheus.GaugeValue,float64(value),)
		}
	}
	return nil
}
