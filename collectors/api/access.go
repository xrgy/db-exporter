package api

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"os/exec"
	"strings"
	"errors"
	"database/sql"
	"fmt"
)
type AccessReq struct {
	monitorInfo map[string]string `json:"monitorInfo"`
}

type AccessResp struct {
	Result map[string]string `json:"result"`
}

func getResponse(w http.ResponseWriter,err error)  {
	resultMap := make(map[string]string,2)
	if err!=nil {
		resultMap["accessible"] = "false"
		resultMap["message"] = err.Error()
	}else {
		resultMap["accessible"] = "true"
		resultMap["message"] = ""
	}
	accessResp := AccessResp{Result:resultMap}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(accessResp)
}

func MysqlAccess(w http.ResponseWriter,r *http.Request)  {
	var info AccessReq
    body,err := ioutil.ReadAll(r.Body)
	if err!=nil {
		log.Printf("get body data error:"+err.Error())
	}
	err = json.Unmarshal(body,&info)
	ip := info.monitorInfo["ip"]
	username := info.monitorInfo["username"]
	password := info.monitorInfo["password"]
	databasename := info.monitorInfo["databasename"]
	port := info.monitorInfo["port"]
	if ip =="" || username == "" ||password == "" || port == "" || databasename == "" {
		log.Printf("illegal request bosy:%s",string(body))
		http.Error(w,"illegal request bosy",http.StatusBadRequest)
		return
	}
	command := "ping -i 0.3 -w 5 "+ip+" -c 3 | tail -n 2"
	cmd := exec.Command("/bin/sh","-c",command)
	ret,_ := cmd.Output()
	s := string(ret)
	if strings.Contains(s,"100% packet loss") {
		getResponse(w,errors.New("ip doesn't exist."))
		return
	}
	driverDataSource := username + ":"+password+"@("+ip+":"+port+")/"+databasename
	db,err := sql.Open("mysql",driverDataSource)
	if err!=nil {
		fmt.Sprintf("Error opening connection to database:",err)
		http.Error(w,"",http.StatusBadRequest)
	}
	err = db.Ping()
	getResponse(w,err)
}
