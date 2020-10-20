package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var (
	isDeploy = os.Getenv("IS_DEPLOY")
)

func deployOrLocal(local string, deploy string) string {
	if isDeploy != "" {
		return deploy
	} else {
		return local
	}
}

func addProxy(prefix string, path string) error {
	target, err := url.Parse(path)
	if err != nil {
		return err
	}
	http.Handle(prefix, http.StripPrefix(prefix, httputil.NewSingleHostReverseProxy(target)))
	return nil
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func main() {
	panicIf(addProxy("/user/", deployOrLocal("http://localhost:8080/", "http://auth:8080/")))
	panicIf(addProxy("/picture/", deployOrLocal("http://localhost:8081/", "http://web:8081/")))
	panicIf(addProxy("/upload/", deployOrLocal("http://localhost:8083/", "http://uploader:8083/")))
	panicIf(addProxy("/els/", deployOrLocal("http://localhost:9200/", "http://elastic:9200/")))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))
	http.HandleFunc("/", hello)
	//fmt.Print("started server")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil))
}
