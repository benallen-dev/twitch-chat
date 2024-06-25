package main

import (
	"log"
	"fmt"
	"net/http"
)

// HTTP receiver for token, unused right now
func receiver () {

	// Just use the default mux
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, "%+v", map[string]string{
			"path": r.URL.Path,
			"query": r.URL.RawQuery,
			"method": r.Method,
			"request_uri": r.RequestURI,
		})

		fmt.Println("Got a request")
		fmt.Println(r.RequestURI)
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
