package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"golang.org/x/net/http2"
)

type Response struct {
	Message string
}

func main() {
	cert := flag.String("cert", "cert.pem", "The path to the certificate")
	key := flag.String("key", "key.pem", "The path to the certificate key")
	host := flag.String("host", "localhost", "The host to issue requests against")

	flag.Parse()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	certificate, _ := tls.LoadX509KeyPair(*cert, *key)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{certificate},
	}
	if len(certificate.Certificate) > 0 {
		tlsConfig.BuildNameToCertificate()
	}
	transport := &http2.Transport{
		TLSClientConfig: tlsConfig,
	}
	HTTPClient := &http.Client{Transport: transport}

	makeRequests(HTTPClient, host)
	var userInput string
	for {
		fmt.Scanln(&userInput)
		makeRequests(HTTPClient, host)
	}
}

func makeRequests(HTTPClient *http.Client, host *string) {
	for i := 0; i < 2000; i++ {
		go func() {
			payload := []byte("test test test test test test test test test test test test test test test test test test test test test test test test")
			reader := bytes.NewReader(payload)
			url := fmt.Sprintf("https://%s:5000/", *host)
			req, err := http.NewRequest("POST", url, reader)
			if err != nil {
				panic(err)
			}
			req.Header.Set("Content-Type", "text/plain")
			req.ContentLength = int64(len(payload))
			res, err := HTTPClient.Do(req)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()
			response := &Response{}
			decoder := json.NewDecoder(res.Body)
			if err := decoder.Decode(&response); err != nil && err != io.EOF {
				panic(err)
			}
			fmt.Println(res.Proto, response.Message, time.Now().Format(time.RubyDate))
		}()
	}
}
