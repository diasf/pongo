package tex

import (
	"errors"
	"image"
	"image/png"
	"io"

	"golang.org/x/mobile/gl"
)

type Texture struct {
	Id     gl.Texture
	Levels map[int]TexData
	target gl.Enum
}

type TexData interface {
	Bin() ([]byte, error)
	Format() gl.Enum
	Width() int
	Height() int
}

func GenTexture() *Texture {
	return &Texture{
		Id:     gl.GenTexture(),
		Levels: map[int]TexData{},
		target: gl.TEXTURE_2D,
	}
}

func (t *Texture) AddTexData(level int, data TexData) {
	t.Levels[level] = data
}

func (t *Texture) Upload() error {
	gl.BindTexture(t.target, t.Id)
	for lvl, data := range t.Levels {
		bin, err := data.Bin()
		if err != nil {
			return err
		}
		gl.TexImage2D(t.target, lvl, data.Width(), data.Height(), data.Format(), gl.UNSIGNED_BYTE, bin)
	}
	return nil
}

func NewTexDataFromPNG(r io.Reader) TexData {
	return &texDataPNG{reader: r}
}

type texDataPNG struct {
	reader io.Reader
	format gl.Enum
	width  int
	height int
}

func (t *texDataPNG) Format() gl.Enum {
	return t.format
}

func (t *texDataPNG) Width() int {
	return t.width
}

func (t *texDataPNG) Height() int {
	return t.height
}

func (t *texDataPNG) Bin() (data []byte, err error) {
	img, err := png.Decode(t.reader)
	if err != nil {
		return
	}

	rgbaImg, ok := img.(*image.NRGBA)
	if !ok {
		return nil, errors.New("Texture must be an NRGBA image")
	}

	t.format = gl.RGBA
	t.width, t.height = img.Bounds().Dx(), img.Bounds().Dy()

	// flip horizontally
	lineLen := t.width * 4
	dataSize := lineLen * t.height

	data = make([]byte, dataSize)
	dest := dataSize - lineLen
	for src := 0; src < len(rgbaImg.Pix); src += rgbaImg.Stride {
		copy(data[dest:dest+lineLen], rgbaImg.Pix[src:src+rgbaImg.Stride])
		dest -= lineLen
	}
	return
}
