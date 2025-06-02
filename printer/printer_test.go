package printer

import (
	"fmt"
	image2 "image"
	"image/draw"
	"image/jpeg"
	"os"
	"testing"

	bgr2 "github.com/sipkg/golang-win32-printer/image/bgr"
	"github.com/sipkg/golang-win32-printer/win32"
)

type imagePrinter struct{}

func (ip *imagePrinter) Print(p *Printer) {
	file, err := os.Open("C:\\Users\\wangjun\\Desktop\\pdf\\sekiro.png")
	fmt.Print(err)
	image, err := jpeg.Decode(file)
	fmt.Print(err)
	bgr := bgr2.NewBGRImage(image.Bounds())
	draw.Draw(bgr, image.Bounds(), image, image2.Point{}, draw.Src)
	src := bgr2.ReverseDIB(bgr.Pix, image.Bounds().Dx(), image.Bounds().Dy(), 24)
	left, _ := p.GetMarginLeft()
	top, _ := p.GetMarginLeft()
	width, _ := p.GetWidthPixel()
	height, _ := p.GetHeightPixel()
	win32.DrawDIImage(p.hdc, left, top, width, height, 0, 0, int32(image.Bounds().Dx()), int32(image.Bounds().Dy()), src)
}

func TestPrinter(t *testing.T) {
	printer := Printer{}
	err := printer.InitPrinter("Microsoft Print to PDF")
	if err != nil {
		t.Errorf("%v", err)
	}
	printer.Print(&imagePrinter{})
}
