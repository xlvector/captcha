package cv

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