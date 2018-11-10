package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
	"db-exporter/config"
	"database/sql"
	"fmt"
	"time"
	"log"
)

type MysqlCollector struct {
	Target string
}
var (
	before int64
	after int64
	mysql_monitorstatus = prometheus.NewDesc(prometheus.BuildFQName("mysql","","monitorstatus"),
		"mysql monitorstatus",nil,nil)
	mysql_scrapeDurationDesc = prometheus.NewDesc(prometheus.BuildFQName("mysql","","scrape_duration_seconds"),
		"mysql: Duration of a collector scrape.",nil,nil)
	mysql_scrapeErrors = prometheus.NewCounter(prometheus.CounterOpts{"mysql","","scrape_errors_total",
	"Total number of times an errors occured scraping a mysql",nil})
)
func (c MysqlCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("dummy", "dummy", nil, nil)
}

func (c MysqlCollector) Collect(ch chan<- prometheus.Metric) {
	getMysqlData(ch,c.Target)
}
func getMysqlData(ch chan<- prometheus.Metric,target string) error {
	monitor_info := config.GetMonitorInfo(target)
	ip := monitor_info.IP
	username := monitor_info.Params_maps["username"]
	password := monitor_info.Params_maps["password"]
	databasename := monitor_info.Params_maps["databasename"]
	port := monitor_info.Params_maps["port"]
	if ip =="" || username == "" ||password == "" || port == "" || databasename == "" {
		mysql_scrapeErrors.Inc()
		ch <- prometheus.MustNewConstMetric(mysql_monitorstatus,prometheus.GaugeValue,0)
		return nil
	}
	driverDataSource := username + ":"+password+"@("+ip+":"+port+")/"+databasename
	before = time.Now().UnixNano()
	db,err := sql.Open("mysql",driverDataSource)
	if err!=nil {
		mysql_scrapeErrors.Inc()
		fmt.Sprintf("Error opening connection to database:",err)
		ch <- prometheus.MustNewConstMetric(mysql_monitorstatus,prometheus.GaugeValue,0)
		return nil
	}
	err = db.Ping()
	if err!=nil {
		mysql_scrapeErrors.Inc()
		fmt.Sprintf("mysql error:",err)
		ch <- prometheus.MustNewConstMetric(mysql_monitorstatus,prometheus.GaugeValue,0)
		return nil
	}
	ch <- prometheus.MustNewConstMetric(mysql_monitorstatus,prometheus.GaugeValue,1)
	after = time.Now().UnixNano()
	mysqlScrapeTarget(db,ch)
	return nil
}
func mysqlScrapeTarget(db *sql.DB, ch chan<- prometheus.Metric) {
	var err error
	if err = ScrapeInfoSchema(db,ch);err !=nil{
		mysql_scrapeErrors.Inc()
		log.Printf("Error scraping for collect.infoSchema")
	}
	if err = ScrapeOperationSchema(db,ch);err !=nil{
		mysql_scrapeErrors.Inc()
		log.Printf("Error scraping for collect.operationSchema")
	}
	if err = ScrapeConnectionSchema(db,ch);err !=nil{
		mysql_scrapeErrors.Inc()
		log.Printf("Error scraping for collect.connectionSchema")
	}
	if err = ScrapeRequestSchema(db,ch);err !=nil{
		mysql_scrapeErrors.Inc()
		log.Printf("Error scraping for collect.requestSchema")
	}
	if err = ScrapeThreadSchema(db,ch);err !=nil{
		mysql_scrapeErrors.Inc()
		log.Printf("Error scraping for collect.threadSchema")
	}
	if err = ScrapeOpenFiles(db,ch);err !=nil{
		mysql_scrapeErrors.Inc()
		log.Printf("Error scraping for collect.openfiels")
	}
	if err = ScrapeTable(db,ch);err !=nil{
		mysql_scrapeErrors.Inc()
		log.Printf("Error scraping for collect.table")
	}
	if err = ScrapeSlow(db,ch);err !=nil{
		mysql_scrapeErrors.Inc()
		log.Printf("Error scraping for collect.slow")
	}
	if err = ScrapeKeySchema(db,ch);err !=nil{
		mysql_scrapeErrors.Inc()
		log.Printf("Error scraping for collect.keySchema")
	}
	if err = ScrapeCacheSchema(db,ch);err !=nil{
		mysql_scrapeErrors.Inc()
		log.Printf("Error scraping for collect.cacheSchema")
	}
	if err = ScrapeDBSchema(db,ch);err !=nil{
		mysql_scrapeErrors.Inc()
		log.Printf("Error scraping for collect.dbSchema")
	}
	if err = ScrapeLockSchema(db,ch);err !=nil{
		mysql_scrapeErrors.Inc()
		log.Printf("Error scraping for collect.lockSchema")
	}
	if err = ScrapeHandlerSchema(db,ch);err !=nil{
		mysql_scrapeErrors.Inc()
		log.Printf("Error scraping for collect.handlerSchema")
	}
	if err = ScrapeResponseSchema(db,ch,after-before);err !=nil{
		mysql_scrapeErrors.Inc()
		log.Printf("Error scraping for collect.responseSchema")
	}
	db.Close()
}
