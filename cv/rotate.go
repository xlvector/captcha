package cv

import (
	"math"
)

func Rotate(img *BinaryImage, theta float64) *BinaryImage {
	left := 1000.0
	top := 1000.0
	right := -1000.0
	bottom := -1000.0
	for x := 0; x < img.Width; x++ {
		for y := 0; y < img.Height; y++ {
			if !img.IsOpen(x, y) {
				continue
			}
			xx := float64(x)
			yy := float64(y)
			rx := xx*math.Cos(theta) - yy*math.Sin(theta)
			ry := xx*math.Sin(theta) + yy*math.Cos(theta)
			left = math.Min(left, rx)
			right = math.Max(right, rx)
			top = math.Min(top, ry)
			bottom = math.Max(bottom, ry)
		}
	}

	width := int(right - left + 1)
	height := int(bottom - top + 1)
	ret := NewBinaryImage(width, height)
	for x := 0; x < img.Width; x++ {
		for y := 0; y < img.Height; y++ {
			if !img.IsOpen(x, y) {
				continue
			}
			for xx := float64(x) - 0.5; xx <= float64(x)+0.5; xx += 0.5 {
				for yy := float64(y) - 0.5; yy <= float64(y)+0.5; yy += 0.5 {
					rx := xx*math.Cos(theta) - yy*math.Sin(theta)
					ry := xx*math.Sin(theta) + yy*math.Cos(theta)
					ret.Open(int(rx-left+0.5), int(ry-top+0.5))
				}
			}
		}
	}
	return ret
}
