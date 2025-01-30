// Simple demo that can be cross compiled from Linux to
// Windows and executed via Wine with a printer named `PDF`
//
// ```go
// GOOS=windows GOARCH=amd64 go build cmd
// ```
package main

import (
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"syscall"

	"print/image/bgr"
	"print/layout"
	"print/win32"
)

const (
	textHeight = 128
	margin     = 150
)

func main() {
	printName := "Microsoft Print to PDF"
	// printName := "PDF"
	dc, err := win32.CreateDC(printName)
	if err != nil {
		log.Fatalf("CreateDC failed: %s", err)
	}
	err = win32.StartDCPrinter(dc, "gdiDoc")
	if err != nil {
		log.Fatalf("StartDCPrinter failed: %s", err)
	}
	err = win32.StartPage(dc)
	if err != nil {
		log.Fatalf("StartPage failed: %s", err)
	}

	win32.SetFont(dc, win32.CourierNew)
	width, err := win32.GetDeviceCaps(dc, win32.HORZRES)
	if err != nil {
		log.Fatalf("Retreiving page width failed: %s", err)
	}
	height, err := win32.GetDeviceCaps(dc, win32.VERTRES)
	if err != nil {
		log.Fatalf("Retreiving page height failed: %s", err)
	}
	log.Printf("Page dimensions: %d x %d pixels", width, height)

	// Load an image and convert it to be printable
	file, err := os.Open("demo.png")
	if err != nil {
		log.Fatalf("Could not load image: %s", err)
	}
	img, err := png.Decode(file)
	if err != nil {
		log.Fatalf("Could not decode image: %s", err)
	}

	imgWidth := uint32(img.Bounds().Dx()) * 5
	imgHeight := uint32(img.Bounds().Dy()) * 5
	imgX := layout.CenterElement(uint32(width), imgWidth)
	imgY := layout.CenterElement(uint32(height), imgHeight)

	imgbgr := bgr.NewBGRImage(img.Bounds())
	draw.Draw(imgbgr, img.Bounds(), img, image.Point{0, 0}, draw.Src)
	src := bgr.ReverseDIB(imgbgr.Pix, img.Bounds().Dx(), img.Bounds().Dy(), 24)
	err = win32.DrawDIImage(dc, imgX, imgY, imgWidth, imgHeight, 0, 0, int32(img.Bounds().Dx()), int32(img.Bounds().Dy()), src)
	if err != nil {
		log.Printf("DrawDIImage failed: %s", err)
	}

	oldcol, err := win32.SetTextColor(dc, win32.RGB(200, 20, 80))
	if err != nil {
		log.Printf("SetTextColor failed: %s", err)
	}
	log.Printf("Before color was %v", oldcol)

	_, err = win32.SetTextSize(dc, textHeight)
	if err != nil {
		log.Printf("SetTextSize failed: %s", err)
	}

	text := "Hello world"
	textWidth, _, err := win32.GetTextExtentPoint32(syscall.Handle(dc), text)
	if err != nil {
		log.Printf("Get text dimensions failed: %s", err)
	}
	x := layout.CenterElement(width, textWidth)
	startY := uint32(10)
	err = win32.TextOut(dc, x, startY, text, uint32(len(text)))
	if err != nil {
		log.Printf("TextOut failed: %s", err)
	}

	err = win32.EndPage(dc)
	if err != nil {
		log.Printf("EndPage failed: %s", err)
	}
	err = win32.EndDoc(dc)
	if err != nil {
		log.Printf("EndDoc failed: %s", err)
	}
	err = win32.DeleteDC(dc)
	if err != nil {
		log.Printf("DeleteDC failed: %s", err)
	}
}
