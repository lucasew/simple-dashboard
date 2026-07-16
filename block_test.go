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

func renderBlockHelper(t *testing.T, block RenderableBlock) string {
	t.Helper()
	ctx := NewRequestContext(context.Background())
	var buf bytes.Buffer
	err := block.RenderBlock(ctx, &buf)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	return buf.String()
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
	output := renderBlockHelper(t, block)
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
	output := renderBlockHelper(t, block)
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

func TestSectionAsRenderBlock_NonPositiveSizes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		keys map[string]string
		want string
	}{
		{
			name: "zero size_x",
			keys: map[string]string{
				"label":            "L",
				"background_color": "red",
				"size_x":           "0",
			},
			want: "size_x must be positive",
		},
		{
			name: "negative size_y",
			keys: map[string]string{
				"label":            "L",
				"background_color": "red",
				"size_y":           "-2",
			},
			want: "size_y must be positive",
		},
		{
			name: "zero size_x on background image",
			keys: map[string]string{
				"background_image": "http://example.com/i.png",
				"size_x":           "0",
				"size_y":           "1",
			},
			want: "size_x must be positive",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			section := MockSectionProvider(tt.keys)
			_, err := SectionAsRenderBlock(section)
			if err == nil {
				t.Fatal("expected error for non-positive size")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error %q does not contain %q", err.Error(), tt.want)
			}
		})
	}
}
