package win32

import (
	"fmt"
	image2 "image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"testing"

	bgr2 "github.com/clementuu/golang-win32-printer/image/bgr"
	"golang.org/x/image/bmp"
)

func TestCreateDC(t *testing.T) {
	gotDc, err := CreateDC("Microsoft Print to PDF")
	fmt.Print(gotDc)
	fmt.Print(err)
}

func TestDeleteDC(t *testing.T) {
	gotDc, err := CreateDC("Microsoft Print to PDF")
	fmt.Print(gotDc)
	fmt.Print(err)
	er := DeleteDC(gotDc)
	fmt.Print(er)
}

func TestResetDC(t *testing.T) {
	gotDc, err := CreateDC("Microsoft Print to PDF")
	fmt.Print(gotDc)
	fmt.Print(err)
	err = SetPixel(gotDc, 23, 10, 0)
	fmt.Print(err)
	c, err := GetPixel(gotDc, 23, 10)
	fmt.Print(c)
	fmt.Print(err)
	er := ResetDC(gotDc)
	fmt.Print(err)
	fmt.Print(er)
}

func TestImageEncode(t *testing.T) {

	printName := "Microsoft Print to PDF"
	dc, err := CreateDC(printName)
	fmt.Print(err)
	StartDCPrinter(dc, "gdiDoc")
	StartPage(dc)

	pix := []byte{255, 255, 255, 150, 146, 246, 36, 28, 237, 0, 0, 0,
		36, 36, 36, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		73, 73, 73, 24, 24, 24, 133, 133, 133, 0, 0, 0}
	fmt.Print(err)
	err = DrawDIImage(dc, 600, 600, 720, 360, 0, 0, 3, 3, pix)
	fmt.Print(err)
	EndPage(dc)
	EndDoc(dc)
	DeleteDC(dc)
}

func TestPrintLine(t *testing.T) {
	printName := "Microsoft Print to PDF"
	dc, err := CreateDC(printName)
	fmt.Print(err)
	StartDCPrinter(dc, "gdiDoc")
	StartPage(dc)
	p, err := MoveTo(dc, 120, 150)
	fmt.Print(p)
	fmt.Print(err)
	LineTo(dc, 500, 150)
	EndPage(dc)
	EndDoc(dc)
	DeleteDC(dc)
}

func TestGetDeviceCaps(t *testing.T) {
	hdc, err := CreateDC("Microsoft Print to PDF")
	if err != nil {
		t.Errorf("CreateDC() error = %v ", err)
	}
	GetDeviceCaps(hdc, HORZRES)
	GetDeviceCaps(hdc, VERTRES)
}

func TestBmp(t *testing.T) {
	file, err := os.Open("C:\\Users\\Desktop\\test.bmp")
	t.Logf("%v", err)
	bmp.Decode(file)
}

func TestPrintPNG(t *testing.T) {
	printName := "Microsoft Print to PDF"
	dc, err := CreateDC(printName)
	fmt.Print(err)
	StartDCPrinter(dc, "gdiDoc")
	StartPage(dc)
	file, err := os.Open("C:\\Users\\Desktop\\test.png")
	fmt.Print(err)
	image, err := png.Decode(file)
	fmt.Print(err)
	bgr := bgr2.NewBGRImage(image.Bounds())
	draw.Draw(bgr, image.Bounds(), image, image2.Point{0, 0}, draw.Src)
	src := bgr2.ReverseDIB(bgr.Pix, image.Bounds().Dx(), image.Bounds().Dy(), 24)
	DrawDIImage(dc, 0, uint32(image.Bounds().Dy())*10, uint32(image.Bounds().Dx())*10, uint32(image.Bounds().Dy())*10, 0, 0, int32(image.Bounds().Dx()), int32(image.Bounds().Dy()), src)
	EndPage(dc)
	EndDoc(dc)
	DeleteDC(dc)
}

func TestPrintJPG(t *testing.T) {
	printName := "Microsoft Print to PDF"
	dc, err := CreateDC(printName)
	fmt.Print(err)
	StartDCPrinter(dc, "gdiDoc")
	StartPage(dc)

	file, err := os.Open("C:\\Users\\Desktop\\test.jpg")
	fmt.Print(err)
	image, err := jpeg.Decode(file)
	fmt.Print(err)
	bgr := bgr2.NewBGRImage(image.Bounds())
	draw.Draw(bgr, image.Bounds(), image, image2.Point{0, 0}, draw.Src)
	src := bgr2.ReverseDIB(bgr.Pix, image.Bounds().Dx(), image.Bounds().Dy(), 24)
	DrawDIImage(dc, 0, uint32(image.Bounds().Dy()), uint32(image.Bounds().Dx()), uint32(image.Bounds().Dy()), 0, 0, int32(image.Bounds().Dx()), int32(image.Bounds().Dy()), src)
	EndPage(dc)
	EndDoc(dc)
	DeleteDC(dc)
}

func TestPrintBMP(t *testing.T) {
	printName := "Microsoft Print to PDF"
	dc, err := CreateDC(printName)
	fmt.Print(err)
	StartDCPrinter(dc, "gdiDoc")
	StartPage(dc)

	file, err := os.Open("C:\\Users\\Desktop\\test.jpg")
	fmt.Print(err)
	image, err := jpeg.Decode(file)
	fmt.Print(err)
	bgr := bgr2.NewBGRImage(image.Bounds())
	draw.Draw(bgr, image.Bounds(), image, image2.Point{0, 0}, draw.Src)
	src := bgr2.ReverseDIB(bgr.Pix, image.Bounds().Dx(), image.Bounds().Dy(), 24)
	DrawDIImage(dc, 0, uint32(image.Bounds().Dy()), uint32(image.Bounds().Dx()), uint32(image.Bounds().Dy()), 0, 0, int32(image.Bounds().Dx()), int32(image.Bounds().Dy()), src)
	EndPage(dc)
	EndDoc(dc)
	DeleteDC(dc)
}

func TestRGBA(t *testing.T) {
	/*image2.NewRGBA64()
	image2.NewRGBA()
	image2.NewYCbCr()*/
}
