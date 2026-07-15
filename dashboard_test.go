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
