package cv

import (
	"image"
	"image/color"
	"sort"
)

type UInt32Sorter []uint32

func (ms UInt32Sorter) Len() int {
	return len(ms)
}

func (ms UInt32Sorter) Less(i, j int) bool {
	return ms[i] < ms[j]
}

func (ms UInt32Sorter) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func MiddleValue(a UInt32Sorter) uint16 {
	sort.Sort(a)
	return uint16(a[len(a)/2])
}

type MeanShift struct {
	K int
}

func (self *MeanShift) Process(img image.Image) image.Image {
	for k := 0; k < self.K; k++ {
		next := image.NewNRGBA64(img.Bounds())
		for x := 0; x < img.Bounds().Dx(); x++ {
			for y := 0; y < img.Bounds().Dy(); y++ {
				rs := []uint32{}
				gs := []uint32{}
				bs := []uint32{}
				as := []uint32{}
				for i := Max(0, x-1); i <= Min(img.Bounds().Dx()-1, x+1); i++ {
					for j := Max(0, y-1); j <= Min(img.Bounds().Dy()-1, y+1); j++ {
						r, g, b, a := img.At(i, j).RGBA()
						rs = append(rs, r)
						gs = append(gs, g)
						bs = append(bs, b)
						as = append(as, a)
					}
				}
				next.Set(x, y, color.RGBA64{R: MiddleValue(rs), G: MiddleValue(gs), B: MiddleValue(bs), A: MiddleValue(as)})
			}
		}
		img = next
	}
	return img
}

type AverageSmooth struct {

}

func (self *AverageSmooth) Process(img image.Image) image.Image {
	ret := image.NewNRGBA64(img.Bounds())
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			rs := 0
			gs := 0
			bs := 0
			as := 0
			n := 0
			for i := Max(0, x-1); i <= Min(img.Bounds().Dx()-1, x+1); i++ {
				for j := Max(0, y-1); j <= Min(img.Bounds().Dy()-1, y+1); j++ {
					c := 1
					if i == x && j == y {
						c = 8
					}
					r, g, b, a := img.At(i, j).RGBA()
					rs += int(r) * c
					gs += int(g) * c
					bs += int(b) * c
					as += int(a) * c
					n += c
				}
			}
			ret.Set(x, y, color.RGBA64{R: uint16(rs/n), G: uint16(gs/n), B: uint16(bs/n), A: uint16(as/n)})
		}
	}
	return ret
}

type BinaryImageMeanSmooth struct {
	Threshold int
}

func (self * BinaryImageMeanSmooth) Process(img *BinaryImage) *BinaryImage {
	ret := NewBinaryImage(img.Width, img.Height)
	for x := 0; x < img.Width; x++ {
		for y := 0; y < img.Height; y++ {
			total := 0
			n := 0
			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					if i == 0 && j == 0 {
						break
					}
					total += img.Get(x+i, y+j)
					n += 1
				}
			}
			if total > self.Threshold {
				ret.Open(x, y)
			}
		}
	}
	return ret
}

type RemoveIsolatePoints struct {

}

func (self * RemoveIsolatePoints) Process(img *BinaryImage) *BinaryImage {
	ret := CopyBinaryImage(img)
	for x := 0; x < img.Width; x++ {
		for y := 0; y < img.Height; y++ {
			if !img.IsOpen(x, y) {
				continue
			}
			total := 0
			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					if i == 0 && j == 0 {
						continue
					}
					total += img.Get(x+i, y+j)
				}
			}
			if total == 0 {
				ret.Close(x, y)
			}
		}
	}
	return ret
}

type RemoveXAxis struct {
	K int
}

func (self * RemoveXAxis) Process(img *BinaryImage) *BinaryImage {
	ret := NewBinaryImage(img.Width, img.Height)
	for x := 0; x < img.Width; x++ {
		sum := 0
		for y := 0; y < img.Height; y++ {
			sum += img.Get(x, y)
		}
		if sum > self.K {
			for y := 0; y < img.Height; y++{
				if img.IsOpen(x, y){
					ret.Open(x, y)
				}
			}
		}
	}
	return ret
}

