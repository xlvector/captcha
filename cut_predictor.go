package captcha

import (
	"captcha/cv"
	"log"
)

type FastCutBasedPredictor struct {

}

func (self *FastCutBasedPredictor) XHis(img *cv.BinaryImage) []int {
	h := []int{}
	for x := 0; x < img.Width; x++{
		h = append(h, 0)
		for y := 0; y < img.Height; y++{
			h[x] += img.Get(x, y)
		}
	}
	return h
}

func (self *FastCutBasedPredictor) FindNextSplitPoints(img *cv.BinaryImage, sx int, h []int) []int {
	ret := []int{}
	for x := cv.Max(1, sx + 1); x < img.Width - 1 && x < sx + img.Height + 10; x++{
		if x == img.Width - 2 {
			ret = append(ret, x + 1)
			continue
		}
		if h[x] == 0 {
			ret = append(ret, x)
			continue
		}
		if (h[x] < 5 && (h[x] < h[x - 1] && h[x] <= h[x+1]) || (h[x] <= h[x - 1] && h[x] < h[x+1])) {
			ret = append(ret, x)
			continue
		}
	}
	return ret
}

func (self *FastCutBasedPredictor) CutMatrixByRect(img *cv.BinaryImage, left, top, right, bottom int) *cv.BinaryImage {
	if right <= left || bottom <= top {
		return nil
	}
	ret := cv.NewBinaryImage(right-left, bottom-top)
	for x := left; x < right && x < img.Width; x++ {
		for y := top; y < bottom && y < img.Height; y++ {
			if img.IsOpen(x, y) {
				ret.Open(x-left, y-top)
			}
		}
	}
	return ret
}

func (self *FastCutBasedPredictor) Guess(img *cv.BinaryImage, mki *MaskIndex, chType int) []*Result {
	h := self.XHis(img)
	scaler := cv.ScaleBinaryImage{Height : SCALE_HEIGHT}
	bounder := cv.BoundBinaryImage{}
	curX := 0
	ret := NewResult("", 0.0)
	failCount := 0
	for {
		if curX >= img.Width {
			break
		}
		curX, _ = FindBeginPos(img, curX)
		
		maxSim := 0.0
		var bestCutImg *cv.BinaryImage
		var bestLocalResults []*Result
		nextX := 0
		endXs := self.FindNextSplitPoints(img, curX, h)
		for beginX := 0; beginX < 3; beginX++ {
			for k, endX := range endXs {
				if k > 20 {
					break
				}
				width := endX - curX - beginX + 1
				if width <= 3 {
					continue
				}

				cutImg := self.CutMatrixByRect(img, curX+beginX, 0, endX, img.Height)
				if cutImg == nil || cutImg.FrontSize() < 10 {
					continue
				}
				cutImg = bounder.Process(cutImg)
				if cutImg.Height < 5 {
					continue
				}
				cutImg = scaler.Process(cutImg)
				if cutImg == nil {
					continue
				}
				localResults := mki.FindBestMatchedMasks(cutImg, chType, 1)
				if localResults == nil || len(localResults) == 0{
					continue
				}
				if maxSim < localResults[0].Weight {
					maxSim = localResults[0].Weight
					bestLocalResults = localResults
					nextX = endX
					bestCutImg = cutImg
				}
			}
			if maxSim > 0.8 {
				break
			}
		}
		if maxSim < 0.6 {
			curX += 1
			failCount += 1
			continue
		}
		curX = nextX
		ret.Label += bestLocalResults[0].Label
		ret.Weight += bestLocalResults[0].Weight
		ret.AddComponent(bestCutImg)
		log.Println(bestLocalResults[0].Label, maxSim)
		if len(ret.Label) > 8 || failCount > 5 {
			return []*Result{}
		}
	}
	return []*Result{ret}
}
