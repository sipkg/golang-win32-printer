// win32/textsize.go

package win32

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Constants
const (
	HFONT = 6 // Object type for getCurrentFont
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
func SetTextSize(hdc HDC, size int32) error {
	// Get current font
	var lf LOGFONT
	font, _, err := syscall.SyscallN(procGetCurrentObject.Addr(), uintptr(hdc), uintptr(OBJ_FONT))
	if font == 0 {
		return err
	}

	// Get font information
	ret, _, err := syscall.SyscallN(procGetObject.Addr(),
		font,
		uintptr(unsafe.Sizeof(lf)),
		uintptr(unsafe.Pointer(&lf)),
	)
	if ret == 0 {
		return err
	}

	// Update height (negative value for character height)
	lf.Height = -size
	// Update width proportionally (adjust the factor if needed)
	lf.Width = size / 2

	// Create new font
	newFont, _, err := syscall.SyscallN(procCreateFontIndirect.Addr(), uintptr(unsafe.Pointer(&lf)))
	if newFont == 0 {
		return err
	}

	// Select new font into DC
	oldFont, _, err := syscall.SyscallN(procSelectObject.Addr(),
		uintptr(hdc),
		newFont,
	)
	if oldFont == 0 {
		syscall.SyscallN(procDeleteObject.Addr(), newFont)
		return err
	}

	// Delete old font if it exists
	if font != 0 {
		syscall.SyscallN(procDeleteObject.Addr(), font)
	}

	return nil
}

func SetBoldFont(hdc HDC, bold bool) error {
	// Get current font
	var lf LOGFONT
	font, _, err := syscall.SyscallN(procGetCurrentObject.Addr(), uintptr(hdc), uintptr(OBJ_FONT))
	if font == 0 {
		return err
	}

	// Get font information
	ret, _, err := syscall.SyscallN(procGetObject.Addr(),
		font,
		uintptr(unsafe.Sizeof(lf)),
		uintptr(unsafe.Pointer(&lf)),
	)
	if ret == 0 {
		return err
	}

	// Update weight to bold
	if bold {
		lf.Weight = 700 // FW_BOLD
	} else {
		lf.Weight = 400 //FW_NORMAL
	}

	// Create new font
	newFont, _, err := syscall.SyscallN(procCreateFontIndirect.Addr(), uintptr(unsafe.Pointer(&lf)))
	if newFont == 0 {
		return err
	}

	// Select new font into DC
	oldFont, _, err := syscall.SyscallN(procSelectObject.Addr(),
		uintptr(hdc),
		newFont,
	)
	if oldFont == 0 {
		syscall.SyscallN(procDeleteObject.Addr(), newFont)
		return err
	}

	// Delete old font if it exists
	if font != 0 {
		syscall.SyscallN(procDeleteObject.Addr(), font)
	}

	return nil
}

func SetItalicFont(hdc HDC, italic bool) error {
	// Get current font
	var lf LOGFONT
	font, _, err := syscall.SyscallN(procGetCurrentObject.Addr(), uintptr(hdc), uintptr(OBJ_FONT))
	if font == 0 {
		return err
	}

	// Get font information
	ret, _, err := syscall.SyscallN(procGetObject.Addr(),
		font,
		uintptr(unsafe.Sizeof(lf)),
		uintptr(unsafe.Pointer(&lf)),
	)
	if ret == 0 {
		return err
	}

	// Update italic property
	if italic {
		lf.Italic = 1
	} else {
		lf.Italic = 0
	}

	// Create new font
	newFont, _, err := syscall.SyscallN(procCreateFontIndirect.Addr(), uintptr(unsafe.Pointer(&lf)))
	if newFont == 0 {
		return err
	}

	// Select new font into DC
	oldFont, _, err := syscall.SyscallN(procSelectObject.Addr(),
		uintptr(hdc),
		newFont,
	)
	if oldFont == 0 {
		syscall.SyscallN(procDeleteObject.Addr(), newFont)
		return err
	}

	// Delete old font if it exists
	if font != 0 {
		syscall.SyscallN(procDeleteObject.Addr(), font)
	}

	return nil
}
