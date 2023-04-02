package main

import (
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

	// switching to more complet definition
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
			targetUrl.RawQuery = ""
			req.URL.Path = "" + req.URL.Path
			req.Host = targetUrl.Host

		}
		fmt.Println(req.URL)
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		// For some reason the "Content-Type" header is missing here, but is present in the header of browser.
		// if resp.Header.Get("Content-Type") == "application/json" {

		// Commenting this conditionals because  the response status in currently 304.
		// if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var data interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return err //Getting unexpected end of JSON input error.
		}

		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return err //Getting EOF
		}

		// }
		// }
		return nil
	}

	http.ListenAndServe(":8080", proxy)

}

func randomValue() string {
	rand.Seed(time.Now().Unix())
	randomValues := []string{"foo", "bar", "slug"}
	return randomValues[rand.Intn(len(randomValues))]

}
