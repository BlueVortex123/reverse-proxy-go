package main

import (
	"fmt"
	"html"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Test route, %q\n", html.EscapeString(r.URL.Path))
	})

	http.ListenAndServe(":8080", nil)

}
