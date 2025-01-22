package win32

import (
	"image"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	gdi32 = syscall.NewLazyDLL("Gdi32.dll")

	procCreateDCW = gdi32.NewProc("CreateDCW")
	procDeleteDCW = gdi32.NewProc("DeleteDC")
	procResetDCW  = gdi32.NewProc("ResetDCW")

	procTextOutW      = gdi32.NewProc("TextOutW")
	procSetPixel      = gdi32.NewProc("SetPixel")
	procGetPixel      = gdi32.NewProc("GetPixel")
	procMoveTo        = gdi32.NewProc("MoveToEx")
	procLineTo        = gdi32.NewProc("LineTo")
	procGetDeviceCaps = gdi32.NewProc("GetDeviceCaps")
	procStartDocW     = gdi32.NewProc("StartDocW")
	procStartPage     = gdi32.NewProc("StartPage")
	procEndDoc        = gdi32.NewProc("EndDoc")
	procEndPage       = gdi32.NewProc("EndPage")
	procStretchDIBits = gdi32.NewProc("StretchDIBits")

	procSetTextColor = gdi32.NewProc("SetTextColor")

	procGetStockObject = gdi32.NewProc("GetStockObject")

	procGetCurrentObject   = gdi32.NewProc("GetCurrentObject")
	procGetObject          = gdi32.NewProc("GetObjectW")
	procSelectObject       = gdi32.NewProc("SelectObject")
	procDeleteObject       = gdi32.NewProc("DeleteObject")
	procCreateFontIndirect = gdi32.NewProc("CreateFontIndirectW")
)

type HDC syscall.Handle

// CreateDC
func CreateDC(printerName string) (dc HDC, err error) {
	driver := windows.StringToUTF16Ptr("WINSPOOL")
	device := windows.StringToUTF16Ptr(printerName)
	r1, _, e1 := syscall.SyscallN(procCreateDCW.Addr(), uintptr(unsafe.Pointer(driver)), uintptr(unsafe.Pointer(device)), uintptr(0), uintptr(0))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return HDC(r1), err
}

// ResetDC https://learn.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-resetdcw
func ResetDC(dc HDC) (err error) {
	r1, _, e1 := syscall.SyscallN(procResetDCW.Addr(), uintptr(dc), uintptr(0))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}

// DeleteDC https://learn.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-deletedc
func DeleteDC(dc HDC) (err error) {
	r1, _, e1 := syscall.SyscallN(procDeleteDCW.Addr(), uintptr(dc))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}

// SetPixel https://learn.microsoft.com/zh-cn/windows/win32/api/wingdi/nf-wingdi-setpixel
func SetPixel(dc HDC, x, y uint32, color uint32) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetPixel.Addr(), uintptr(dc), uintptr(x), uintptr(y), uintptr(color))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}

// GetPixel https://learn.microsoft.com/zh-cn/windows/win32/api/wingdi/nf-wingdi-getpixel
func GetPixel(dc HDC, x, y uint32) (color uint32, err error) {
	r1, _, e1 := syscall.SyscallN(procGetPixel.Addr(), uintptr(dc), uintptr(x), uintptr(y))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return uint32(r1), err
}

func TextOut(dc HDC, x, y uint32, text string, len uint32) (err error) {
	str := windows.StringToUTF16Ptr(text)
	r1, _, e1 := syscall.SyscallN(procTextOutW.Addr(), uintptr(dc), uintptr(x), uintptr(y), uintptr(unsafe.Pointer(str)), uintptr(len))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}

type DOCINFOA struct {
	Size     uintptr
	DocName  *uint16
	Output   *uint16
	DataType *uint16
	Type     uint32
}

func StartDCPrinter(dc HDC, docName string) (err error) {
	size := unsafe.Sizeof(&DOCINFOA{})
	doc := &DOCINFOA{
		Size:     size,
		DocName:  windows.StringToUTF16Ptr(docName),
		Output:   nil,
		DataType: nil,
		Type:     0,
	}
	return StartDoc(dc, doc)
}

// https://learn.microsoft.com/zh-cn/windows/win32/api/wingdi/nf-wingdi-startdocw
func StartDoc(dc HDC, doc *DOCINFOA) (err error) {
	r1, _, e1 := syscall.SyscallN(procStartDocW.Addr(), uintptr(dc), uintptr(unsafe.Pointer(doc)))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}

// https://learn.microsoft.com/zh-cn/windows/win32/api/wingdi/nf-wingdi-startpage
func StartPage(dc HDC) (err error) {
	r1, _, e1 := syscall.SyscallN(procStartPage.Addr(), uintptr(dc))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}

