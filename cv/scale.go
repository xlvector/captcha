package cv

type ScaleBinaryImage struct {
	Height int
}

func (self *ScaleBinaryImage) Process(img *BinaryImage) *BinaryImage {
	height := self.Height
	if img.Height <= 0 || img.Width <= 0 {
		return nil
	}
	width := img.Width * height / img.Height
	ret := NewBinaryImage(width, height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			xx1 := x * img.Width / width
			xx2 := (x + 1) * img.Width / width
			yy1 := y * img.Height / height
			yy2 := (y + 1) * img.Height / height
			ave := 0
			n := 0
			for i := xx1; i < xx2 && i < img.Width; i++ {
				for j := yy1; j < yy2 && j < img.Height; j++ {
					ave += img.Get(i, j)
					n += 1
				}
			}
			if n == 0 && img.IsOpen(xx1, yy1) {
				ret.Open(x, y)
			}
			if n > 0 && ave*2 >= n {
				ret.Open(x, y)
			}
		}
	}
	return ret
}
