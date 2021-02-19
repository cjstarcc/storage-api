package objects

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"storageApi/conf"
	"storageApi/util/heartbeart"
	"storageApi/util/objectstream"
	"strings"
)

func Handler(w http.ResponseWriter,r *http.Request){
	m :=r.Method
	if m == http.MethodPut{
		Put(w,r)
		return
	}
	if m == http.MethodGet {
		Get(w,r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func Put(w http.ResponseWriter,r *http.Request)  {
	object := strings.Split(r.URL.EscapedPath(),"/")[2]
	c,err := storeObject(r.Body,object)
	if err!=nil{
		log.Println(err)
	}
	w.WriteHeader(c)
}

func storeObject(r io.Reader,object string) (int,error) {
	stream,err := putStream(object)
	if err!=nil{
		return http.StatusServiceUnavailable,err
	}
	io.Copy(stream,r)
	err = stream.Close()
	if err!=nil{
		return http.StatusInternalServerError,err
	}
	return http.StatusOK,nil
}

func putStream(object string) (*objectstream.PutStream,error) {
	server := heartbeart.ChooseRandomDataServer()
	if server == ""{
		return nil ,fmt.Errorf("Cannot find any dataServer")
	}
	return objectstream.NewPutStream(server,object),nil
}

func Get(w http.ResponseWriter,r *http.Request)  {
	f,err:=os.Open(conf.GetConfig().Env.Dir+"/objects/"+strings.Split(r.URL.EscapedPath(),"/")[2])
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	io.Copy(w,f)
}

