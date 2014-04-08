package captcha

import (
	"image"
	"captcha/cv"
)

type Decoder struct {
	ImageProcessors []cv.ImageProcessor
	BiColorProcessor cv.BiColorProcessor
	BinaryImageProcessors []cv.BinaryImageProcessor
	BinaryImagePredictor BinaryImagePredictor
	ImagePredictor ImagePredictor
}

func (self *Decoder) Predict(img image.Image, mki *MaskIndex, chType int) []*Result {
	if self.ImageProcessors != nil {
		for _, p := range self.ImageProcessors{
			img = p.Process(img)
		}
	}
	if self.BiColorProcessor != nil {
		biImg := self.BiColorProcessor.Process(img)
		if self.BinaryImageProcessors != nil {
			for _, p := range self.BinaryImageProcessors {
				biImg = p.Process(biImg)
			}
		}
		biImg.Save("bi")
		return self.BinaryImagePredictor.Guess(biImg, mki, chType)
	} else {
		return self.ImagePredictor.Guess(img, mki, chType)
	}
}

