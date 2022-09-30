package Picture

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"
	"strings"
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
func (p Picture) advancedMeasureString(label string, d *font.Drawer, xOffset int) []string {
	labels := strings.Split(label, " ")

	var newLabels []string
	newLabel := ""
	for _, word := range labels {

		measure := d.MeasureString(newLabel + word).Round()
		if measure >= p.img.Bounds().Size().X-xOffset {
			newLabels = append(newLabels, newLabel)
			newLabel = ""
			//newLabel = newLabel + "\n" + word
		} else {
			newLabel = newLabel + " " + word
		}

	}
	newLabels = append(newLabels, newLabel)

	return newLabels
}

func (p Picture) getFace(fontsize float64) font.Face {
	f, err := truetype.Parse(gobold.TTF)
	if err != nil {
		log.Panic(err)
	}
	return truetype.NewFace(f, &truetype.Options{
		Size:    fontsize, //TODO: check size
		DPI:     72,
		Hinting: 0,
	})
}
func (p Picture) AddLabelCenterHorizontalWithOffsetBottom(label string, xOffset, y int, color color.Color, fontsize float64) fixed.Int26_6 {
	return p.AddLabelCenterHorizontalWithOffset(label, xOffset, y-p.getFace(fontsize).Metrics().Height.Ceil(), color, fontsize)
}
func (p Picture) AddLabelCenterHorizontalWithOffset(label string, xOffset, y int, color color.Color, fontsize float64) fixed.Int26_6 {
	center := p.img.Bounds().Size().Div(2)
	point := fixed.Point26_6{
		X: fixed.I(center.X + xOffset/2),
		Y: fixed.I(y),
	}

	d := &font.Drawer{
		Dst:  p.img,
		Src:  image.NewUniform(color),
		Face: p.getFace(fontsize),
		Dot:  point,
	}
	measure := d.MeasureString(label).Round()
	if measure > p.img.Bounds().Size().X-xOffset {

		labels := p.advancedMeasureString(label, d, xOffset)

		for i, s := range labels {
			d.Dot = fixed.Point26_6{
				X: fixed.I(center.X+xOffset/2) - d.MeasureString(s)/2,
				Y: fixed.I(y) + d.Face.Metrics().Height.Mul(fixed.I(i+1)),
			}
			d.DrawString(s)
		}
		return d.Face.Metrics().Height.Mul(fixed.I(len(labels)))
	}

	d.Dot = fixed.Point26_6{
		X: fixed.I(center.X+xOffset/2) - d.MeasureString(label)/2,
		Y: fixed.I(y) + d.Face.Metrics().Height,
	}

	d.DrawString(label)

	return d.Face.Metrics().Height
}
func (p Picture) AddLabelCenterHorizontal(label string, y int, color color.Color, fontSize float64) fixed.Int26_6 {
	return p.AddLabelCenterHorizontalWithOffset(label, 0, y, color, fontSize)
}

func (p Picture) AddLabel(x, y int, fontSize float64, color color.Color, label string) fixed.Int26_6 {
	d := p.GetDrawer(fontSize, color)
	d.Dot = fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}
	d.Dot.Y = fixed.I(y) + d.Face.Metrics().Height
	d.DrawString(label)

	return d.Face.Metrics().Height
}

func (p Picture) GetDrawer(fontSize float64, color color.Color) *font.Drawer {
	d := &font.Drawer{
		Dst:  p.img,
		Src:  image.NewUniform(color),
		Face: p.getFace(fontSize),
		Dot: fixed.Point26_6{
			X: fixed.I(0),
			Y: fixed.I(0),
		},
	}
	return d
}

func (p Picture) GetImage() *image.RGBA {
	return p.img
}
func (p Picture) DrawLine(y int, color color.Color, width int) {
	x := p.img.Bounds().Size().X
	for i := 0; i < x; i++ {
		for w := 0; w < width; w++ {
			p.img.Set(i, y+w, color)
		}
	}
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
		log.Println("Error on GetImageFromURL\nURL:" + url + "\n")
		log.Panic(err)
	}
	defer res.Body.Close()
	img, _, err := image.Decode(res.Body)
	if err != nil {
		log.Println(err)
	}
	return img
}

func GetImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}
