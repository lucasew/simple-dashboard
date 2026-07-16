package godashboard

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lucasew/gocfg"
)

// errBlock fails RenderBlock with a fixed error (for ServeHTTP error-path tests).
type errBlock struct {
	err error
}

func (e errBlock) SizeX() int { return 1 }
func (e errBlock) SizeY() int { return 1 }
func (e errBlock) RenderBlock(*RequestContext, io.Writer) error {
	return e.err
}

func TestParseReloadTimeoutMs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		set  bool
		want int
	}{
		{name: "missing uses default", set: false, want: defaultReloadTimeoutMs},
		{name: "valid custom", set: true, raw: "2500", want: 2500},
		{name: "zero falls back", set: true, raw: "0", want: defaultReloadTimeoutMs},
		{name: "negative falls back", set: true, raw: "-5", want: defaultReloadTimeoutMs},
		{name: "non-numeric falls back", set: true, raw: "fast", want: defaultReloadTimeoutMs},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := gocfg.NewConfig()
			if tt.set {
				cfg.RawSet("", "reload_timeout", tt.raw)
			}
			got := parseReloadTimeoutMs(cfg)
			if got != tt.want {
				t.Fatalf("parseReloadTimeoutMs() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestServeHTTP_UsesConfiguredReloadTimeout(t *testing.T) {
	t.Parallel()

	cfg := gocfg.NewConfig()
	cfg.RawSet("", "reload_timeout", "3500")
	cfg.RawSet("cpu", "label", "ok")
	cfg.RawSet("cpu", "background_color", "red")

	h := NewGoDashboard(cfg)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, "setTimeout(() => window.location.reload(true), 3500)") {
		t.Fatalf("response missing configured reload timeout; body=%s", body)
	}
	if strings.Contains(body, "setTimeout(() => window.location.reload(true), 1000)") {
		t.Fatalf("response still uses hardcoded 1000ms reload; body=%s", body)
	}
}

func TestServeHTTP_SetsHTMLContentType(t *testing.T) {
	t.Parallel()

	h := NewGoDashboardFromBlocks()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	ct := rec.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("Content-Type = %q, want text/html", ct)
	}
	if !strings.Contains(ct, "charset=utf-8") {
		t.Fatalf("Content-Type = %q, want charset=utf-8", ct)
	}
}

func TestNewGoDashboardFromBlocks_DefaultReloadTimeout(t *testing.T) {
	t.Parallel()

	h := NewGoDashboardFromBlocks()
	d, ok := h.(*GoDashboard)
	if !ok {
		t.Fatalf("expected *GoDashboard, got %T", h)
	}
	if d.reloadTimeoutMs != defaultReloadTimeoutMs {
		t.Fatalf("reloadTimeoutMs = %d, want %d", d.reloadTimeoutMs, defaultReloadTimeoutMs)
	}
}

func TestServeHTTP_EscapesErrorInClientAlert(t *testing.T) {
	t.Parallel()

	// Crafted payload would break out of a raw alert(`...`) if not encoded.
	payload := "x`);</script><script>/*xss*/"
	h := NewGoDashboardFromBlocks(errBlock{err: errors.New(payload)})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	body := rec.Body.String()
	if strings.Contains(body, "alert(`") {
		t.Fatalf("error path still uses unescaped template literal; body=%s", body)
	}
	if strings.Contains(body, "</script><script>") {
		t.Fatalf("raw error allowed script/HTML breakout; body=%s", body)
	}
	if !strings.Contains(body, "<script>alert(") {
		t.Fatalf("expected client alert script; body=%s", body)
	}
	// encoding/json HTML-escapes < and >; message text must still appear.
	if !strings.Contains(body, "xss") {
		t.Fatalf("expected encoded error payload in body; body=%s", body)
	}
	if !strings.Contains(body, `\u003c/script\u003e`) {
		t.Fatalf("expected HTML-escaped script close sequence; body=%s", body)
	}
}

func TestWriteClientErrorAlert_JSONEncodesMessage(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	writeClientErrorAlert(rec, errors.New(`he said "hi"`))
	body := rec.Body.String()
	if !strings.HasPrefix(body, "<script>alert(") || !strings.HasSuffix(body, ")</script>") {
		t.Fatalf("unexpected wrapper: %s", body)
	}
	if strings.Contains(body, "alert(`") {
		t.Fatalf("must not use raw template literal: %s", body)
	}
	// Quotes in the message must be JSON-escaped inside the alert argument.
	if !strings.Contains(body, `\"hi\"`) && !strings.Contains(body, `\u0022hi\u0022`) {
		t.Fatalf("expected JSON-escaped quotes in alert arg; body=%s", body)
	}
}

func TestNewGoDashboard_StableSectionOrder(t *testing.T) {
	t.Parallel()

	// Insert sections out of alphabetical order; render order must still be
	// sorted by section name (gocfg.Config is a map with random range order).
	cfg := gocfg.NewConfig()
	cfg.RawSet("zebra", "label", "ZEBRA")
	cfg.RawSet("zebra", "background_color", "black")
	cfg.RawSet("alpha", "label", "ALPHA")
	cfg.RawSet("alpha", "background_color", "white")
	cfg.RawSet("middle", "label", "MIDDLE")
	cfg.RawSet("middle", "background_color", "gray")
	cfg.RawSet("", "reload_timeout", "1000") // global section must not become a block

	h := NewGoDashboard(cfg)
	d, ok := h.(*GoDashboard)
	if !ok {
		t.Fatalf("expected *GoDashboard, got %T", h)
	}
	if len(d.blocks) != 3 {
		t.Fatalf("blocks = %d, want 3 (global section excluded)", len(d.blocks))
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	body := rec.Body.String()

	alphaAt := strings.Index(body, "ALPHA")
	middleAt := strings.Index(body, "MIDDLE")
	zebraAt := strings.Index(body, "ZEBRA")
	if alphaAt < 0 || middleAt < 0 || zebraAt < 0 {
		t.Fatalf("missing expected labels in body=%s", body)
	}
	if !(alphaAt < middleAt && middleAt < zebraAt) {
		t.Fatalf("block order not alphabetical by section name: alpha=%d middle=%d zebra=%d body=%s",
			alphaAt, middleAt, zebraAt, body)
	}
}
