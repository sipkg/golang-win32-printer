// win32/textsize.go

package win32

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Constants
const (
	HFONT = 6 // Object type for getCurrentFont
)

const (
	Arial         = "Arial"
	TimesNewRoman = "Times New Roman"
	CourierNew    = "Courier New"
	Verdana       = "Verdana"
)

// LOGFONT contains information about a logical font
type LOGFONT struct {
	Height         int32
	Width          int32
	Escapement     int32
	Orientation    int32
	Weight         int32
	Italic         byte
	Underline      byte
	StrikeOut      byte
	CharSet        byte
	OutPrecision   byte
	ClipPrecision  byte
	Quality        byte
	PitchAndFamily byte
	FaceName       [32]uint16
}

// COLORREF represents a Windows color value in BGR format
type COLORREF uint32

// RGB creates a COLORREF value from red, green, and blue components
func RGB(r, g, b byte) COLORREF {
	return COLORREF(uint32(b)<<16 | uint32(g)<<8 | uint32(r))
}

// SetTextColor sets the text color for the specified device context
// Parameters:
//   - hdc: Handle to the device context
//   - color: Color value created using RGB()
//
// Returns:
//   - Previous text color value
//   - error: nil if successful, error object otherwise
func SetTextColor(hdc HDC, color COLORREF) (COLORREF, error) {
	ret, _, err := syscall.SyscallN(procSetTextColor.Addr(), uintptr(hdc), uintptr(color))

	if ret == ^uintptr(0) { // CLR_INVALID
		if err != windows.ERROR_SUCCESS {
			return 0, err
		}
	}

	return COLORREF(ret), nil
}

// SetTextSize changes the text height for the specified device context
// Parameters:
//   - hdc: Handle to the device context
//   - size: Text height in logical units (points)
//
// Returns:
//   - error: nil if successful, error object otherwise
func SetTextSize(hdc HDC, size int32) (i int32, err error) {
	// Get current font
	var lf LOGFONT
	font, _, errno := syscall.SyscallN(procGetCurrentObject.Addr(), uintptr(hdc), uintptr(OBJ_FONT))
	if errno != 0 {
		err = errno
	}
	if err != nil {
		return 0, fmt.Errorf("syscall getcurrobject : %v", err)
	}

	// Get font information
	_, _, errno = syscall.SyscallN(procGetObject.Addr(),
		font,
		uintptr(unsafe.Sizeof(lf)),
		uintptr(unsafe.Pointer(&lf)),
	)
	if errno != 0 {
		err = errno
	}
	if err != nil {
		return 0, fmt.Errorf("syscall getobject : %v", err)
	}

	originalHeight := lf.Height

	// Update height (negative value for character height)
	lf.Height = -size
	// Update width proportionally (adjust the factor if needed)
	lf.Width = size / 2

	// Create new font
	newFont, _, errno := syscall.SyscallN(procCreateFontIndirect.Addr(), uintptr(unsafe.Pointer(&lf)))
	if errno != 0 {
		err = errno
	}
	if err != nil {
		return 0, fmt.Errorf("syscall createfont : %v", err)
	}

	// Select new font into DC
	oldFont, _, errno := syscall.SyscallN(procSelectObject.Addr(),
		uintptr(hdc),
		newFont,
	)
	if errno != 0 {
		err = errno
	}
	if err != nil || oldFont == 0 {
		syscall.SyscallN(procDeleteObject.Addr(), newFont)
		return 0, fmt.Errorf("syscall deleteobject : %v", err)
	}

	// Delete old font if it exists
	if font != 0 {
		syscall.SyscallN(procDeleteObject.Addr(), font)
	}

	return -originalHeight, nil
}

