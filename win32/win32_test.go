package win32

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"
)

func TestOpenPrinter(t *testing.T) {
	got, err := OpenPrinter("Microsoft Print to PDF")
	fmt.Print(got)
	fmt.Print(err)
}

func TestClosePrinter(t *testing.T) {
	got, err := OpenPrinter("Microsoft Print to PDF")
	fmt.Print(got)
	fmt.Print(err)

	ClosePrinter(got)
}

func TestImage(t *testing.T) {

}

func TestStartDocPrinter(t *testing.T) {
	got, err := OpenPrinter("Microsoft Print to PDF")
	fmt.Print(got)
	fmt.Print(err)
	doc := &DOC_INFO_1{
		DocName:    windows.StringToUTF16Ptr("TEST"),
		OutputFile: windows.StringToUTF16Ptr(""),
		Datatype:   windows.StringToUTF16Ptr("RAW"),
	}
	err = StartDocPrinter(got, 1, doc)
	fmt.Print(err)
}

func TestEndDocPrinter(t *testing.T) {
	got, err := OpenPrinter("Microsoft Print to PDF")
	fmt.Print(got)
	fmt.Print(err)
	doc := &DOC_INFO_1{
		DocName:    windows.StringToUTF16Ptr("TEST"),
		OutputFile: windows.StringToUTF16Ptr(""),
		Datatype:   windows.StringToUTF16Ptr("RAW"),
	}
	err = StartDocPrinter(got, 1, doc)
	fmt.Print(err)
	err = EndDocPrinter(got)
	fmt.Print(err)
}

func TestStartPagePrinter(t *testing.T) {
	got, err := OpenPrinter("Microsoft Print to PDF")
	fmt.Print(got)
	fmt.Print(err)
	doc := &DOC_INFO_1{
		DocName:    windows.StringToUTF16Ptr("TEST"),
		OutputFile: windows.StringToUTF16Ptr(""),
		Datatype:   windows.StringToUTF16Ptr("RAW"),
	}
	err = StartDocPrinter(got, 1, doc)
	fmt.Print(err)
	err = StartPagePrinter(got)
	fmt.Print(err)
	err = EndPagePrinter(got)
	fmt.Print(err)
	err = StartPagePrinter(got)
	fmt.Print(err)
	err = EndPagePrinter(got)
	fmt.Print(err)
	err = StartPagePrinter(got)
	fmt.Print(err)
	err = EndPagePrinter(got)
	fmt.Print(err)
	err = EndDocPrinter(got)
	fmt.Print(err)
}

func TestPrint(t *testing.T) {
	printName := "Microsoft Print to PDF"
	got, err := OpenPrinter(printName)
	doc := &DOC_INFO_1{
		DocName:    windows.StringToUTF16Ptr("myDoc"),
		OutputFile: windows.StringToUTF16Ptr(""),
		Datatype:   windows.StringToUTF16Ptr("RAW"),
	}
	err = StartDocPrinter(got, 1, doc)
	fmt.Print(err)
	err = StartPagePrinter(got)
	fmt.Print(err)
	dc, err := CreateDC(printName)
	TextOut(dc, 1, 1, "TEST", 4)
	err = EndPagePrinter(got)
	fmt.Print(err)
	err = EndDocPrinter(got)
	fmt.Print(err)
	DeleteDC(dc)
	err = ClosePrinter(got)
	fmt.Print(err)
}

func TestGdiPrint(t *testing.T) {
	printName := "Microsoft Print to PDF"
	dc, err := CreateDC(printName)
	fmt.Print(err)
	StartDCPrinter(dc, "gdiDoc")
	StartPage(dc)
	TextOut(dc, 10, 10, "哈哈哈哈", 4)
	EndPage(dc)
	EndDoc(dc)
	DeleteDC(dc)
}

func imageToRGBA(img image.Image) *[]uint8 {
	sz := img.Bounds()
	raw := make([]uint8, (sz.Max.X-sz.Min.X)*(sz.Max.Y-sz.Min.Y)*4)
	idx := 0
	for y := sz.Min.Y; y < sz.Max.Y; y++ {
		for x := sz.Min.X; x < sz.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			raw[idx], raw[idx+1], raw[idx+2], raw[idx+3] = uint8(r), uint8(g), uint8(b), uint8(a)
			idx += 4
		}
	}
	return &raw
}

func TestEnumPrinter(t *testing.T) {
	var need uint32 = 0
	var returned uint32 = 0
	EnumPrinter(PRINTER_ENUM_LOCAL|PRINTER_ENUM_CONNECTIONS, nil, 4, nil, 0, &need, &returned)
	buf := make([]byte, need)
	EnumPrinter(PRINTER_ENUM_LOCAL|PRINTER_ENUM_CONNECTIONS, nil, 4, &buf[0], need, &need, &returned)
	ps := (*[1024]PRINT_INFO_4)(unsafe.Pointer(&buf[0]))[:returned:returned]

	printInfos := make([]PrinterInfo, 0, returned)
	for _, p := range ps {
		printInfos = append(printInfos, PrinterInfo{PrinterName: windows.UTF16PtrToString(p.PrinterName),
			ServerName: windows.UTF16PtrToString(p.ServerName),
			Attributes: p.Attributes})
	}
	fmt.Println(printInfos)
}

func TestDefault(t *testing.T) {
	printer, err := Default()
	fmt.Println(err)
	fmt.Println(printer)
	SetDefaultPrinter(windows.StringToUTF16Ptr("Microsoft Print to PDF"))
}
