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

	"print/image/bgr"
	"print/win32"
)

func main() {
	// printName := "Microsoft Print to PDF"
	printName := "PDF"
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

	// Load an image and convert it to be printable
	file, err := os.Open("demo.png")
	if err != nil {
		log.Fatalf("Could not load image: %s", err)
	}
	img, err := png.Decode(file)
	if err != nil {
		log.Fatalf("Could not decode image: %s", err)
	}
	imgbgr := bgr.NewBGRImage(img.Bounds())
	draw.Draw(imgbgr, img.Bounds(), img, image.Point{0, 0}, draw.Src)
	src := bgr.ReverseDIB(imgbgr.Pix, img.Bounds().Dx(), img.Bounds().Dy(), 24)
	err = win32.DrawDIImage(dc, 0, uint32(img.Bounds().Dy())*10, uint32(img.Bounds().Dx())*10, uint32(img.Bounds().Dy())*10, 0, 0, int32(img.Bounds().Dx()), int32(img.Bounds().Dy()), src)
	if err != nil {
		log.Printf("DrawDIImage failed: %s", err)
	}

	oldcol, err := win32.SetTextColor(dc, win32.RGB(250, 35, 20))
	if err != nil {
		log.Printf("SetTextColor failed: %s", err)
	}
	log.Printf("Before color was %v", oldcol)

	err = win32.SetTextSize(dc, 128)
	if err != nil {
		log.Printf("SetTextColor failed: %s", err)
	}

	err = win32.TextOut(dc, 120, 50, "HELLO WORLD", 11)
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
