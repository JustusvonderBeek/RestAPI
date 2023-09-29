package main

import (
	"flag"
	"io"
	"log"
	"os"
)

func main() {
	setupLogger("api.log")

	addr := flag.String("a", "0.0.0.0", "Listen address")
	port := flag.String("p", "50002", "Listen port")
	overwrite := flag.Bool("e", false, "Overwrite the existing vocabulary")
	client := flag.Bool("c", false, "If set start as client and make request")
	flag.Parse()

	configuration := Configuration{
		IP_Address:  *addr,
		Listen_Port: *port,
		Overwrite:   *overwrite,
		Client:      *client,
	}

	// Starting the main server and waiting for request
	// net.Listen()
	if configuration.Client {
		startingClient(configuration)
	} else {
		startingServer(configuration)
	}
}

func setupLogger(logfile string) {
	logFile, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}
