package godashboard

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/lucasew/gocfg"
)

// defaultReloadTimeoutMs is used when reload_timeout is missing or invalid.
const defaultReloadTimeoutMs = 1000

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
	blocks          []RenderableBlock
	reloadTimeoutMs int
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
	d := &GoDashboard{
		blocks:          blocks,
		reloadTimeoutMs: parseReloadTimeoutMs(cfg),
	}
	return d
}

func parseReloadTimeoutMs(cfg gocfg.Config) int {
	if !cfg.RawHasKey("", "reload_timeout") {
		return defaultReloadTimeoutMs
	}
	raw := cfg.RawGet("", "reload_timeout")
	ms, err := strconv.Atoi(raw)
	if err != nil || ms <= 0 {
		log.Printf("invalid reload_timeout %q, using default %dms", raw, defaultReloadTimeoutMs)
		return defaultReloadTimeoutMs
	}
	return ms
}

func NewGoDashboardFromBlocks(blocks ...RenderableBlock) http.Handler {
	return &GoDashboard{
		blocks:          blocks,
		reloadTimeoutMs: defaultReloadTimeoutMs,
	}
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
	_, err = fmt.Fprintf(w, "<script>setTimeout(() => window.location.reload(true), %d)</script>\n", g.reloadTimeoutMs)
	if err != nil {
		goto handle_err
	}
	return
handle_err:
	// If we can't write to w, there's not much we can do.
	_, _ = fmt.Fprintf(w, "<script>alert(`%s`)</script>", err.Error())
}
