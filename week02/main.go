package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	var f http.HandlerFunc

	f = func(res http.ResponseWriter, req *http.Request) {
		h := req.Header.Clone()

		for k, v := range h {
			//接收客户端 request，并将 request 中带的 header 写入 response header
			res.Header().Set(k, strings.Join(v, ","))
		}

		//读取当前系统的环境变量中的 VERSION 配置，并写入 response header
		res.Header().Set("VERSION", os.Getenv("VERSION"))
		//Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
		res.Write([]byte("hello,index"))

		log.Printf("客户端IP:%s,HTTT返回码:%d", req.Host, http.StatusOK)
	}

	//f实现了http.Handler接口
	http.Handle("/", f)
	//当访问 localhost/healthz 时，应返回 200
	http.HandleFunc("/healthz", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("200"))
	})
	http.ListenAndServe(":80", nil)

}
