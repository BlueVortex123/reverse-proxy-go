package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	targetUrl, err := url.Parse("https://jsonplaceholder.typicode.com")
	if err != nil {
		fmt.Println(err)
	}

	// getting 403 Forbidden for using only NewSingleHostReverseProxy.
	// proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	// switching to more complet definition
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "https",
		Host:   "jsonplaceholder.typicode.com",
	})

	proxy.Director = func(req *http.Request) {
		if req.Method == "GET" {
			req.URL.Scheme = targetUrl.Scheme // for targe scheme error handling
			req.URL.Host = targetUrl.Host     // for  http: no Host in request URL error handling
			targetUrl.RawQuery = ""
			req.URL.Path = "" + req.URL.Path
			req.Host = targetUrl.Host

		}
		fmt.Println(req.URL)
	}

	http.ListenAndServe(":8080", proxy)

}
