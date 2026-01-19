package godashboard

import (
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

type ErrorBlock struct {
	Msg string
}

func (e ErrorBlock) SizeX() int { return 1 }
func (e ErrorBlock) SizeY() int { return 1 }
func (e ErrorBlock) RenderBlock(ctx *RequestContext, w io.Writer) error {
	return errors.New(e.Msg)
}

func TestXSSVulnerability(t *testing.T) {
	maliciousMsg := "`); alert('XSS');//"
	block := ErrorBlock{Msg: maliciousMsg}
	dashboard := NewGoDashboardFromBlocks(block)

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	dashboard.ServeHTTP(w, req)

	resp := w.Result()
	bodyBytes, _ := io.ReadAll(resp.Body)
	body := string(bodyBytes)

	// Expected escaped output via json.Marshal
	// The string "`); alert('XSS');//" marshals to ""`); alert('XSS');//""
	// Note: json.Marshal wraps the string in double quotes.
	escaped := "\"`); alert('XSS');//\""
	expectedSafe := "<script>alert(" + escaped + ")</script>"

	if strings.Contains(body, expectedSafe) {
		t.Log("Fix verified: Payload is safely escaped.")
	} else {
		t.Errorf("Fix FAILED. Body: %s\nExpected to contain: %s", body, expectedSafe)
	}
}
