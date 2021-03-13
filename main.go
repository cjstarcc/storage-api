package main

import (
	"log"
	"net/http"
	"storageApi/conf"
	"storageApi/util/heartbeart"
	"storageApi/util/www/locate"
	"storageApi/util/www/objects"
	"storageApi/util/www/versions"
	"strconv"
)

func main() {
	conf.InitConfig()
	go heartbeart.ListenHeartBeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Hander)
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(conf.GetConfig().Env.Port), nil))
}
