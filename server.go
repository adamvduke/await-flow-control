package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	serverCert := flag.String("cert", "cert.pem", "The path to the certificate")
	serverKey := flag.String("key", "key.pem", "The path to the certificate key")
	readBody := flag.Bool("read-body", false, "Wether the server should read the incoming request bodies or not")

	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		statusCode := 200

		body := "{\"message\": \"OK\"}"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		fmt.Fprintf(w, body)
		if *readBody {
			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}
			fmt.Println(r.Proto, statusCode, time.Now().Format(time.RubyDate), string(reqBody))
		} else {
			fmt.Println(r.Proto, statusCode)
		}
	})

	log.Fatal(http.ListenAndServeTLS(":5000", *serverCert, *serverKey, nil))
}
