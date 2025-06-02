package main

import (
	"log"
	"syscall"

	"github.com/clementuu/golang-win32-printer/layout"
	"github.com/clementuu/golang-win32-printer/win32"
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

	width, err := win32.GetDeviceCaps(dc, win32.HORZRES)
	if err != nil {
		log.Fatalf("Retreiving page width failed: %s", err)
	}
	_, err = win32.GetDeviceCaps(dc, win32.VERTRES)
	if err != nil {
		log.Fatalf("Retreiving page height failed: %s", err)
	}

	oldcol, err := win32.SetTextColor(dc, win32.RGB(0, 0, 0))
	if err != nil {
		log.Printf("SetTextColor failed: %s", err)
	}
	log.Printf("Before color was %v", oldcol)

	_, err = win32.SetTextSize(dc, textHeight)
	if err != nil {
		log.Printf("SetTextSize failed: %s", err)
	}

	text := "Original Size"
	textWidth, textHeight, err := win32.GetTextExtentPoint32(syscall.Handle(dc), text)
	if err != nil {
		log.Printf("Get text dimensions failed: %s", err)
	}
	x := layout.CenterElement(width, textWidth)
	startY := uint32(10)
	err = win32.TextOut(dc, x, startY, text, uint32(len(text)))
	if err != nil {
		log.Printf("TextOut failed: %s", err)
	}

	startY += textHeight + 20

	originalHeight, err := win32.SetTextSize(dc, int32(textHeight+100))
	if err != nil {
		log.Printf("SetTextSize failed: %s", err)
	}

	text = "Bigger text"
	textWidth, textHeight, err = win32.GetTextExtentPoint32(syscall.Handle(dc), text)
	if err != nil {
		log.Printf("Get text dimensions failed: %s", err)
	}
	x = layout.CenterElement(width, textWidth)
	err = win32.TextOut(dc, x, startY, text, uint32(len(text)))
	if err != nil {
		log.Printf("TextOut failed: %s", err)
	}

	startY += textHeight + 20

	_, err = win32.SetTextSize(dc, originalHeight)
	if err != nil {
		log.Printf("SetTextSize failed: %s", err)
	}

	text = "Original Size Again"
	textWidth, _, err = win32.GetTextExtentPoint32(syscall.Handle(dc), text)
	if err != nil {
		log.Printf("Get text dimensions failed: %s", err)
	}
	x = layout.CenterElement(width, textWidth)
	err = win32.TextOut(dc, x, startY, text, uint32(len(text)))
	if err != nil {
		log.Printf("TextOut failed: %s", err)
	}

	startY += textHeight + 20

	err = win32.SetBoldFont(dc, true)
	if err != nil {
		log.Printf("SetBold failed: %s", err)
	}

	text = "Bolder text"
	textWidth, _, err = win32.GetTextExtentPoint32(syscall.Handle(dc), text)
	if err != nil {
		log.Printf("Get text dimensions failed: %s", err)
	}
	x = layout.CenterElement(width, textWidth)
	err = win32.TextOut(dc, x, startY, text, uint32(len(text)))
	if err != nil {
		log.Printf("TextOut failed: %s", err)
	}

	startY += textHeight + 20

	err = win32.SetBoldFont(dc, false)
	if err != nil {
		log.Printf("SetBold failed: %s", err)
	}

	text = "Thiner text"
	textWidth, _, err = win32.GetTextExtentPoint32(syscall.Handle(dc), text)
	if err != nil {
		log.Printf("Get text dimensions failed: %s", err)
	}
	x = layout.CenterElement(width, textWidth)
	err = win32.TextOut(dc, x, startY, text, uint32(len(text)))
	if err != nil {
		log.Printf("TextOut failed: %s", err)
	}

	startY += textHeight + 20

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