// https://learn.microsoft.com/zh-cn/windows/win32/api/wingdi/nf-wingdi-enddoc
func EndDoc(dc HDC) (err error) {
	r1, _, e1 := syscall.SyscallN(procEndDoc.Addr(), uintptr(dc))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}

func EndPage(dc HDC) (err error) {
	r1, _, e1 := syscall.SyscallN(procEndPage.Addr(), uintptr(dc))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}

/* Device Parameters for GetDeviceCaps() */

type PropType uint32

const (
	DRIVERVERSION         PropType = 0  /* Device driver version                    */
	TECHNOLOGYPropType    PropType = 2  /* Device classification                    */
	HORZSIZE              PropType = 4  /* Horizontal size in millimeters           */
	VERTSIZE              PropType = 6  /* Vertical size in millimeters             */
	HORZRES               PropType = 8  /* Horizontal width in pixels               */
	VERTRES               PropType = 10 /* Vertical height in pixels                */
	BITSPIXEL             PropType = 12 /* Number of bits per pixel                 */
	PLANES                PropType = 14 /* Number of planes                         */
	NUMBRUSHES            PropType = 16 /* Number of brushes the device has         */
	NUMPENS               PropType = 18 /* Number of pens the device has            */
	NUMMARKERS            PropType = 20 /* Number of markers the device has         */
	NUMFONTS              PropType = 22 /* Number of fonts the device has           */
	NUMCOLORS             PropType = 24 /* Number of colors the device supports     */
	PDEVICESIZE           PropType = 26 /* Size required for device descriptor      */
	CURVECAPS             PropType = 28 /* Curve capabilities                       */
	LINECAPSPropType      PropType = 30 /* Line capabilities                        */
	POLYGONALCAPSPropType PropType = 32 /* Polygonal capabilities                   */
	TEXTCAPSPropType      PropType = 34 /* Text capabilities                        */
	CLIPCAPSPropType      PropType = 36 /* Clipping capabilities                    */
	RASTERCAPS            PropType = 38 /* Bitblt capabilities                      */
	ASPECTX               PropType = 40 /* Length of the X leg                      */
	ASPECTY               PropType = 42 /* Length of the Y leg                      */
	ASPECTXYPropType      PropType = 44 /* Length of the hypotenuse                 */

	LOGPIXELSX PropType = 88 /* Logical pixels/inch in X                 */
	LOGPIXELSY PropType = 90 /* Logical pixels/inch in Y                 */

	SIZEPALETTE PropType = 104 /* Number of entries in physical palette    */
	NUMRESERVED PropType = 106 /* Number of reserved entries in palette    */
	COLORRES    PropType = 108 /* Actual color resolution                  */

	// Printing related DeviceCaps. These replace the appropriate Escapes

	PHYSICALWIDTH   PropType = 110 /* Physical Width in device units           */
	PHYSICALHEIGHT  PropType = 111 /* Physical Height in device units          */
	PHYSICALOFFSETX PropType = 112 /* Physical Printable Area x margin         */
	PHYSICALOFFSETY PropType = 113 /* Physical Printable Area y margin         */
	SCALINGFACTORX  PropType = 114 /* Scaling factor x                         */
	SCALINGFACTORY  PropType = 115 /* Scaling factor y                         */

	// Display driver specific

	VREFRESH PropType = 116 /* Current vertical refresh rate of the    */
	/* display device (for displays only) in Hz*/
	DESKTOPVERTRES PropType = 117 /* Horizontal width of entire desktop in   */
	/* pixels                                  */
	DESKTOPHORZRES PropType = 118 /* Vertical height of entire desktop in    */
	/* pixels                                  */
	BLTALIGNMENT PropType = 119 /* Preferred blt alignment                 */

	SHADEBLENDCAPS PropType = 120 /* Shading and blending caps               */
	COLORMGMTCAPS  PropType = 121 /* Color Management caps                   */
)

func GetDeviceCaps(dc HDC, index PropType) (number uint32, err error) {
	r1, _, e1 := syscall.SyscallN(procGetDeviceCaps.Addr(), uintptr(dc), uintptr(index))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return uint32(r1), err
}

type BITMAPINFOHEADER struct {
	biSize          uint32
	biWidth         int32
	biHeight        int32
	biPlanes        uint16
	biBitCount      uint16
	biCompression   BiCompression
	biSizeImage     uint32
	biXPelsPerMeter uint32
	biYPelsPerMeter uint32
	biClrUsed       uint32
	biClrImportant  uint32
}

