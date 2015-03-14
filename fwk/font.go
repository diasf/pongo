package fwk

import (
	"image"
	"image/draw"
	"io/ioutil"
	"log"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
)

type Font struct {
	font   *truetype.Font
	config FontConfig
}

type FontConfig struct {
	Dpi     float64
	Size    float64
	Spacing float64
}

func NewFont(fontfile string, param ...FontParam) (font *Font, err error) {
	var trueTypeFont *truetype.Font
	if trueTypeFont, err = loadFont(fontfile); err == nil {
		font = &Font{font: trueTypeFont, config: newFontConfig(param...)}
	}
	return
}

func loadFont(fontfile string) (font *truetype.Font, err error) {
	var fontBytes []byte
	if fontBytes, err = ioutil.ReadFile(fontfile); err == nil {
		font, err = freetype.ParseFont(fontBytes)
	}
	return
}

type FontParam func(config *FontConfig)

func newFontConfig(params ...FontParam) FontConfig {
	config := FontConfig{
		Dpi:     72,
		Size:    12,
		Spacing: 1.2,
	}
	for _, p := range params {
		p(&config)
	}
	return config
}

func FontSize(size float64) FontParam {
	return func(config *FontConfig) {
		config.Size = size
	}
}

func FontDpi(dpi float64) FontParam {
	return func(config *FontConfig) {
		config.Dpi = dpi
	}
}

func FontSpacing(spacing float64) FontParam {
	return func(config *FontConfig) {
		config.Spacing = spacing
	}
}

func (f *Font) PrintLines(text []string) draw.Image {
	// Initialize the context.
	fg, bg := image.Black, image.White

	// context
	width, height := f.calcDimensions(text)
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	// draw background
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Over)

	// setup the freetype context
	c := freetype.NewContext()
	c.SetDPI(f.config.Dpi)
	c.SetFont(f.font)
	c.SetFontSize(f.config.Size)
	c.SetSrc(fg)
	c.SetDst(rgba)
	c.SetClip(rgba.Bounds())

	cspacing := c.PointToFix32(f.config.Size * f.config.Spacing)

	// Draw the text.
	pt := freetype.Pt(0, int(f.config.Size))
	for _, s := range text {
		if _, err := c.DrawString(s, pt); err != nil {
			log.Println(err)
			return nil
		}
		pt.Y += cspacing
	}

	return rgba
}

func (f *Font) calcDimensions(text []string) (width, height int) {
	size := int(f.config.Size)
	for _, line := range text {
		lineWidth := 0
		lineHeight := 0
		for _, ch := range line {
			ind := f.font.Index(ch)
			lineWidth += int(f.font.HMetric(int32(size), ind).AdvanceWidth)
			if hm := int(f.font.VMetric(int32(size), ind).AdvanceHeight); hm > lineHeight {
				lineHeight = hm
			}
		}
		if lineWidth > width {
			width = lineWidth
		}
		height += lineHeight + int(f.config.Size*f.config.Spacing-f.config.Size)
	}
	return
}
