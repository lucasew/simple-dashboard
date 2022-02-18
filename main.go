package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/lucasew/gocfg"
)

var PORT int
var config gocfg.Config

const htmlBefore = `
<html>
    <head>
        <meta name="viewport" content="width=device-width, initial-scale=1">
    </head>
    <body>
    <style>
    body {
        margin: 0;
        padding: 0;
        width: 100vw;
        height: 100vh;
        max-height: 100vh;
    }
    svg {
        display: block;
        margin: 0;
        padding: 0;
        max-height: 20vh;
    }
    .box {
        display: table;
        float: left;
        margin: 0;
        padding: 0;
    }

    .box h1 {
        display: table-cell;
        font-weight: normal;
        font-size: 2rem;
        margin: 0;
        text-align: center;
        vertical-align: middle;
    }
    </style>
`
const htmlAfter = `
    <script>setTimeout(() => window.location.reload(), 1000)</script>
    </body>
</html>
`

var blocks []RenderableBlock

func init() {
    // flag.StringVar(&TEMPLATE_FILE, "t", "index.html", "template file to use for dashboard")
    var configFile string
    flag.StringVar(&configFile, "c", "config.cfg", "Config file with the blocks defined")
    flag.IntVar(&PORT, "p", 8080, "Port to listen for connections")
    flag.Parse()
    f, err := os.Open(configFile)
    if err != nil {
        panic(err)
    }
    config = gocfg.NewConfig()
    err = config.InjestReader(f)
    if err != nil {
        panic(err)
    }
    for k, v := range config {
        if k == "" {
            continue
        }
        block, err := SectionAsRenderBlock(v)
        if err != nil {
            panic(fmt.Errorf("while loading section '%s': %w", k, err))
        }
        blocks = append(blocks, block)
        fmt.Printf("setting up section %s\n", k)
    }
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        var err error
        reqContext := NewRequestContext(r.Context())

        _, err = fmt.Fprint(w, htmlBefore)
        if err != nil { goto handle_err }

        for _, block := range blocks {
            err = block.RenderBlock(reqContext, w)
            if err != nil { goto handle_err }
        }

        _, err = fmt.Fprint(w, htmlAfter)
        if err != nil { goto handle_err }

        return
        handle_err:
        fmt.Fprintf(w, "<script>alert(`%s`)</script>", err.Error())
    })
    err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
    if err != nil {
        panic(err)
    }
}
