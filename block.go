package godashboard

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"strconv"

	"github.com/lucasew/gocfg"
)

type RenderableBlock interface {
	SizeX() int
	SizeY() int
	RenderBlock(*RequestContext, io.Writer) error
}

type BackgroundImageBlock struct {
	sx, sy    int
	image_url *template.Template
}

type LabelBlock struct {
	sx, sy           int
	label            *template.Template
	background_color *template.Template
}

func SectionAsRenderBlock(section gocfg.SectionProvider) (RenderableBlock, error) {
	size_x := 1
	size_y := 1
	if section.RawHasKey("size_x") {
		r, err := strconv.Atoi(section.RawGet("size_x"))
		if err != nil {
			return nil, fmt.Errorf("while getting size_x: %w", err)
		}
		size_x = r
	}
	if section.RawHasKey("size_y") {
		r, err := strconv.Atoi(section.RawGet("size_y"))
		if err != nil {
			return nil, fmt.Errorf("while getting size_y: %w", err)
		}
		size_y = r
	}
	if section.RawHasKey("background_image") && section.RawHasKey("background_color") {
		return nil, fmt.Errorf("a section cant have both a background_image and a background_color")
	}
	if section.RawHasKey("background_image") && section.RawHasKey("label") {
		return nil, fmt.Errorf("a section cant have both a background_image and a label")
	}
	if section.RawHasKey("background_image") {
		tpl, err := template.New("background_image").Parse(section.RawGet("background_image"))
		if err != nil {
			return nil, fmt.Errorf("invalid background_image template: %w", err)
		}
		return BackgroundImageBlock{
			sx: size_x, sy: size_y,
			image_url: tpl,
		}, nil
	}
	if section.RawHasKey("label") {
		tpl_label, err := template.New("label").Parse(section.RawGet("label"))
		if err != nil {
			return nil, fmt.Errorf("invalid label template: %w", err)
		}
		tpl_color, err := template.New("bg_color").Parse(section.RawGet("background_color"))
		if err != nil {
			return nil, fmt.Errorf("invalid label template: %w", err)
		}
		return LabelBlock{
			sx: size_x, sy: size_y,
			label:            tpl_label,
			background_color: tpl_color,
		}, nil
	}
	return nil, errors.New("unmatched block type")
}

func (b BackgroundImageBlock) SizeX() int {
	return b.sx
}

func (l LabelBlock) SizeX() int {
	return l.sx
}

func (b BackgroundImageBlock) SizeY() int {
	return b.sy
}

func (l LabelBlock) SizeY() int {
	return l.sy
}

func (b BackgroundImageBlock) RenderBlock(ctx *RequestContext, w io.Writer) error {
	fmt.Fprintf(w, `<img class="box" width="%d" height="%d" src="`, b.sx*ctx.sizeBaseline, b.sy*ctx.sizeBaseline)
	err := b.image_url.Execute(w, ctx)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, `">`)
	return nil
}

func (l LabelBlock) RenderBlock(ctx *RequestContext, w io.Writer) error {
	sx := l.sx * ctx.sizeBaseline
	sy := l.sy * ctx.sizeBaseline
	fmt.Fprintf(w, `<svg class="box" viewBox="0 0 %d %d" width="%d" height="%d">`, sx, sy, sx, sy)
	fmt.Fprint(w, `<rect style="fill:`)
	err := l.background_color.Execute(w, ctx)
	if err != nil {
		return err
	}
	fmt.Fprint(w, `;" x="0%" y="0%" height="100%" width="100%" />`)
	fmt.Fprint(w, `<text x="50%" y="50%" font-size="2rem" text-anchor="middle" dominant-baseline="middle">`)
	err = l.label.Execute(w, ctx)
	if err != nil {
		return err
	}
	fmt.Fprint(w, `</text></svg>`)
	return nil
}
