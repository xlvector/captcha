package captcha

import (
	"captcha/cv"
	"log"
	"image"
	_ "strconv"
)

type BinaryImageConnectedComponentPredictor struct {
	Dx, Dy int
}

func (self *BinaryImageConnectedComponentPredictor) Guess(img *cv.BinaryImage, mki *MaskIndex, chType int) []*Result {
	scaler := cv.ScaleBinaryImage{Height : SCALE_HEIGHT}
	ccs := cv.ExtractAllConnectedComponentInBinaryImage(img, self.Dx, self.Dy)
	ret := NewResult("", 0.0)
	if len(ccs) < 4 {
		return []*Result{}
	}
	for _, cc := range ccs {
		if cc.Size < 50 || cc.Width < 3 || cc.Height < 8 {
			continue
		}
		unScaleCC := cv.CopyBinaryImage(cc)
		cc = scaler.Process(cc)
		localResults := mki.FindBestMatchedMasks(cc, chType, 0)
		if localResults == nil || len(localResults) == 0{
			continue
		}
		if localResults[0].Weight < 0.65 {
			continue
		}
		log.Println(localResults[0].Label, localResults[0].Weight)

		ret.Label += localResults[0].Label
		ret.Weight += localResults[0].Weight
		ret.AddComponent(unScaleCC)
		if len(ret.Label) > 8 {
			return []*Result{}
		}
	}
	return []*Result{ret}
}

type ConnectedComponentPredictor struct {
	Dx, Dy int
}

func (self *ConnectedComponentPredictor) Guess(img image.Image, mki *MaskIndex, chType int) []*Result {
	scaler := cv.ScaleBinaryImage{Height : SCALE_HEIGHT}
	ccs := cv.ExtractAllConnectedComponent(img, self.Dx, self.Dy)
	ret := NewResult("", 0.0)
	if len(ccs) < 4 {
		return []*Result{}
	}
	imgSize := img.Bounds().Dx() * img.Bounds().Dy()
	for _, cc := range ccs {
		if cc.Size < 50 || cc.Size + 100 > imgSize || cc.Width + 10 > img.Bounds().Dx() || cc.Width < 3 || cc.Height < 8 {
			continue
		}
		unScaleCC := cv.CopyBinaryImage(cc)
		cc = scaler.Process(cc)
		
		localResults := mki.FindBestMatchedMasks(cc, chType, 0)
		if localResults == nil || len(localResults) == 0{
			continue
		}
		log.Println(localResults[0].Label, localResults[0].Weight)
		if localResults[0].Weight < 0.65 {
			continue
		}
		ret.Label += localResults[0].Label
		ret.Weight += localResults[0].Weight
		ret.AddComponent(unScaleCC)
		if len(ret.Label) > 8 {
			return []*Result{}
		}
	}
	return []*Result{ret}
}
