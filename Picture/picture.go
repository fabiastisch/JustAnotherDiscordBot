package Picture

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"
)

type Picture struct {
	img *image.RGBA
}

func New(sizeX, sizeY int) Picture {
	return Picture{
		image.NewRGBA(image.Rectangle{
			Min: image.Point{},
			Max: image.Point{X: sizeX, Y: sizeY},
		}),
	}
}

func (p Picture) Background(color color.RGBA) {
	draw.Draw(p.img, p.img.Bounds(), &image.Uniform{C: color}, image.Point{}, draw.Src)
}

func (p Picture) DrawImage(image2 image.Image, posX, posY int) {
	draw.Draw(p.img, p.img.Bounds(), image2, image.Point{X: -posX, Y: -posY}, draw.Src)
}
func (p Picture) DrawImageBottomRight(image2 image.Image) {
	bounds := p.img.Bounds()

	draw.Draw(p.img, bounds, image2, bounds.Size().Sub(image2.Bounds().Size()).Mul(-1), draw.Over)
}
func (p Picture) DrawImageCenter(image2 image.Image) {
	bounds := p.img.Bounds()
	center := bounds.Size().Div(2)
	center2 := image2.Bounds().Size().Div(2)
	pos := center.Sub(center2)

	draw.Draw(p.img, bounds, image2, pos.Mul(-1), draw.Over)
}

func (p Picture) DrawImageBottomCenter(image2 image.Image) {
	bounds := p.img.Bounds()
	bottomCenter := bounds.Size().Div(2)
	bottomCenter.Y = bounds.Size().Y
	pos := bottomCenter.Sub(image.Point{
		X: image2.Bounds().Size().Div(2).X,
		Y: image2.Bounds().Size().Y,
	})

	draw.Draw(p.img, bounds, image2, pos.Mul(-1), draw.Over)
}

func (p Picture) AddLabelCenterHorizontal(label string, y int, color color.Color) {
	center := p.img.Bounds().Size().Div(2)
	point := fixed.Point26_6{
		X: fixed.I(center.X),
		Y: fixed.I(y),
	}

	f, err := truetype.Parse(gobold.TTF)
	if err != nil {
		log.Panic(err)
	}

	d := &font.Drawer{
		Dst: p.img,
		Src: image.NewUniform(color),
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    28, //TODO: check size
			DPI:     72,
			Hinting: 0,
		}),
		Dot: point,
	}

	d.Dot = fixed.Point26_6{
		X: fixed.I(center.X) - d.MeasureString(label)/2,
		Y: fixed.I(y) + d.Face.Metrics().Height,
	}

	d.DrawString(label)

}

func (p Picture) AddLabel(x, y int, label string) {
	col := color.Black
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	d := &font.Drawer{
		Dst:  p.img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func (p Picture) GetImage() image.Image {
	return p.img
}

func (p Picture) ToReader() *os.File {
	reader, writer, err := os.Pipe()
	log.Println("ToReader")
	if err != nil {
		log.Panic(err)
	}
	err = png.Encode(writer, p.img)
	log.Println("Encoded")

	if err != nil {
		writer.Close()
		log.Panicln(err)
	}
	log.Println("Encodeasd")

	if err = writer.Close(); err != nil {
		log.Panicln(err)
	}
	log.Println("x")

	return reader
}
func GetImageFromURL(url string) image.Image {

	res, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}
	defer res.Body.Close()
	img, _, err := image.Decode(res.Body)
	if err != nil {
		log.Println(err)
	}
	return img
}
