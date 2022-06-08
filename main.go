package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"net/http/httputil"
	"fmt"
	"io"
	"bytes"
	"time"
)

var cache = make(map[string]string)

func InterceptResponse(r *http.Response) error {
	log.Printf("%s", "Response received for path " + r.Request.URL.Path)
	body, err := io.ReadAll(r.Body)
	
	if(err != nil) {
		return err
	}

	r.Body.Close()
	cache[r.Request.URL.Path] = string(body)

	buf := bytes.NewBufferString("")
    buf.Write(body)
	r.Body = ioutil.NopCloser(buf)

	r.Header.Add("seedien-cache", "miss")
	
	return nil
}

func reverseProxy(res http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse("http://localhost:3000")

	proxy := httputil.NewSingleHostReverseProxy(url)

	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	proxy.ModifyResponse = InterceptResponse

	proxy.ServeHTTP(res, req)
}

func getMatchingCacheObject(key string) string {
	if value, ok := cache[key]; ok {
		return value
	}
	return ""
}

func handleRequest(res http.ResponseWriter, req *http.Request) {
	start := time.Now()
	log.Printf("Request: %s %s\n", req.Method, req.URL.Path)
	cached := getMatchingCacheObject(req.URL.Path)

	if len(cached) != 0 {
		log.Printf("Cache hit")
		res.Header().Set("seedien-cache", "hit")
		fmt.Fprintf(res, "%s", cached)
		elapsed := time.Since(start)
    	log.Printf("Hit latency %s", elapsed)
	} else {
		log.Printf("Cache miss")
		reverseProxy(res, req)
		elapsed := time.Since(start)
    	log.Printf("Miss latency %s", elapsed)
	}
}

func getListenAddress() string {
	return ":" + "8000"
}

func main() {
	http.HandleFunc("/", handleRequest)
	if err := http.ListenAndServe(getListenAddress(), nil); err != nil {
		panic(err);
	}
}