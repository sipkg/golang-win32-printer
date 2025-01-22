package main

import (
	"fmt"
	"image/draw"
	"image/png"
	"os"

	image2 "image"
	bgr2 "print/image/bgr"
	"print/win32"
)

func main() {
	// printName := "Microsoft Print to PDF"
	printName := "PDF"
	dc, err := win32.CreateDC(printName)
	fmt.Print(err)
	win32.StartDCPrinter(dc, "gdiDoc")
	win32.StartPage(dc)
	file, err := os.Open("demo.png")
	fmt.Print(err)
	image, err := png.Decode(file)
	fmt.Print(err)
	bgr := bgr2.NewBGRImage(image.Bounds())
	draw.Draw(bgr, image.Bounds(), image, image2.Point{0, 0}, draw.Src)
	src := bgr2.ReverseDIB(bgr.Pix, image.Bounds().Dx(), image.Bounds().Dy(), 24)
	win32.DrawDIImage(dc, 0, uint32(image.Bounds().Dy())*10, uint32(image.Bounds().Dx())*10, uint32(image.Bounds().Dy())*10, 0, 0, int32(image.Bounds().Dx()), int32(image.Bounds().Dy()), src)
	win32.TextOut(dc, 10, 10, "SOFTINNOV", 9)
	win32.EndPage(dc)
	win32.EndDoc(dc)
	win32.DeleteDC(dc)
}
