package main

import (
	"flag"
	"log"
	"time"

	"github.com/OkaniYoshiii/sqlite-go/internal/api"
)

var address = flag.String("address", "127.0.0.1:8000", "Specifies the TCP address for the server to listen on, in the form “host:port”. ")
var readTimeout = flag.Int("readtimeout", 10000, "The maximum duration in milliseconds for reading the entire request, including the body.")
var readHeaderTimeout = flag.Int("readheadertimeout", 2000, "The maximum duration in milliseconds for reading the headers.")
var writeTimeout = flag.Int("writetimeout", 3000, "The maximum duration in milliseconds before timing out writes of the response.")
var idleTimeout = time.Millisecond * 100

func main() {
	flag.Parse()

	if err := api.Run(*address, *readTimeout, *readHeaderTimeout, *writeTimeout, int(idleTimeout)); err != nil {
		log.Fatal(err)
	}
}
