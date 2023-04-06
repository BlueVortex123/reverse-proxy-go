package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

func main() {
	targetUrl := &url.URL{
		Scheme: "http",
		Host:   "ergast.com",
	}

	// getting 403 Forbidden for using only NewSingleHostReverseProxy.
	// proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	// Switching to more complex definition
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	// director := proxy.Director
	// proxy.Director = func(req *http.Request) {
	// 	// director(req)
	// 	if req.Method == "GET" {
	// 		req.URL.Scheme = targetUrl.Scheme // for targe scheme error handling
	// 		req.URL.Host = targetUrl.Host     // for  http: no Host in request URL error handling
	// 		req.Host = targetUrl.Host

	// 	}
	// 	req.Header.Add("Accept", "application/json")
	// 	fmt.Println(req.URL)
	// }

	// proxy.Transport = &captureTransport{
	// 	Transport: http.DefaultTransport,
	// }

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://ergast.com"+r.RequestURI, http.StatusTemporaryRedirect)
	})

	proxy.Transport = &captureTransport{http.DefaultTransport}

	http.ListenAndServe(":8080", proxy)
}

type captureTransport struct {
	Transport http.RoundTripper
}

func (ct *captureTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = ct.Transport.RoundTrip(req)
	fmt.Println(resp)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	body = bytes.Replace(body, []byte("MRData"), []byte(randomValue()), -1)
	body_data := ioutil.NopCloser(bytes.NewReader(body))
	resp.Body = body_data
	resp.ContentLength = int64(len(body))
	resp.Header.Set("Content-Length", strconv.Itoa(len(body)))
	return resp, nil
}

func randomValue() string {
	rand.Seed(time.Now().Unix())
	randomValues := []string{"foo", "bar", "slug"}
	return randomValues[rand.Intn(len(randomValues))]

}
