package captcha

import (
	"captcha/cv"
	"log"
	"sort"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Mask struct {
	Img   *cv.BinaryImage
	Label string
}

type MaskIndex struct {
	masks []*Mask
	index map[int][]int
}

func NewMaskIndex() *MaskIndex {
	ret := MaskIndex{}
	ret.masks = []*Mask{}
	ret.index = make(map[int][]int)
	return &ret
}

func (self *MaskIndex) Clear() {
	self.masks = []*Mask{}
	self.index = make(map[int][]int)
}

func (self *Mask) MatchChType(chType int) bool {
	if chType == MIX {
		return true
	} else if chType == NUMBER {
		return self.Label >= "0" && self.Label <= "9"
	} else if chType == ALPHA {
		return self.Label >= "a" && self.Label <= "z"
	}
	return false
}

func (self *MaskIndex) Masks() []*Mask {
	return self.masks
}

func (self *MaskIndex) AddMask(mk *Mask) {
	self.masks = append(self.masks, mk)
	h := mk.Img.FeatureEncode()
	_, ok := self.index[h]
	if !ok {
		self.index[h] = []int{}
	}
	self.index[h] = append(self.index[h], len(self.masks)-1)
}

func (self *MaskIndex) QueryMask(img *cv.BinaryImage) ([]int, bool) {
	h := img.FeatureEncode()
	ret, ok := self.index[h]
	return ret, ok
}

func (self *MaskIndex) FindBestMatchedMasks(img *cv.BinaryImage, chType int, widthWeight int) []*Result {
	h := img.FeatureEncode()
	ids, ok := self.index[h]
	if !ok {
		return nil
	}

	results := ResultSorter{}
	for _, i := range ids {
		mk := self.masks[i]
		if !mk.MatchChType(chType) {
			continue
		}
		sim := mk.Match(img, 0, 0)
		sim *= float64(img.Width)
		sim /= float64(img.Width + widthWeight)
		results = append(results, NewResult(mk.Label, sim))
	}
	sort.Sort(results)
	return results
}

func (self *Mask) Match(src *cv.BinaryImage, srcLeft int, srcTop int) float64 {
	ret := 0.0
	for x := srcLeft; x < srcLeft+self.Img.Width && x < src.Width; x++ {
		for y := srcTop; y < srcTop+self.Img.Height && y < src.Height; y++ {
			if src.IsOpen(x, y) && self.Img.IsOpen(x-srcLeft, y-srcTop) {
				ret += 1.0
			}
		}
	}
	ret *= 2.0
	ret /= float64(self.Img.FrontSize() + src.FrontSize() + 1)
	return ret
}

func (self *Mask) ToString() string {
	ret := self.Label
	ret += "\t"
	ret += strconv.Itoa(self.Img.Width)
	ret += "\t"
	ret += strconv.Itoa(self.Img.Height)
	for x := 0; x < self.Img.Width; x++{
		for y := 0; y < self.Img.Height; y++{
			if !self.Img.IsOpen(x, y) {
				continue
			}
			ret += "\t"
			ret += strconv.Itoa(x)
			ret += ":"
			ret += strconv.Itoa(y)
		}
	}
	return ret
}

func (self *Mask) FromString(line string) {
	tks := strings.Split(line, "\t")
	self.Label = tks[0]
	width, _ := strconv.Atoi(tks[1])
	height, _ := strconv.Atoi(tks[2])
	self.Img = cv.NewBinaryImage(width, height)
	for i := 3; i < len(tks); i++ {
		kv := strings.Split(tks[i], ":")
		x, _ := strconv.Atoi(kv[0])
		y, _ := strconv.Atoi(kv[1])
		self.Img.Open(x, y)
	}
}

func NewMask(line string) *Mask {
	ret := Mask{}
	ret.FromString(line)
	return &ret
}

func LoadMasks(folder string) *MaskIndex {
	scaler := cv.ScaleBinaryImage{Height : SCALE_HEIGHT}
	bounder := cv.BoundBinaryImage{}
	ret := NewMaskIndex()
	n := 0
	err := filepath.Walk(folder, func(root string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			label := f.Name()
			if len(label) > 2{
				return nil
			}
			log.Println(label, root)
			filepath.Walk(root, func(root2 string, f2 os.FileInfo, err2 error) error {
				if f2 == nil {
					return err
				}
				if f2.IsDir() {
					return nil
				}
				img := LoadImage(root2)
				if img != nil {
					mk := Mask{}
					mk.Img = scaler.Process(bounder.Process(cv.Convert2BinaryImage(img)))
					mk.Label = label
					ret.AddMask(&mk)
					n += 1
				}
				return nil
			})
			return nil
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	log.Println("mask count", n)
	return ret
}
