package cv

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
)

type BinaryImage struct {
	Width, Height, Size int
	Data                []byte
}

func NewBinaryImage(width, height int) *BinaryImage {
	ret := BinaryImage{}
	ret.Width = width
	ret.Height = height
	ret.Size = width * height
	ret.Data = []byte{}
	for i := 0; i < ret.Size; i++ {
		ret.Data = append(ret.Data, byte(0))
	}
	return &ret
}

func CopyBinaryImage(src *BinaryImage) *BinaryImage {
	ret := BinaryImage{}
	ret.Width = src.Width
	ret.Height = src.Height
	ret.Size = src.Size
	ret.Data = []byte{}
	for i := 0; i < ret.Size; i++ {
		ret.Data = append(ret.Data, src.Data[i])
	}
	return &ret
}

func (self *BinaryImage) Save(name string) {
	if self.Width <= 0 || self.Height <= 0 {
		return
	}
	ret := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{self.Width, self.Height}})
	for x := 0; x < self.Width; x++ {
		for y := 0; y < self.Height; y++ {
			if self.IsOpen(x, y) {
				ret.SetGray(x, y, color.Gray{Y: 255})
			} else {
				ret.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}
	f, err := os.Create(name + ".png")
	if err != nil {
		return
	}
	defer f.Close()
	png.Encode(f, ret)
}

func (self *BinaryImage) Open(x, y int) {
	if x < 0 || y < 0 || x >= self.Width || y >= self.Height {
		return
	}
	self.Data[y*self.Width+x] = 1
}

func (self *BinaryImage) Close(x, y int) {
	if x < 0 || y < 0 || x >= self.Width || y >= self.Height {
		return
	}
	self.Data[y*self.Width+x] = 0
}

func (self *BinaryImage) IsOpen(x, y int) bool {
	if x < 0 || y < 0 || x >= self.Width || y >= self.Height {
		return false
	}
	if self.Data[y*self.Width+x] == 1 {
		return true
	} else {
		return false
	}
}

func (self *BinaryImage) Get(x, y int) int {
	if self.IsOpen(x, y) {
		return 1
	} else {
		return 0
	}
}

func (self *BinaryImage) Encode() string {
	n := int64(0)
	for _, b := range self.Data {
		n *= 11
		n += int64(b)
	}
	if n < 0 {
		n *= -1
	}
	return strconv.FormatInt(n, 10)
}

func (self *BinaryImage) FrontSize() int {
	ret := 0
	if self.Data == nil {
		return ret
	}
	for _, b := range self.Data {
		ret += int(b)
	}
	return ret
}

func (self *BinaryImage) FeatureEncode() int {
	/*left := 0
	right := 0
	for x := 0; x * 2 < self.Width; x++{
		for y := 0; y < self.Height; y++{
			if self.IsOpen(x, y){
				left += 1
			}
			if self.IsOpen(self.Width - x - 1, y){
				right += 1
			}
		}
	}
	ret := 0
	if left < right / 2 {
		ret = 0
	} else if right < left / 2 {
		ret = 1
	} else {
		ret = 2
	}

	top := 0
	bottom := 0
	for x := 0; x < self.Width; x++{
		for y := 0; y * 2 < self.Height; y++{
			if self.IsOpen(x, y){
				top += 1
			}
			if self.IsOpen(x, self.Height - 1 - y){
				bottom += 1
			}
		}
	}
	ret *= 10
	if top < bottom / 2 {
		ret += 0
	} else if bottom < top / 2 {
		ret += 1
	} else {
		ret += 2
	}

	ret *= 10
	if self.Width * 4 < self.Height {
		ret += 1
	} else {
		ret += 2
	}
	
	ret *= 10
	ret += self.FrontSize() * 3 / self.Size

	
	maxX := 0
	meanX := 0
	for x := 0; x < self.Width; x++{
		num := 0
		for y := 0; y < self.Height; y++{
			num += self.Get(x, y)
		}
		maxX = Max(maxX, num)
		meanX += num
	}
	meanX /= self.Width
	ret *= 10
	ret += (maxX + 1) / (meanX + 1)
	*/
	return 0
}

func (self *BinaryImage) ExtractFeatures() []float64 {
	ret := []float64{}
	ret = append(ret, float64(self.Width)/float64(self.Width+self.Height))

	for x := 0; x < 32; x++ {
		for y := 0; y < 32; y++ {
			xx1 := x * self.Width / 32
			xx2 := (x + 1) * self.Width / 32
			yy1 := y * self.Height / 32
			yy2 := (y + 1) * self.Height / 32
			ave := 0
			n := 0
			for i := xx1; i < xx2 && i < self.Width; i++ {
				for j := yy1; j < yy2 && j < self.Height; j++ {
					ave += self.Get(i, j)
					n += 1
				}
			}
			open := 0.0
			if n == 0 && self.IsOpen(xx1, yy1) {
				open = 1.0
			} else {
				open = float64(ave) / float64(n)
			}
			ret = append(ret, open)
		}
	}

	return ret
}
