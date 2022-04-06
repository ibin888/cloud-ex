package main

import (
	"k8s.io/klog/v2"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
)

func main() {
	var index http.HandlerFunc
	/*
		接收客户端 request，并将 request 中带的 header 写入 response header
		读取当前系统的环境变量中的 VERSION 配置，并写入 response header
		Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
	*/
	index = func(w http.ResponseWriter, req *http.Request) {
		h := req.Header.Clone()
		for k, v := range h {
			w.Header().Set(k, strings.Join(v, ","))
		}

		w.Header().Set("VERSION", os.Getenv("VERSION"))

		w.Write([]byte("hello,index"))

		klog.Infof("客户端IP:%s,HTTT返回码:%d", getClientIP(req), http.StatusOK)
	}

	healthz := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200"))
		klog.Infof("客户端IP:%s,HTTT返回码:%d", getClientIP(req), http.StatusOK)
	}

	var httpPort string
	if runtime.GOOS == "windows" {
		httpPort = "8080"
	} else {
		httpPort = os.Getenv("httpPort")
	}

	http.Handle("/", index)
	http.HandleFunc("/healthz", healthz)

	klog.Infof("程序运行在%s端口", httpPort)

	http.ListenAndServe(":"+httpPort, nil)
}

//获取客户端IP
func getClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}
