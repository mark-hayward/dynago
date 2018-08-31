package main

import (
	"fmt"
	"flag"
	"net/http"
)

var (
	HTTPAddr = flag.String("http", "0.0.0.0:8001", "Address to listen for HTTP requests on")
)

func main() {

	fmt.Println("Registering the collector")

	http.HandleFunc("/process", Collector)

//		s3Configuration *S3Config, coldef []ColumnDefinition, dyndb *DynamoDBConfig)

	http.Handle("/", http.FileServer(http.Dir("./web")))
	// Start the HTTP server!
	fmt.Println("HTTP server listening on", *HTTPAddr)
	if err := http.ListenAndServe(*HTTPAddr, nil); err != nil {
		fmt.Println(err.Error())
	}


}

