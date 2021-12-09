package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

var TEMPLATE_FILE string
var PORT int
var TEMPLATE *template.Template

func loadTemplate() error {
    f, err := os.Open(TEMPLATE_FILE)
    defer f.Close()
    if err != nil {
        return err
    }
    tmplBytes, err := ioutil.ReadAll(f)
    if err != nil {
        return err
    }
    tmpl, err := template.New("dashboard").Parse(string(tmplBytes))
    if err != nil {
        return err
    }
    TEMPLATE = tmpl
    return nil
}

func init() {
    flag.StringVar(&TEMPLATE_FILE, "t", "index.html", "template file to use for dashboard")
    flag.IntVar(&PORT, "p", 8080, "Port to listen for connections")
    flag.Parse()
    err := loadTemplate()
    if err != nil {
        panic(err)
    }
}
func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        reqContext := NewRequestContext(r.Context())
        err := TEMPLATE.Execute(w, reqContext)
        if err != nil {
            fmt.Printf("error: %s", err.Error())
            return
        }
    })
    err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
    if err != nil {
        panic(err)
    }
}
