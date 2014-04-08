package cv

import (
	_ "strconv"
)

type GrayImage struct {
	Width, Height, Size int
	Data                []byte
}

func NewGrayImage(width, height int) *GrayImage {
	ret := GrayImage{}
	ret.Width = width
	ret.Height = height
	ret.Size = width * height
	ret.Data = []byte{}
	for i := 0; i < ret.Size; i++ {
		ret.Data = append(ret.Data, byte(0))
	}
	return &ret
}

func CopyGrayImage(src *GrayImage) *GrayImage {
	ret := GrayImage{}
	ret.Width = src.Width
	ret.Height = src.Height
	ret.Size = src.Size
	ret.Data = []byte{}
	for i := 0; i < ret.Size; i++ {
		ret.Data = append(ret.Data, src.Data[i])
	}
	return &ret
}

func (self *GrayImage) Set(x, y int, val byte) {
	if x < 0 || y < 0 || x >= self.Width || y >= self.Height {
		return
	}
	self.Data[y*self.Width+x] = val
}

func (self *GrayImage) Get(x, y int) int {
	if x < 0 || y < 0 || x >= self.Width || y >= self.Height {
		return 0
	}
	return int(self.Data[y*self.Width+x])
}