type RGBQUAD struct {
	RgbBlue     uint8
	RgbGreen    uint8
	RgbRed      uint8
	RgbReserved uint8
}

type BITMAPINFO struct {
	BmiHeader BITMAPINFOHEADER
	BmiColors RGBQUAD
}

type DIBColors uint32

const (
	DIB_RGB_COLORS  DIBColors = 0x00
	DIB_PAL_COLORS  DIBColors = 0x01
	DIB_PAL_INDICES DIBColors = 0x02
)

type RasterOperationCode uint32

const (
	SRCCOPY     RasterOperationCode = 0x00CC0020 /* dest = source                   */
	SRCPAINT    RasterOperationCode = 0x00EE0086 /* dest = source OR dest           */
	SRCAND      RasterOperationCode = 0x008800C6 /* dest = source AND dest          */
	SRCINVERT   RasterOperationCode = 0x00660046 /* dest = source XOR dest          */
	SRCERASE    RasterOperationCode = 0x00440328 /* dest = source AND (NOT dest )   */
	NOTSRCCOPY  RasterOperationCode = 0x00330008 /* dest = (NOT source)             */
	NOTSRCERASE RasterOperationCode = 0x001100A6 /* dest = (NOT src) AND (NOT dest) */
	MERGECOPY   RasterOperationCode = 0x00C000CA /* dest = (source AND pattern)     */
	MERGEPAINT  RasterOperationCode = 0x00BB0226 /* dest = (NOT source) OR dest     */
	PATCOPY     RasterOperationCode = 0x00F00021 /* dest = pattern                  */
	PATPAINT    RasterOperationCode = 0x00FB0A09 /* dest = DPSnoo                   */
	PATINVERT   RasterOperationCode = 0x005A0049 /* dest = pattern XOR dest         */
	DSTINVERT   RasterOperationCode = 0x00550009 /* dest = (NOT dest)               */
	BLACKNESS   RasterOperationCode = 0x00000042 /* dest = BLACK                    */
	WHITENESS   RasterOperationCode = 0x00FF0062 /* dest = WHITE                    */
)

type BiCompression uint32

const (
	BI_RGB       BiCompression = 0
	BI_RLE8      BiCompression = 1
	BI_RLE4      BiCompression = 2
	BI_BITFIELDS BiCompression = 3
	BI_JPEG      BiCompression = 4
	BI_PNG       BiCompression = 5
)

func DrawDIImage(dc HDC, dx, dy, dw, dh, sx, sy uint32, sw, sh int32, image []byte) (err error) {
	header := BITMAPINFOHEADER{
		biSize:          0,
		biWidth:         sw,
		biHeight:        -1 * sh,
		biPlanes:        1,
		biBitCount:      24,
		biCompression:   BI_RGB,
		biSizeImage:     0,
		biXPelsPerMeter: 0,
		biYPelsPerMeter: 0,
		biClrUsed:       0,
		biClrImportant:  0,
	}
	header.biSize = uint32(unsafe.Sizeof(header))
	bitmap := &BITMAPINFO{
		BmiHeader: header,
		BmiColors: RGBQUAD{},
	}

	return StretchDIBits(dc, dx, dy, dw, dh, sx, sy, sw, sh, image, bitmap, DIB_RGB_COLORS, SRCCOPY)
}

func StretchDIBits(dc HDC, dx, dy, dw, dh, sx, sy uint32, sw, sh int32, image []byte, bitmap *BITMAPINFO, color DIBColors, operation RasterOperationCode) (err error) {
	r1, _, e1 := syscall.SyscallN(procStretchDIBits.Addr(), uintptr(dc),
		uintptr(dx), uintptr(dy), uintptr(dw), uintptr(dh),
		uintptr(sx), uintptr(sy), uintptr(sw), uintptr(sh),
		uintptr(unsafe.Pointer(&image[0])), uintptr(unsafe.Pointer(bitmap)), uintptr(color), uintptr(operation))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}

func MoveTo(dc HDC, x, y uint32) (rect *image.Point, err error) {
	var point image.Point
	r1, _, e1 := syscall.SyscallN(procMoveTo.Addr(), uintptr(dc), uintptr(x), uintptr(y), uintptr(unsafe.Pointer(&point)))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return &point, err
}

func LineTo(dc HDC, x, y uint32) (err error) {
	r1, _, e1 := syscall.SyscallN(procLineTo.Addr(), uintptr(dc), uintptr(x), uintptr(y))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}
