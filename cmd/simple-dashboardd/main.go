package main

import (
	"flag"
	"fmt"
	"github.com/lucasew/go-getlistener"
	"github.com/lucasew/gocfg"
	"github.com/lucasew/godashboard"
	"log"
	"net/http"
	"os"
)

const banner = `
         (_)               | |          | |         | |   | |                       | |
      ___ _ _ __ ___  _ __ | | ___    __| | __ _ ___| |__ | |__   ___   __ _ _ __ __| |
     / __| | '_ ' _ \| '_ \| |/ _ \  / _  |/ _  / __| '_ \| '_ \ / _ \ / _  | '__/ _  |
     \__ \ | | | | | | |_) | |  __/ | (_| | (_| \__ \ | | | |_) | (_) | (_| | | | (_| |
     |___/_|_| |_| |_| .__/|_|\___|  \__,_|\__,_|___/_| |_|_.__/ \___/ \__,_|_|  \__,_|
                     | |                                                               
                     |_|                                                               
`

func main() {
	var configFile string
	var PORT int
	flag.StringVar(&configFile, "c", "config.cfg", "Config file with the blocks defined")
	flag.IntVar(&getlistener.PORT, "p", getlistener.PORT, "Port to listen for connections")
	flag.Parse()
	if getlistener.PORT == 0 {
		getlistener.PORT = 8080
	}
	f, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	config := gocfg.NewConfig()
	err = config.InjestReader(f)
	if err != nil {
		panic(err)
	}
	println(banner)
	log.Printf("Listening in port %d", PORT)
	dashboard := godashboard.NewGoDashboard(config)
	ln, err := getlistener.GetListener()
	if err != nil {
		panic(err)
	}
	err = http.Serve(ln, dashboard)
	if err != nil {
		panic(err)
	}

}
