package cv

import (
	"log"
)

func Neighbor8Encoding(img *BinaryImage, x, y int) int {
	ret := 0
	dx := []int{1,1,0,-1,-1,-1,0,1}
	dy := []int{0,-1,-1,-1,0,1,1,1}

	for i := 0; i < 8; i++{
		xx := x + dx[i]
		yy := y + dy[i]
		ret *= 2
		if img.IsOpen(xx, yy) {
			ret += 1
		}
	}
	return ret
}

func Neighbor8Values(img *BinaryImage, x, y int) []int {
	ret := []int{}
	dx := []int{0,1,1,0,-1,-1,-1,0,1}
	dy := []int{0,0,-1,-1,-1,0,1,1,1}

	for i := 0; i < 9; i++{
		xx := x + dx[i]
		yy := y + dy[i]
		ret = append(ret, img.Get(xx, yy))
	}
	return ret
}

func RemoveLineWithNPixelOpen(img *BinaryImage, n int) *BinaryImage{
	for x := 0; x < img.Width; x++ {
		sum := 0
		for y := 0; y < img.Height; y++{
			sum += img.Get(x, y)
		}
		if sum <= n {
			for y := 0; y < img.Height; y++{
				img.Close(x, y)
			}
		} 
	}
	return img
}

type Thining struct {

}

func (self *Thining) processPoint(img *BinaryImage, x, y int) (int) {
	a := Neighbor8Values(img, x, y)
	sigma := a[1] + a[2] + a[3] + a[4] + a[5] + a[6] + a[7] + a[8]
	chi := 0
	for i := 1; i < 8; i += 2 {
		if a[i] != a[(i+2) % 8] {
			chi += 1
		}
	}
	if a[1] == 0 && a[2] == 1 && a[3] == 0 {
		chi += 2
	}
	if a[3] == 0 && a[4] == 1 && a[5] == 0 {
		chi += 2
	}
	if a[5] == 0 && a[6] == 1 && a[7] == 0 {
		chi += 2
	}
	if a[7] == 0 && a[8] == 1 && a[1] == 0 {
		chi += 2
	}
	if a[0] == 1 && chi == 2 && sigma != 1 {
		return 0
	}
	return 1
}

func (self *Thining) Process(img *BinaryImage) *BinaryImage {
	for {
		finish := true
		for x := 0; x < img.Width; x++ {
			for y := 0; y < img.Height; y++{
				if img.IsOpen(x, y) {
					v := self.processPoint(img, x, y)
					if v == 0 {
						img.Close(x, y)
						finish = false
					}
				}
			}
		}
		if finish {
			break
		}
	}
	log.Println("finish")
	return img
}

type Erosion struct{
	Mask []Point
}

func (self *Erosion) Process(img *BinaryImage) *BinaryImage{
	ret := NewBinaryImage(img.Width, img.Height)
	for x := 0; x < img.Width; x++{
		for y := 0; y < img.Height; y++{
			if !img.IsOpen(x, y) {
				continue
			}
			ok := true
			for _, p := range self.Mask {
				if !img.IsOpen(x + p.X, y + p.Y){
					ok = false
					break
				}
			}
			if ok {
				ret.Open(x, y)
			}
		}
	}
	return ret
}