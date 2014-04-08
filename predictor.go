package captcha

import (
	"captcha/cv"
	"image"
)

type Result struct {
	Label  string
	Weight float64
	Components []*cv.BinaryImage
}

func NewResult(label string, weight float64) *Result {
	ret := Result{}
	ret.Label = label
	ret.Weight = weight
	ret.Components = []*cv.BinaryImage{}
	return &ret
}

func (self *Result) GetWeight() float64 {
	return self.Weight / float64(len(self.Label))
}

func (self *Result) AddComponent(img *cv.BinaryImage){
	self.Components = append(self.Components, img)
}

type ResultSorter []*Result

func (ms ResultSorter) Len() int {
	return len(ms)
}

func (ms ResultSorter) Less(i, j int) bool {
	return ms[i].GetWeight() > ms[j].GetWeight()
}

func (ms ResultSorter) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func GenerateResults(results [][]*Result) []*Result {
	if len(results) <= 0 {
		return nil
	} else if len(results) == 1 {
		return results[0]
	} else {
		rights := GenerateResults(results[1:])
		ret := []*Result{}
		for _, prefix := range results[0] {
			for _, right := range rights{
				ret = append(ret, NewResult(prefix.Label + right.Label, prefix.Weight + right.Weight))
			}
		}
		return ret
	}
}

type BinaryImagePredictor interface {
	Guess(img *cv.BinaryImage, mki *MaskIndex, chType int) []*Result
}

type ImagePredictor interface {
	Guess(img image.Image, mki *MaskIndex, chType int) []*Result
}
