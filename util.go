package captcha

import (
	"captcha/cv"
)

func FindBeginPos(img *cv.BinaryImage, beginX int) (int, int) {
	for x := beginX; x < img.Width; x++ {
		for y := 0; y < img.Height; y++ {
			if img.IsOpen(x, y) {
				return x, y
			}
		}
	}
	return img.Width - 1, img.Height - 1
}

func FindEndPos(img *cv.BinaryImage, beginX int) int {
	for x := beginX; x < img.Width; x++ {
		sum := 0
		for y := 0; y < img.Height; y++ {
			sum += img.Get(x, y)
		}
		if sum == 0 {
			return x
		}
	}
	return img.Width
}
