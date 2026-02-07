package godashboard

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// MockSectionProvider implements gocfg.SectionProvider interface based on usage in block.go
type MockSectionProvider map[string]string

func (m MockSectionProvider) RawGet(key string) string {
	return m[key]
}

func (m MockSectionProvider) RawHasKey(key string) bool {
	_, ok := m[key]
	return ok
}

func (m MockSectionProvider) RawSet(key, value string) bool {
	m[key] = value
	return true
}

func TestSectionAsRenderBlock_BackgroundImage(t *testing.T) {
	section := MockSectionProvider{
		"background_image": "http://example.com/image.png",
		"size_x":           "2",
		"size_y":           "3",
	}

	block, err := SectionAsRenderBlock(section)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if block.SizeX() != 2 {
		t.Errorf("Expected SizeX 2, got %d", block.SizeX())
	}
	if block.SizeY() != 3 {
		t.Errorf("Expected SizeY 3, got %d", block.SizeY())
	}

	// Test Render
	ctx := NewRequestContext(context.Background())
	var buf bytes.Buffer
	err = block.RenderBlock(ctx, &buf)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	output := buf.String()
	expectedWidth := "200"
	if !strings.Contains(output, `width="`+expectedWidth+`"`) {
		t.Errorf("Expected width %s in output, got %s", expectedWidth, output)
	}
	if !strings.Contains(output, `src="http://example.com/image.png"`) {
		t.Errorf("Expected src in output, got %s", output)
	}
}

func TestSectionAsRenderBlock_Label(t *testing.T) {
	section := MockSectionProvider{
		"label":            "My Label",
		"background_color": "red",
	}

	block, err := SectionAsRenderBlock(section)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if block.SizeX() != 1 { // Default
		t.Errorf("Expected SizeX 1, got %d", block.SizeX())
	}

	// Test Render
	ctx := NewRequestContext(context.Background())
	var buf bytes.Buffer
	err = block.RenderBlock(ctx, &buf)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "My Label") {
		t.Errorf("Expected label text, got %s", output)
	}
	if !strings.Contains(output, "fill:red") {
		t.Errorf("Expected fill:red, got %s", output)
	}
}

func TestSectionAsRenderBlock_Errors(t *testing.T) {
	// Conflict
	section := MockSectionProvider{
		"label":            "L",
		"background_image": "I",
	}
	_, err := SectionAsRenderBlock(section)
	if err == nil {
		t.Error("Expected error for conflicting keys")
	}

	// Invalid Size
	section2 := MockSectionProvider{
		"label":            "L",
		"background_color": "red",
		"size_x":           "invalid",
	}
	_, err = SectionAsRenderBlock(section2)
	if err == nil {
		t.Error("Expected error for invalid size_x")
	}
}

func TestSectionAsRenderBlock_InvalidBackgroundColor(t *testing.T) {
	section := MockSectionProvider{
		"label":            "L",
		"background_color": "{{",
	}
	_, err := SectionAsRenderBlock(section)
	if err == nil {
		t.Error("Expected error for invalid background_color template")
	} else if !strings.Contains(err.Error(), "invalid background_color template") {
		t.Errorf("Expected error message to contain 'invalid background_color template', got '%s'", err.Error())
	}
}
