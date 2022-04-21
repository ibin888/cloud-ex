package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"httptest/metrics"
	"k8s.io/klog/v2"
	"math/rand"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

func main() {
	var index http.HandlerFunc
	/*
		接收客户端 request，并将 request 中带的 header 写入 response header
		读取当前系统的环境变量中的 VERSION 配置，并写入 response header
		Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
	*/
	index = func(w http.ResponseWriter, req *http.Request) {
		timer := metrics.NewTimer()
		defer timer.ObserveTotal()
		delay := randInt(0, 2000)
		time.Sleep(time.Millisecond * time.Duration(delay))
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

	metrics.Register()
	http.Handle("/", index)
	http.HandleFunc("/healthz", healthz)
	//添加metric,promethues默认放了一些指标
	http.Handle("/metrics", promhttp.Handler())
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

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
