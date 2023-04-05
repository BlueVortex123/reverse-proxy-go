package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func main() {
	targetUrl, err := url.Parse("https://jsonplaceholder.typicode.com")
	if err != nil {
		fmt.Println(err)
	}

	// getting 403 Forbidden for using only NewSingleHostReverseProxy.
	// proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	// Switching to more complex definition
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "https",
		Host:   "jsonplaceholder.typicode.com",
	})

	// director := proxy.Director
	proxy.Director = func(req *http.Request) {
		// director(req)
		if req.Method == "GET" {
			req.URL.Scheme = targetUrl.Scheme // for targe scheme error handling
			req.URL.Host = targetUrl.Host     // for  http: no Host in request URL error handling
			req.Host = targetUrl.Host

		}
		req.Header.Add("Accept", "application/json")
		fmt.Println(req.URL)
	}

	proxy.Transport = &captureTransport{
		Transport: http.DefaultTransport,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://jsonplaceholder.typicode.com"+r.RequestURI, http.StatusTemporaryRedirect)
	})

	http.ListenAndServe(":8080", proxy)
}

type captureTransport struct {
	Transport http.RoundTripper
}

func (ct *captureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := ct.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		return nil, err
	}

	resp.Body.Close()

	// if resp.Header.Get("Content-Type") == "application/json; charset=utf-8" && resp.ContentLength != 0 {

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		// log.Printf("error decoding: %v", err)
		// if e, ok := err.(*json.SyntaxError); ok {
		// 	log.Printf("syntax error of byte offset %d", e.Offset)
		// }
		// log.Printf("Response: %q", body)
		return nil, err
	}

	data["title"] = "test"
	fmt.Printf("%v\n", data)

	body, err = json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// ctx := req.Context()
	// ctx = context.WithValue(ctx, "response_body", body)
	// req = req.WithContext(ctx)

	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	// }
	return resp, nil
}

func randomValue() string {
	rand.Seed(time.Now().Unix())
	randomValues := []string{"foo", "bar", "slug"}
	return randomValues[rand.Intn(len(randomValues))]

}
