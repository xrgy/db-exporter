package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/gorilla/mux"
	"strings"
	"fmt"
	"db-exporter/config"
	"db-exporter/collectors/api"
	"db-exporter/collectors"
)


var listenAddress = kingpin.Flag("web.listen-address","Address to listen on for web insterface and " +
	"telemetry.").Default(":9103").String()


func init()  {
	config.GetDBHandle()
}
func runCollector(collector prometheus.Collector,w http.ResponseWriter,r *http.Request)  {
	registry:= prometheus.NewRegistry()
	registry.MustRegister(collector)
	gatherers := prometheus.Gatherers{
		prometheus.DefaultGatherer,
		registry,
	}
	h:=promhttp.HandlerFor(gatherers,promhttp.HandlerOpts{})
	h.ServeHTTP(w,r)
}
func main() {
	kingpin.Parse()
	r := mux.NewRouter()
	r.HandleFunc("/mysql",handler)
	r.HandleFunc("/api/v1/mysql/access",api.MysqlAccess).Methods(http.MethodPost)
	http.ListenAndServe(*listenAddress,r)
}
func handler(w http.ResponseWriter,r *http.Request)  {
	var collectorType prometheus.Collector
	target:= r.URL.Query().Get("target")
	if target=="" {
		http.Error(w,"'target' parameter must be specified",400)
		return
	}
	atr:=strings.Split(fmt.Sprintf("%s",r.URL),"?")[0]
	fmt.Sprintf(atr)
	switch strings.Split(fmt.Sprintf("%s",r.URL),"?")[0] {
	case "/mysql":
		collectorType = collectors.MysqlCollector{target}
		break
	default:
		break
	}
	runCollector(collectorType,w,r)
}

