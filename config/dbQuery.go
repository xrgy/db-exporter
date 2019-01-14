package config

import (
	"os"
	"log"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/coreos/etcd/client"
	"context"
	"github.com/ghodss/yaml"
	yaml2 "gopkg.in/yaml.v2"
)

var db *sql.DB

type ETCDParameter struct {
	Kind     string            `yaml:"kind"`
	Metadata map[string]string `yaml:"metadata"`
	Spec     Spec              `yaml:"spec"`
	Status   map[string]string `yaml:"status"`
}
type Spec struct {
	Ports     map[string]string `yaml:"ports"`
	Selector  map[string]string `yaml:"selector"`
	ClusterIP string            `yaml:"clusterIP"`
	Stype     string            `yaml:"type"`
}
type ConnectInfoData struct {
	Uuid        string
	IP          string
	Databasename string
	Username string
	Password string
	Port string
	//Params_maps map[string]string
}
type Monitor_info []byte
type ConnectInfo struct {
	uuid   string
	ip     string
	databasename string
	username string
	password string
	port string
	m_info Monitor_info
}

func readEtcdInfo(cfg client.Config, servicename string) string {
	c, err := client.New(cfg)
	if err != nil {
		log.Printf("%s", err.Error())
		return ""
	}
	//m := make(map[string]string)
	kapi := client.NewKeysAPI(c)
	resp1, err := kapi.Get(context.Background(), "/registry/services/specs/default/"+servicename, nil)
	if err != nil {
		return ""
	} else {
		log.Printf("etcd node value:"+resp1.Node.Value)
		param := &ETCDParameter{}
		v_rw := []byte(resp1.Node.Value)
		y_rw, err := yaml.JSONToYAML(v_rw)
		if err != nil {
			log.Printf("%s", err.Error())
			return ""
		}

		yaml2.Unmarshal(y_rw, &param)
		return param.Spec.ClusterIP

	}

}
func GetDBHandle() *sql.DB {
	var err error
	//cfg := client.Config{
	//	Endpoints:               []string{"http://" + os.Getenv("ETCD_ENDPOINT")},
	//	Transport:               client.DefaultTransport,
	//	HeaderTimeoutPerRequest: time.Second,
	//}
	DBUsername := os.Getenv("DB_USERNAME")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBEndpoint := os.Getenv("DB_ENDPOINT")
	DBDatabase := os.Getenv("DB_DATABASE")
	//servicename := strings.Split(DBEndpoint, ":")[0]
	//serviceport := strings.Split(DBEndpoint, ":")[1]
	//ip := readEtcdInfo(cfg, servicename)
	//dsn := DBUsername + ":" + DBPassword + "@(" + ip + ":" + serviceport + ")/" + DBDatabase
	dsn := DBUsername + ":" + DBPassword + "@(" + DBEndpoint + ")/" + DBDatabase

	log.Printf("connectionurl:"+dsn)
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
func GetMonitorInfo(id string) (ConnectInfoData, error) {
	info, err := queryConnectInfo(id)
	if err != nil {
		log.Printf("queryConnectInfo error", err.Error())
		return ConnectInfoData{}, nil
	}
	//m := info.m_info
	//m_info_map := make(map[string]string)
	//if len(m) != 0 {
	//	err := json.Unmarshal(m, &m_info_map)
	//	if err != nil {
	//		log.Printf("Unmarshal error")
	//	}
	//}
	con_info_data := ConnectInfoData{
		"",
		info.ip,
		info.databasename,
		info.username,
		info.password,
		info.port,
	}
	return con_info_data, nil
}
func queryConnectInfo(id string) (ConnectInfo, error) {
	rows, err := db.Query("select ip,databasename,username,password,port from tbl_db_monitor_record where uuid=?", id)
	info := ConnectInfo{}
	if err != nil {
		log.Printf("query error")
		return info, err
	} else {
		for rows.Next() {
			err = rows.Scan(&info.ip,&info.databasename,&info.username,&info.password,&info.port)
		}
		defer rows.Close()
		return info, nil
	}
}
func CloseDBHandle() {
	db.Close()
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
