# golang-win32-printer

Use Go language to encapsulate win32 print API to support printing
pictures, strings, and files.

# win32 function encapsulated

Please refer to the win32 API document for details.

1. DeleteDC
2. CreateDC
3. TextOut
4. StartDoc
5. EndDoc
6. StartPage
7. EndPage
8. OpenPrinter
9. ClosePrinter
10. StartDocPrinter
11. CloseDocPrinter
12. StartPagePrinter
13. ClosePagePrinter
14. ResetDC
15. SetPixel
16. GetPixel
17. GetDeviceCaps
18. StretchDIBits
19. MoveTo
20. LineTo
21. EnumPrinter
22. GetDefaultPrinter
23. SetDefaultPrinter

## Package Structure

- golang-win32-printer
  - image: BGR format image wrapper, supports 24-bit BPP
  - printer: win32 API logic wrapper
  - win32: system call API encapsulation (inclugind gdi32)

## Current Printing Flow

BGRImage encapsulation handles drawing functions to write images, text, lines,
and rectangles into a temporary BGRImage buffer, then copies to printer HDC
output via StretchDIBits

## Code Example

See examples folder.

Features under development. Contributors welcome.
