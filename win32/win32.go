package win32

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

var (
	printspool32 = syscall.NewLazyDLL("winspool.drv")

	procOpenPrinter        = printspool32.NewProc("OpenPrinterW")
	procClosePrinter       = printspool32.NewProc("ClosePrinter")
	procStartDocPrinter    = printspool32.NewProc("StartDocPrinterW")
	procStartPagePrinter   = printspool32.NewProc("StartPagePrinter")
	procEndDocPrinter      = printspool32.NewProc("EndDocPrinter")
	procEndPagePrinter     = printspool32.NewProc("EndPagePrinter")
	procWritePrinter       = printspool32.NewProc("WritePrinter")
	procEnumPrintersW      = printspool32.NewProc("EnumPrintersW")
	procGetDefaultPrinterW = printspool32.NewProc("GetDefaultPrinterW")
	procSetDefaultPrinter  = printspool32.NewProc("SetDefaultPrinterW")
)

type Printer syscall.Handle

func OpenPrinter(name string) (Printer, error) {
	var printHandler Printer
	nameptr := uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(name)))
	r1, _, e1 := syscall.SyscallN(procOpenPrinter.Addr(), nameptr, uintptr(unsafe.Pointer(&printHandler)), uintptr(0))
	var err error
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return printHandler, err
}

func ClosePrinter(handle Printer) (err error) {
	r1, _, e1 := syscall.SyscallN(procClosePrinter.Addr(), uintptr(handle))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}

// https://learn.microsoft.com/en-us/windows/win32/printdocs/doc-info-1
type DOC_INFO_1 struct {
	DocName    *uint16
	OutputFile *uint16
	Datatype   *uint16
}

// https://learn.microsoft.com/en-us/windows/win32/printdocs/startdocprinter
func StartDocPrinter(handle Printer, level uint32, doc *DOC_INFO_1) (err error) {
	r1, _, e1 := syscall.SyscallN(procStartDocPrinter.Addr(), uintptr(handle), uintptr(level), uintptr(unsafe.Pointer(doc)))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}

// https://learn.microsoft.com/zh-cn/windows/win32/printdocs/enddocprinter
func EndDocPrinter(handle Printer) (err error) {
	r1, _, e1 := syscall.SyscallN(procEndDocPrinter.Addr(), uintptr(handle))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}

func StartPagePrinter(handle Printer) (err error) {
	r1, _, e1 := syscall.SyscallN(procStartPagePrinter.Addr(), uintptr(handle))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}

func EndPagePrinter(handle Printer) (err error) {
	r1, _, e1 := syscall.SyscallN(procEndPagePrinter.Addr(), uintptr(handle))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}

func WritePrinter(handle Printer, buf *byte, bufN uint32, written *uint32) (err error) {
	r1, _, e1 := syscall.SyscallN(procWritePrinter.Addr(), uintptr(handle), uintptr(unsafe.Pointer(buf)), uintptr(bufN), uintptr(unsafe.Pointer(written)))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}

type EnumFlag uint32

const (
	PRINTER_ENUM_DEFAULT     EnumFlag = 0x00000001
	PRINTER_ENUM_LOCAL       EnumFlag = 0x00000002
	PRINTER_ENUM_CONNECTIONS EnumFlag = 0x00000004
	PRINTER_ENUM_FAVORITE    EnumFlag = 0x00000004
	PRINTER_ENUM_NAME        EnumFlag = 0x00000008
	PRINTER_ENUM_REMOTE      EnumFlag = 0x00000010
	PRINTER_ENUM_SHARED      EnumFlag = 0x00000020
	PRINTER_ENUM_NETWORK     EnumFlag = 0x00000040
)

type PRINT_INFO_4 struct {
	PrinterName *uint16
	ServerName  *uint16
	Attributes  uint32
}

type PrinterInfo struct {
	PrinterName string
	ServerName  string
	Attributes  uint32
}

const (
	PRINTER_ATTRIBUTE_LOCAL   uint32 = 0x00000040
	PRINTER_ATTRIBUTE_NETWORK uint32 = 0x00000010
)

// EnumPrinter https://learn.microsoft.com/zh-cn/windows/win32/printdocs/enumprinters
func EnumPrinter(flag EnumFlag, name *uint16, level uint32, info *byte, bufLen uint32, bufSize *uint32, returnLen *uint32) (err error) {
	r1, _, e1 := syscall.SyscallN(procEnumPrintersW.Addr(),
		uintptr(flag), uintptr(unsafe.Pointer(name)), uintptr(level), uintptr(unsafe.Pointer(info)),
		uintptr(bufLen), uintptr(unsafe.Pointer(bufSize)), uintptr(unsafe.Pointer(returnLen)))
	if r1 == 0 {
		if e1 != syscall.ERROR_INSUFFICIENT_BUFFER {
			if e1 != 0 {
				err = error(e1)
			} else {
				err = syscall.EINVAL
			}
		}
	}
	return err
}

func Default() (string, error) {
	b := make([]uint16, 3)
	n := uint32(len(b))
	err := GetDefaultPrinter(&b[0], &n)
	if err != nil {
		if err != syscall.ERROR_INSUFFICIENT_BUFFER {
			return "", err
		}
		b = make([]uint16, n)
		err = GetDefaultPrinter(&b[0], &n)
		if err != nil {
			return "", err
		}
	}
	return syscall.UTF16ToString(b), nil
}

func GetDefaultPrinter(buf *uint16, bufN *uint32) (err error) {
	r1, _, e1 := syscall.SyscallN(procGetDefaultPrinterW.Addr(), uintptr(unsafe.Pointer(buf)), uintptr(unsafe.Pointer(bufN)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetDefaultPrinter(buf *uint16) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetDefaultPrinter.Addr(), uintptr(unsafe.Pointer(buf)))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
