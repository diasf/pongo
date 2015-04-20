package tex

import (
	"image"
	"image/png"
	"io"
	"log"
	"os"

	"golang.org/x/mobile/gl"
)

var currentlyBound gl.Texture

type Texture struct {
	Id       gl.Texture
	Levels   map[int]TexData
	target   gl.Enum
	uploaded bool
}

type TexData interface {
	Bin() ([]byte, error)
	Format() gl.Enum
	Width() int
	Height() int
}

func GenTexture() *Texture {
	return &Texture{
		Id:     gl.CreateTexture(),
		Levels: map[int]TexData{},
		target: gl.TEXTURE_2D,
	}
}

func (t *Texture) AddTexData(level int, data TexData) {
	t.Levels[level] = data
}

func (t *Texture) Upload() error {
	if t.uploaded {
		return nil
	}
	t.Bind()
	for lvl, data := range t.Levels {
		bin, err := data.Bin()
		if err != nil {
			return err
		}
		gl.TexImage2D(t.target, lvl, data.Width(), data.Height(), data.Format(), gl.UNSIGNED_BYTE, bin)
	}
	t.uploaded = true
	return nil
}

func (t *Texture) SetMinFilterNearest() {
	t.Bind()
	gl.TexParameteri(t.target, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
}

func (t *Texture) SetMagFilterNearest() {
	t.Bind()
	gl.TexParameteri(t.target, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
}

func (t *Texture) SetRepeat() {
	t.Bind()
	gl.TexParameteri(t.target, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(t.target, gl.TEXTURE_WRAP_T, gl.REPEAT)
}

func (t *Texture) Bind() {
	if currentlyBound == t.Id {
		return
	}
	gl.BindTexture(t.target, t.Id)
	currentlyBound = t.Id
}

func NewTextureFromPNGFile(fileName string) *Texture {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	texture := GenTexture()
	texture.AddTexData(0, NewTexDataFromPNG(file))
	return texture
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

	pix, stride := getPNGPix(img)

	t.format = gl.RGBA
	t.width, t.height = img.Bounds().Dx(), img.Bounds().Dy()

	// flip horizontally
	lineLen := t.width * 4
	dataSize := lineLen * t.height

	data = make([]byte, dataSize)
	dest := dataSize - lineLen
	for src := 0; src < len(pix); src += stride {
		copy(data[dest:dest+lineLen], pix[src:src+stride])
		dest -= lineLen
	}
	return
}

func getPNGPix(img image.Image) (pix []uint8, stride int) {
	switch pngImg := img.(type) {
	case *image.NRGBA:
		return pngImg.Pix, pngImg.Stride
	case *image.RGBA:
		return pngImg.Pix, pngImg.Stride
	default:
		log.Fatalln("Not supported image format, must be RGBA/NRGBA")
	}
	return
}
