package collectors

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"db-exporter/config"
)

const(
	nameQuery = `SHOW variables like 'version_comment'`
	versionQuery = `SHOW variables like 'version'`
	portQuery = `SHOW variables like 'port'`
	baseDirQuery = `SHOW variables like 'basedir'`
	dataDirQUery  = `SHOW variables like 'datadir'`
)
var (
	mysqldUp = prometheus.NewDesc(prometheus.BuildFQName("mysql","","up"),
		"whether the MySQL server is up",[]string{"mysql_name","mysql_version","mysql_port","mysql_basedir","mysql_datadir"},nil)
)
func ScrapeInfoSchema(db *sql.DB,ch chan<- prometheus.Metric) error {
	name := config.DoQueryWithOneStringResult(db,nameQuery)
	version := config.DoQueryWithOneStringResult(db,versionQuery)
	port := config.DoQueryWithOneStringResult(db,portQuery)
	baseDir := config.DoQueryWithOneStringResult(db,baseDirQuery)
	dataDir := config.DoQueryWithOneStringResult(db,dataDirQUery)
	isUpRows,err := db.Query("SELECT 1")
	if err!=nil {
		ch <- prometheus.MustNewConstMetric(mysqldUp,prometheus.GaugeValue,0,name,version,port,baseDir,dataDir)
		return err
	}
	isUpRows.Close()
	ch <- prometheus.MustNewConstMetric(mysqldUp,prometheus.GaugeValue,1,name,version,port,baseDir,dataDir)
	return nil
}