func SetBoldFont(hdc HDC, bold bool) (err error) {
	// Get current font
	var lf LOGFONT
	font, _, errno := syscall.SyscallN(procGetCurrentObject.Addr(), uintptr(hdc), uintptr(OBJ_FONT))
	if errno != 0 {
		err = errno
	}
	if err != nil {
		return err
	}

	// Get font information
	_, _, errno = syscall.SyscallN(procGetObject.Addr(),
		font,
		uintptr(unsafe.Sizeof(lf)),
		uintptr(unsafe.Pointer(&lf)),
	)
	if errno != 0 {
		err = errno
	}
	if err != nil {
		return err
	}

	// Update weight to bold
	if bold {
		lf.Weight = 700 // FW_BOLD
	} else {
		lf.Weight = 400 //FW_NORMAL
	}

	// Create new font
	newFont, _, errno := syscall.SyscallN(procCreateFontIndirect.Addr(), uintptr(unsafe.Pointer(&lf)))
	if errno != 0 {
		err = errno
	}
	if err != nil {
		return err
	}

	// Select new font into DC
	oldFont, _, errno := syscall.SyscallN(procSelectObject.Addr(),
		uintptr(hdc),
		newFont,
	)
	if errno != 0 {
		err = errno
	}
	if err != nil || oldFont == 0 {
		syscall.SyscallN(procDeleteObject.Addr(), newFont)
		return err
	}

	// Delete old font if it exists
	if font != 0 {
		syscall.SyscallN(procDeleteObject.Addr(), font)
	}

	return nil
}

func SetItalicFont(hdc HDC, italic bool) (err error) {
	// Get current font
	var lf LOGFONT
	font, _, errno := syscall.SyscallN(procGetCurrentObject.Addr(), uintptr(hdc), uintptr(OBJ_FONT))
	if errno != 0 {
		err = errno
	}
	if err != nil {
		return err
	}

	// Get font information
	_, _, errno = syscall.SyscallN(procGetObject.Addr(),
		font,
		uintptr(unsafe.Sizeof(lf)),
		uintptr(unsafe.Pointer(&lf)),
	)
	if errno != 0 {
		err = errno
	}
	if err != nil {
		return err
	}

	// Update italic property
	if italic {
		lf.Italic = 1
	} else {
		lf.Italic = 0
	}

	// Create new font
	newFont, _, errno := syscall.SyscallN(procCreateFontIndirect.Addr(), uintptr(unsafe.Pointer(&lf)))
	if errno != 0 {
		err = errno
	}
	if err != nil {
		return err
	}

	// Select new font into DC
	oldFont, _, errno := syscall.SyscallN(procSelectObject.Addr(),
		uintptr(hdc),
		newFont,
	)
	if errno != 0 {
		err = errno
	}
	if err != nil || oldFont == 0 {
		syscall.SyscallN(procDeleteObject.Addr(), newFont)
		return err
	}

	// Delete old font if it exists
	if font != 0 {
		syscall.SyscallN(procDeleteObject.Addr(), font)
	}

	return nil
}

func SetFont(hdc HDC, fontName string) (err error) {
	// Get current font
	var lf LOGFONT
	font, _, errno := syscall.SyscallN(procGetCurrentObject.Addr(), uintptr(hdc), uintptr(OBJ_FONT))
	if errno != 0 {
		err = errno
	}
	if err != nil {
		return err
	}

	// Get font information
	_, _, errno = syscall.SyscallN(procGetObject.Addr(),
		font,
		uintptr(unsafe.Sizeof(lf)),
		uintptr(unsafe.Pointer(&lf)),
	)
	if errno != 0 {
		err = errno
	}
	if err != nil {
		return err
	}

	// Update font properties
	f, er := syscall.UTF16FromString(fontName)
	if er != nil {
		// Gérer l'erreur ici
		log.Fatalf("Erreur lors de la conversion de la chaîne en UTF-16 : %v", err)
	}
	copy(lf.FaceName[:], f)

	// Create new font
	newFont, _, errno := syscall.SyscallN(procCreateFontIndirect.Addr(), uintptr(unsafe.Pointer(&lf)))
	if errno != 0 {
		err = errno
	}
	if err != nil {
		return err
	}

	// Select new font into DC
	oldFont, _, errno := syscall.SyscallN(procSelectObject.Addr(),
		uintptr(hdc),
		newFont,
	)
	if errno != 0 {
		err = errno
	}
	if err != nil || oldFont == 0 {
		syscall.SyscallN(procDeleteObject.Addr(), newFont)
		return err
	}

	// Delete old font if it exists
	if font != 0 {
		syscall.SyscallN(procDeleteObject.Addr(), font)
	}

	return nil
}
