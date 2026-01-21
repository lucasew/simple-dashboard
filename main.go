package godashboard

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/lucasew/gocfg"
)

var PORT int

const htmlBefore = `
<html>
    <head>
        <meta name="viewport" content="width=device-width, minimum-scale=1.0, maximum-scale=1.0, initial-scale=1.0">
    </head>
    <body id="top">
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
    <script>
        document.body.style.zoom = "100%";
        if (!window.location.href.endsWith("#top")) {
            window.location.href = "#top"
        }
    </script>
    </body>
</html>
`

type GoDashboard struct {
	blocks []RenderableBlock
}

func NewGoDashboard(cfg gocfg.Config) http.Handler {
	blocks := []RenderableBlock{}
	for k, v := range cfg {
		if k == "" {
			continue
		}
		block, err := SectionAsRenderBlock(v)
		if err != nil {
			panic(fmt.Errorf("while loading section '%s': %w", k, err))
		}
		blocks = append(blocks, block)
		log.Printf("setting up section %s\n", k)
	}
	return NewGoDashboardFromBlocks(blocks...)
}

func NewGoDashboardFromBlocks(blocks ...RenderableBlock) http.Handler {
	return &GoDashboard{blocks}
}

func (g *GoDashboard) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("got connection from %s", r.RemoteAddr)
	var err error
	reqContext := NewRequestContext(r.Context())

	_, err = fmt.Fprint(w, htmlBefore)
	if err != nil {
		goto handle_err
	}

	for _, block := range g.blocks {
		err = block.RenderBlock(reqContext, w)
		if err != nil {
			goto handle_err
		}
	}

	_, err = fmt.Fprint(w, htmlAfter)
	if err != nil {
		goto handle_err
	}
	fmt.Fprintln(w, "<script>setTimeout(() => window.location.reload(true), 1000)</script>")
	return
handle_err:
	errBytes, _ := json.Marshal(err.Error())
	fmt.Fprintf(w, "<script>alert(%s)</script>", string(errBytes))
}
