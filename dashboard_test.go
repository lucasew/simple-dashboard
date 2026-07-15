package godashboard

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lucasew/gocfg"
)

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
