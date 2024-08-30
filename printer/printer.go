package printer

import (
	"errors"
	"golang.org/x/sys/windows"
	"print/win32"
	"unsafe"
)

const (
	DOCNAME = "JUNX PRINT"
)

type Printer struct {
	hdc   win32.HDC
	_init bool
}

func (p *Printer) InitPrinter(printerName string) (err error) {
	if p._init {
		return
	}
	p.hdc, err = win32.CreateDC(printerName)
	if err != nil {
		return
	}
	p._init = true
	return
}

func (p *Printer) EnumPrinter() (info []win32.PrinterInfo, err error) {
	var need uint32 = 0
	var returned uint32 = 0
	err = win32.EnumPrinter(win32.PRINTER_ENUM_LOCAL|win32.PRINTER_ENUM_CONNECTIONS, nil, 4, nil, 0, &need, &returned)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, need)
	err = win32.EnumPrinter(win32.PRINTER_ENUM_LOCAL|win32.PRINTER_ENUM_CONNECTIONS, nil, 4, &buf[0], need, &need, &returned)
	if err != nil {
		return nil, err
	}
	ps := (*[1024]win32.PRINT_INFO_4)(unsafe.Pointer(&buf[0]))[:returned:returned]
	printInfos := make([]win32.PrinterInfo, 0, returned)
	for _, p := range ps {
		printInfos = append(printInfos, win32.PrinterInfo{PrinterName: windows.UTF16PtrToString(p.PrinterName),
			ServerName: windows.UTF16PtrToString(p.ServerName),
			Attributes: p.Attributes})
	}
	return printInfos, nil
}

func (p *Printer) GetWidthPixel() (uint32, error) {
	width, err := win32.GetDeviceCaps(p.hdc, win32.HORZRES)
	if err != nil {
		return 0, err
	}
	return width, nil
}

func (p *Printer) GetHeightPixel() (uint32, error) {
	width, err := win32.GetDeviceCaps(p.hdc, win32.VERTRES)
	if err != nil {
		return 0, err
	}
	return width, nil
}

func (p *Printer) GetWidth() (uint32, error) {
	width, err := win32.GetDeviceCaps(p.hdc, win32.HORZSIZE)
	if err != nil {
		return 0, err
	}
	return width, nil
}

func (p *Printer) GetHeight() (uint32, error) {
	width, err := win32.GetDeviceCaps(p.hdc, win32.VERTSIZE)
	if err != nil {
		return 0, err
	}
	return width, nil
}

func (p *Printer) GetBPP() (uint32, error) {
	width, err := win32.GetDeviceCaps(p.hdc, win32.BITSPIXEL)
	if err != nil {
		return 0, err
	}
	return width, nil
}

func (p *Printer) GetMarginLeft() (uint32, error) {
	width, err := win32.GetDeviceCaps(p.hdc, win32.PHYSICALOFFSETX)
	if err != nil {
		return 0, err
	}
	return width, nil
}

func (p *Printer) GetMarginTop() (uint32, error) {
	width, err := win32.GetDeviceCaps(p.hdc, win32.PHYSICALOFFSETY)
	if err != nil {
		return 0, err
	}
	return width, nil
}

func (p *Printer) StartDoc(docName string) error {
	if !p._init {
		return errors.New("must call InitPrinter before")
	}
	size := unsafe.Sizeof(win32.DOCINFOA{})
	doc := &win32.DOCINFOA{
		Size:     size,
		DocName:  windows.StringToUTF16Ptr(docName),
		Output:   nil,
		DataType: nil,
		Type:     0,
	}
	return win32.StartDoc(p.hdc, doc)
}

func (p *Printer) StartPage() error {
	if !p._init {
		return errors.New("must call InitPrinter before")
	}
	return win32.StartPage(p.hdc)
}

func (p *Printer) EndDoc() error {
	if !p._init {
		return errors.New("must call InitPrinter before")
	}
	return win32.EndDoc(p.hdc)
}

func (p *Printer) EndPage() error {
	if !p._init {
		return errors.New("must call InitPrinter before")
	}
	return win32.EndPage(p.hdc)
}

type Printable interface {
	Print(*Printer)
}

func (p *Printer) Print(pt Printable) (err error) {
	if err = p.InitPrinter(""); err != nil {
		return
	}
	if err = p.StartDoc(DOCNAME); err != nil {
		return
	}
	if err = p.StartPage(); err != nil {
		return
	}
	pt.Print(p)
	if err = p.EndPage(); err != nil {
		return
	}
	if err = p.EndDoc(); err != nil {
		return
	}
	return
}
