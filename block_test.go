package godashboard

import (
	"bytes"
	"context"
	"html/template"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/lucasew/gocfg"
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

	// Label without background_color
	section3 := MockSectionProvider{
		"label": "L",
	}
	_, err = SectionAsRenderBlock(section3)
	if err == nil {
		t.Error("Expected error for label without background_color")
	}
	if !strings.Contains(err.Error(), "background_color") {
		t.Errorf("error %q should mention background_color", err.Error())
	}

	// Label with blank background_color (key present, empty value) → same broken fill:;
	section3b := MockSectionProvider{
		"label":            "L",
		"background_color": "   ",
	}
	_, err = SectionAsRenderBlock(section3b)
	if err == nil {
		t.Error("Expected error for blank background_color")
	}
	if !strings.Contains(err.Error(), "background_color") {
		t.Errorf("error %q should mention background_color", err.Error())
	}

	// background_image key present but empty
	section3c := MockSectionProvider{
		"background_image": "  ",
	}
	_, err = SectionAsRenderBlock(section3c)
	if err == nil {
		t.Error("Expected error for empty background_image")
	}
	if !strings.Contains(err.Error(), "background_image") {
		t.Errorf("error %q should mention background_image", err.Error())
	}

	// Neither background_image nor label
	section4 := MockSectionProvider{
		"size_x": "2",
	}
	_, err = SectionAsRenderBlock(section4)
	if err == nil {
		t.Error("Expected error for section without block type keys")
	}
	if !strings.Contains(err.Error(), "background_image or label") {
		t.Errorf("error %q should explain required keys", err.Error())
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

// exampleTempData matches the fields used by config.ini.example [temp] labels.
type exampleTempData struct {
	Temperatures []struct {
		Temperature float64
	}
}

func TestConfigExample_TempLabelHandlesEmptySensors(t *testing.T) {
	t.Parallel()

	f, err := os.Open("config.ini.example")
	if err != nil {
		t.Fatalf("open config.ini.example: %v", err)
	}
	t.Cleanup(func() {
		if cerr := f.Close(); cerr != nil {
			t.Errorf("close config.ini.example: %v", cerr)
		}
	})

	cfg := gocfg.NewConfig()
	if err := cfg.InjestReader(f); err != nil {
		t.Fatalf("parse config.ini.example: %v", err)
	}
	label := cfg.RawGet("temp", "label")
	if label == "" {
		t.Fatal("temp.label missing from config.ini.example")
	}

	tpl, err := template.New("temp").Parse(label)
	if err != nil {
		t.Fatalf("parse temp label template: %v", err)
	}

	// Empty sensor list (VMs, containers, hosts without hwmon): must not error.
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, exampleTempData{}); err != nil {
		t.Fatalf("temp label with empty sensors: %v (label=%q)", err, label)
	}
	if got := buf.String(); got != "n/a" {
		t.Fatalf("temp label empty sensors = %q, want n/a", got)
	}

	// One sensor still formats the reading.
	buf.Reset()
	data := exampleTempData{Temperatures: []struct{ Temperature float64 }{{Temperature: 42.5}}}
	if err := tpl.Execute(&buf, data); err != nil {
		t.Fatalf("temp label with sensor: %v", err)
	}
	if got := buf.String(); got != "42.5 °C" {
		t.Fatalf("temp label with sensor = %q, want %q", got, "42.5 °C")
	}
}

func TestIndexTemperaturesZero_EmptySliceErrors(t *testing.T) {
	t.Parallel()

	// Documents why config.ini.example must guard before index.
	tpl, err := template.New("unsafe").Parse(`{{with index .Temperatures 0}}{{.Temperature}}{{end}}`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	err = tpl.Execute(io.Discard, exampleTempData{})
	if err == nil {
		t.Fatal("expected index on empty Temperatures to error")
	}
}
