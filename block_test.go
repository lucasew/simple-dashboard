package godashboard

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/lucasew/gocfg"
)

func TestLabelBlock_XSS(t *testing.T) {
	// Setup config with a label containing a script tag via template literal.
	// text/template will output this raw. html/template should escape it.
	cfgContent := `
[section1]
label = {{ "<script>alert(1)</script>" }}
background_color = #fff
size_x = 1
size_y = 1
`
	config := gocfg.NewConfig()
	err := config.InjestReader(strings.NewReader(cfgContent))
	if err != nil {
		t.Fatalf("failed to injest config: %v", err)
	}

	section, ok := config["section1"]
	if !ok {
		t.Fatalf("section1 not found in config")
	}

	block, err := SectionAsRenderBlock(section)
	if err != nil {
		t.Fatalf("failed to create block: %v", err)
	}

	ctx := NewRequestContext(context.Background())
	var buf bytes.Buffer
	err = block.RenderBlock(ctx, &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()

	// Check for unescaped script tag
	if strings.Contains(output, "<script>") {
		t.Errorf("Vulnerability reproduced: Output contains unescaped <script> tag: %s", output)
	}

	// Check for escaped script tag (what we want)
	if !strings.Contains(output, "&lt;script&gt;") {
		t.Logf("Output does not contain escaped script tag yet (expected).")
	}
}
