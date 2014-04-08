package cv

type RemoveBinaryImageBorder struct {

}
func (self *RemoveBinaryImageBorder) Process(img *BinaryImage) *BinaryImage {
	left := 0
	top := 0
	right := img.Width
	bottom := img.Height
	for x := 0; x < img.Width; x++ {
		sum := 0
		for y := 0; y < img.Height; y++ {
			sum += img.Get(x, y)
		}
		if sum < img.Height-6 {
			left = x
			break
		}
	}
	for x := img.Width - 1; x >= 0; x-- {
		sum := 0
		for y := 0; y < img.Height; y++ {
			sum += img.Get(x, y)
		}
		if sum < img.Height-6 {
			right = x
			break
		}
	}
	for y := 0; y < img.Height; y++ {
		sum := 0
		for x := 0; x < img.Width; x++ {
			sum += img.Get(x, y)
		}
		if sum < img.Width-8 {
			top = y
			break
		}
	}
	for y := img.Height - 1; y >= 0; y-- {
		sum := 0
		for x := 0; x < img.Width; x++ {
			sum += img.Get(x, y)
		}
		if sum < img.Width-8 {
			bottom = y
			break
		}
	}
	ret := NewBinaryImage(right-left+1, bottom-top+1)
	for x := left; x <= right; x++ {
		for y := top; y <= bottom; y++ {
			if img.IsOpen(x, y) {
				ret.Open(x-left, y-top)
			}
		}
	}
	return ret
}

type BoundBinaryImage struct {
	XMinOpen, YMinOpen int
}

func (self *BoundBinaryImage) Process(img *BinaryImage) *BinaryImage {
	left := img.Width
	top := img.Height
	right := 0
	bottom := 0
	hx := []int{}
	hy := []int{}
	for x := 0; x < img.Width; x++{
		hx = append(hx, 0)
	}
	for y := 0; y < img.Height; y++{
		hy = append(hy, 0)
	}
	for x := 0; x < img.Width; x++ {
		for y := 0; y < img.Height; y++ {
			hx[x] += img.Get(x, y)
			hy[y] += img.Get(x, y)
		}
	}
	for x := 0; x < img.Width; x++{
		if hx[x] > self.XMinOpen {
			left = x
			break
		}
	}
	for x := img.Width - 1; x >= 0; x--{
		if hx[x] > self.XMinOpen {
			right = x
			break
		}
	}
	for y := 0; y < img.Height; y++{
		if hy[y] > self.YMinOpen {
			top = y
			break
		}
	}
	for y := img.Height - 1; y >= 0; y--{
		if hy[y] > self.YMinOpen {
			bottom = y
			break
		}
	}
	if left > right || top > bottom {
		return NewBinaryImage(1, 1)
	}
	ret := NewBinaryImage(right-left+1, bottom-top+1)
	for x := left; x <= right; x++ {
		for y := top; y <= bottom; y++ {
			if img.IsOpen(x, y) {
				ret.Open(x-left, y-top)
			}
		}
	}
	return ret
}
