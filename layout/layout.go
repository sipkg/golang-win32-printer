package layout

func CenterElement(pageWidth, elementWidth uint32) uint32 {
	return (pageWidth - elementWidth) / 2
}

func AlignRight(pageWidth, elementWidth uint32) uint32 {
	return pageWidth - elementWidth
}

func AlignLeft() uint32 {
	return 0
}

func AlignRightFrom(startX, elementWidth uint32) uint32 {
	return startX - elementWidth
}

func AlignLeftFrom(startX uint32) uint32 {
	return startX
}

func CenterElementFrom(startX, containerWidth, elementWidth uint32) uint32 {
	return startX + (containerWidth-elementWidth)/2
}

func AlignTopFrom(startY, elementHeight uint32) uint32 {
	return startY
}

func AlignBottomFrom(startY, containerHeight, elementHeight uint32) uint32 {
	return startY + containerHeight - elementHeight
}

func CenterElementVerticallyFrom(startY, containerHeight, elementHeight uint32) uint32 {
	return startY + (containerHeight-elementHeight)/2
}

func TruncateString(str string, maxLength int) string {
	if len(str) > maxLength {
		return str[:maxLength-1] + "."
	}
	return str
}
