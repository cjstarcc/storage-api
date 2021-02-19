package locate

import (
	"encoding/json"
	"net/http"
	"storageApi/conf"
	"storageApi/rabbitmq"
	"strconv"
	"strings"
	"time"
)

func Hander(w http.ResponseWriter,r *http.Request)  {
	m := r.Method
	if m != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	info := Locate(strings.Split(r.URL.EscapedPath(),"/")[2])
	if len(info) == 0{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b,_ := json.Marshal(info)
	w.Write(b)
}

func Locate(name string) string {
	q:=rabbitmq.New(conf.GetConfig().Env.Rabbitmq)
	q.Publish("dataservers",name)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	msg :=  <-c
	s,_ :=strconv.Unquote(string(msg.Body))
	return s
}
func Existt(name string) bool {
	return Locate(name) !=""
}