package main

import (
    "github.com/lucasew/godashboard"
    "os"
    "github.com/lucasew/gocfg"
    "net/http"
    "flag"
    "fmt"
    "log"
    "github.com/davecgh/go-spew/spew"

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
    flag.IntVar(&PORT, "p", 8080, "Port to listen for connections")
    flag.Parse()
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
    err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), dashboard)
    if err != nil {
        panic(err)
    }

}
