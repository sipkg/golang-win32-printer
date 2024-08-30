package color

import "image/color"

type BGR struct {
	B, G, R uint8
}

func (B BGR) RGBA() (r, g, b, a uint32) {
	r = uint32(B.R)
	r |= r << 8
	g = uint32(B.G)
	g |= g << 8
	b = uint32(B.B)
	b |= b << 8
	a = 0xFF
	return
}

var (
	BGRModel = color.ModelFunc(func(c color.Color) color.Color {
		if y, ok := c.(color.YCbCr); ok {
			r, g, b := color.YCbCrToRGB(y.Y, y.Cb, y.Cr)
			return BGR{b, g, r}
		}
		r, g, b, _ := c.RGBA()
		return BGR{uint8(b), uint8(g), uint8(r)}
	})
)
