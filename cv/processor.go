package cv

import (
	"image"
)

type ImageProcessor interface {
	Process(img image.Image) image.Image
}

type BinaryImageProcessor interface {
	Process(img *BinaryImage) *BinaryImage
}

type BiColorProcessor interface {
	Process(img image.Image) *BinaryImage
}