package bgr

import (
	"image"
	"image/color"

	color2 "github.com/clementuu/golang-win32-printer/image/color"
)

type BGRImage struct {
	Pix    []uint8
	Stride int
	Rect   image.Rectangle
}

func (B BGRImage) ColorModel() color.Model {
	return color2.BGRModel
}

func (B *BGRImage) Bounds() image.Rectangle {
	return B.Rect
}

func (B *BGRImage) At(x, y int) color.Color {
	if !(image.Point{X: x, Y: y}.In(B.Rect)) {
		return color.RGBA{}
	}
	i := B.PixOffset(x, y)
	s := B.Pix[i : i+3 : i+3] // Small cap improves performance, see https://golang.org/issue/27857
	return color2.BGR{s[0], s[1], s[2]}
}

func (p *BGRImage) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3
}

func (B *BGRImage) Set(x, y int, c color.Color) {
	if !(image.Point{X: x, Y: y}.In(B.Rect)) {
		return
	}
	i := B.PixOffset(x, y)
	c1 := color2.BGRModel.Convert(c).(color2.BGR)
	s := B.Pix[i : i+3 : i+3] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = c1.B
	s[1] = c1.G
	s[2] = c1.R
}

func NewBGRImage(r image.Rectangle) *BGRImage {
	return &BGRImage{
		Pix:    make([]uint8, 3*r.Dx()*r.Dy()),
		Stride: 3 * r.Dx(),
		Rect:   r,
	}
}

// 翻转图像  并padding 空白 BMP 必须 DWORD 对齐
func ReverseDIB(imageBits []byte, srcWidth, srcHeight int, bitsPerPixel int) []byte {
	var imgWidthByteSz int

	switch bitsPerPixel {
	case 24:
		imgWidthByteSz = srcWidth * 3
	case 8:
		imgWidthByteSz = srcWidth
	case 1:
		imgWidthByteSz = (srcWidth + 7) / 8
	case 4:
		imgWidthByteSz = (srcWidth + 1) / 2
	default:
		imgWidthByteSz = srcWidth * (bitsPerPixel) / 8
	}

	padBytes := 0
	if imgWidthByteSz%4 != 0 {
		padBytes = 4 - (imgWidthByteSz % 4)
	}

	alignedImage := make([]byte, (imgWidthByteSz+padBytes)*srcHeight)

	if alignedImage != nil {
		for i := srcHeight - 1; i >= 0; i-- {
			copy(alignedImage[i*(imgWidthByteSz+padBytes):], imageBits[i*imgWidthByteSz:(i+1)*imgWidthByteSz])
		}
		return alignedImage
	}
	return nil
}
