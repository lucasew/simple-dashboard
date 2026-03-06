package godashboard

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lucasew/gocfg"
)

type mockSectionProvider map[string]string

func (m mockSectionProvider) RawGet(key string) string {
	return m[key]
}

func (m mockSectionProvider) RawHasKey(key string) bool {
	_, ok := m[key]
	return ok
}

func (m mockSectionProvider) RawSet(key, value string) bool {
	m[key] = value
	return true
}

func TestNewGoDashboard_Order(t *testing.T) {
	// Setup config with keys that would sort differently
	// "z_last", "a_first", "m_middle"

	cfg := make(gocfg.Config)

	cfg["z_last"] = mockSectionProvider{
		"label":            "Z_Last",
		"background_color": "red",
	}
	cfg["a_first"] = mockSectionProvider{
		"label":            "A_First",
		"background_color": "blue",
	}
	cfg["m_middle"] = mockSectionProvider{
		"label":            "M_Middle",
		"background_color": "green",
	}

	handler := NewGoDashboard(cfg)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	body := w.Body.String()

	// Check if "A_First" comes before "M_Middle" which comes before "Z_Last"
	idxA := strings.Index(body, "A_First")
	idxM := strings.Index(body, "M_Middle")
	idxZ := strings.Index(body, "Z_Last")

	if idxA == -1 || idxM == -1 || idxZ == -1 {
		t.Fatal("Missing sections in output")
	}

	if idxA >= idxM || idxM >= idxZ {
		t.Errorf("Expected order A < M < Z, got indices: A=%d, M=%d, Z=%d", idxA, idxM, idxZ)
	}
}
