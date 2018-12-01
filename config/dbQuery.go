package config

import (
	"os"
	"log"
	"time"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
)

var db *sql.DB

type ConnectInfoData struct {
	Uuid        string
	IP          string
	Params_maps map[string]string
}
type Monitor_info []byte
type ConnectInfo struct {
	uuid   string
	ip     string
	m_info Monitor_info
}

func GetDBHandle() *sql.DB {
	var err error
	DBUsername := os.Getenv("DB_USERNAME")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBEndpoint := os.Getenv("DB_ENDPOINT")
	DBDatabase := os.Getenv("DB_DATABASE")
	dsn := DBUsername + ":" + DBPassword + "@(" + DBEndpoint + ")/" + DBDatabase
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("get DB handle error: %v", err)
	}
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(28000 * time.Second)
	err = db.Ping()
	if err != nil {
		log.Printf("connecting DB error: %v ", err)
	}
	return db
}
func GetMonitorInfo(id string) (ConnectInfoData,error) {
	info,err := queryConnectInfo(id)
	if err!=nil {
		log.Printf("queryConnectInfo error",err.Error())
		return ConnectInfoData{},nil
	}
	m := info.m_info
	m_info_map := make(map[string]string)
	if len(m) != 0 {
		err := json.Unmarshal(m, &m_info_map)
		if err != nil {
			log.Printf("Unmarshal error")
		}
	}
	con_info_data := ConnectInfoData{
		"",
		info.ip,
		m_info_map,
	}
	return con_info_data,nil
}
func queryConnectInfo(id string) (ConnectInfo,error) {
	rows, err := db.Query("select ip,monitor_info from tbl_monitor_record where uuid=?", id)
	info := ConnectInfo{}
	if err != nil {
		log.Printf("query error")
		return info,err
	} else {
		for rows.Next() {
			err = rows.Scan(&info.ip, &info.m_info)
		}
		defer rows.Close()
		return info,nil
	}
}
func CloseDBHandle() {
	db.Close()
}
func GetlldpMonitorInfo(id string) []ConnectInfoData {
	rows, err := db.Query("select uuid,ip,monitor_info from tbl_monitor_record where deleted=0 and middle_resource_type_id=?", id)
	if err != nil {
		log.Printf("query error")
	}
	infos := []ConnectInfo{}
	for rows.Next() {
		info := ConnectInfo{}
		err = rows.Scan(&info.uuid, &info.ip, &info.m_info)
		infos = append(infos, info)
	}
	defer rows.Close()
	conninfos := []ConnectInfoData{}
	for i := 0; i < len(infos); i++ {
		info := infos[i]
		m := info.m_info
		m_info_map := make(map[string]string)
		if len(m) != 0 {
			err := json.Unmarshal(m, &m_info_map)
			if err != nil {
				log.Printf("Unmarshal error")
			}
		}
		con_info_data := ConnectInfoData{
			info.uuid,
			info.ip,
			m_info_map,
		}
		conninfos = append(conninfos, con_info_data)
	}
	return conninfos
}

func DoQueryWithTwoResult(desc *prometheus.Desc, db *sql.DB, ch chan<- prometheus.Metric, querystring string) {
	keyRows, err := db.Query(querystring)
	if err != nil {
		log.Printf("query error:", err)
		return
	}
	defer keyRows.Close()
	var str string
	var value uint64
	for keyRows.Next() {
		if err := keyRows.Scan(&str, &value); err != nil {
			log.Printf("error:", err)
			return
		}
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(value), )
	}
}

func DoQueryWithOneResult(desc *prometheus.Desc, db *sql.DB, ch chan<- prometheus.Metric, querystring string) {
	listRows, err := db.Query(querystring)
	if err != nil {
		log.Printf("query error:", err)
		return
	}
	defer listRows.Close()
	var value uint64
	for listRows.Next() {
		if err := listRows.Scan(&value); err != nil {
			log.Printf("error:", err)
			return
		}
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(value), )
	}
}

func DoQueryWithOneStringResult(db *sql.DB, querystring string) string {
	listRows, err := db.Query(querystring)
	if err != nil {
		log.Printf("query error:", err)
		return ""
	}
	defer listRows.Close()
	var str string
	var value string
	for listRows.Next() {
		if err := listRows.Scan(&str, &value); err != nil {
			log.Printf("error:", err)
			return ""
		}
	}
	return value
}
