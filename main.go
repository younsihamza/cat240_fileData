package main

import (
	"cat240/sender"
	"cat240/utils"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":"+os.Getenv("PROMETHUES_PORT"), nil)
	go sender.Sender()   // websockets server
	go utils.ReadData()  // read data from pcap file
	go utils.ParseData() // parse data and send to websockets server
	select {}

}
