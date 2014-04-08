package cv

import (
	"container/list"
	"image"
)

func ExtractAllConnectedComponentInBinaryImage(img *BinaryImage, dx, dy int) []*BinaryImage {
	visited := make(map[int]bool)
	ret := []*BinaryImage{}
	for x := 0; x < img.Width; x++ {
		for y := 0; y < img.Height; y++ {
			if img.IsOpen(x, y) {
				_, ok := visited[EncodeXY(x, y)]
				if ok {
					continue
				}
				var cc *BinaryImage
				cc, visited = extractConnectedComponentInBinaryImage(img, x, y, visited, dx, dy)
				ret = append(ret, cc)
			}
		}
	}
	return ret
}

func HoleNumberInBinaryImage(img *BinaryImage) int {
	visited := make(map[int]bool)
	ret := 0
	for x := 0; x < img.Width; x++ {
		for y := 0; y < img.Height; y++ {
			if !img.IsOpen(x, y) {
				_, ok := visited[EncodeXY(x, y)]
				if ok {
					continue
				}
				_, visited = extractConnectedComponentInBinaryImage(img, x, y, visited, 1, 1)
				isHole := true
				for p, _ := range visited {
					x0, y0 := DecodeXY(p)
					if x0 == 0 || y0 == 0 || x0 == img.Width - 1 || y0 == img.Height - 1{
						isHole = false
						break
					}
				}
				if isHole {
					ret++
				}
			}
		}
	}
	return ret
}

func extractConnectedComponentInBinaryImage(img *BinaryImage, x0, y0 int, visited map[int]bool, dx, dy int) (*BinaryImage, map[int]bool) {
	q := list.New()
	q.PushBack(EncodeXY(x0, y0))
	subVisited := []int{}
	for {
		if q.Len() == 0 {
			break
		}
		p := q.Back().Value.(int)
		x, y := DecodeXY(p)
		q.Remove(q.Back())

		for xx := Max(0, x-dx); xx <= x+dx && xx < img.Width; xx++ {
			for yy := Max(0, y-dy); yy <= y+dy && yy < img.Height; yy++{
				if img.Get(xx, yy) == img.Get(x0, y0) {
					pxy := EncodeXY(xx, yy)
					_, ok := visited[pxy]
					if !ok {
						q.PushBack(pxy)
						visited[pxy] = true
						subVisited = append(subVisited, pxy)
					}
				}
			}
		}

	}

	ret := NewBinaryImage(img.Width, img.Height)
	for _, pxy := range subVisited {
		xx, yy := DecodeXY(pxy)
		ret.Open(xx, yy)
	}
	bounder := BoundBinaryImage{}
	return bounder.Process(ret), visited
}

func Color(img image.Image, x, y int) (int, int, int){
	r, g, b, _ := img.At(x, y).RGBA()
	return int(r/255), int(g/255), int(b/255)
}

func NeighborColorDiff(img image.Image, x, y int) int {
	minR := 255
	minG := 255
	minB := 255
	maxR := 0
	maxG := 0
	maxB := 0
	//for i := Max(0, x - 1); i <= x +1 && i < img.Bounds().Dx(); i++{
	//	for j := Max(0, y - 1); j <= y + 1 && j < img.Bounds().Dy(); j++{
	for i := x; i <= x +1 && i < img.Bounds().Dx(); i++{
		for j := y; j <= y + 1 && j < img.Bounds().Dy(); j++{
			r, g, b := Color(img, i, j)
			minR = Min(minR, r)
			minG = Min(minG, g)
			minB = Min(minB, b)
			maxR = Max(maxR, r)
			maxG = Max(maxG, g)
			maxB = Max(maxB, b)
		}
	}
	return Max3((maxR - minR), (maxG - minG), (maxB - minB))
}

func extractConnectedComponent(img image.Image, x0, y0 int, visited map[int]bool, dx, dy int) (*BinaryImage, map[int]bool) {
	q := list.New()
	q.PushBack(EncodeXY(x0, y0))
	subVisited := []int{}
	r0, g0, b0 := Color(img, x0, y0)
	for {
		if q.Len() == 0 {
			break
		}
		p := q.Back().Value.(int)
		x, y := DecodeXY(p)
		q.Remove(q.Back())
		
		for xx := Max(0, x-dx); xx <= x+dx && xx < img.Bounds().Dx(); xx++ {
			for yy := Max(0, y-dy); yy <= y+dy && yy < img.Bounds().Dy(); yy++{
				r1, g1, b1 := Color(img, xx, yy)
				if Max3(Abs(r0 - r1), Abs(g0 - g1), Abs(b0 - b1)) < 60 {
					pxy := EncodeXY(xx, yy)
					_, ok := visited[pxy]
					if !ok {
						q.PushBack(pxy)
						visited[pxy] = true
						subVisited = append(subVisited, pxy)
					}
				}
			}
		}

	}

	ret := NewBinaryImage(img.Bounds().Dx(), img.Bounds().Dy())
	border := false
	left := 1000
	right := 0
	top := 10000
	bottom := 0
	for _, pxy := range subVisited {
		xx, yy := DecodeXY(pxy)
		if xx == 0 || yy == 0 || xx == img.Bounds().Dx() - 1 || yy == img.Bounds().Dy() - 1{
			border = true
		}
		left = Min(left, xx)
		right = Max(right, xx)
		top = Min(top, yy)
		bottom = Max(bottom, yy)
		ret.Open(xx, yy)
	}
	if border && ((right - left) < 6 || (bottom - top) < 6){
		return nil, visited
	}
	bounder := BoundBinaryImage{}
	return bounder.Process(ret), visited
}

func ExtractAllConnectedComponent(img image.Image, dx, dy int) []*BinaryImage {
	visited := make(map[int]bool)
	ret := []*BinaryImage{}
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			_, ok := visited[EncodeXY(x, y)]
			if ok {
				continue
			}
			if NeighborColorDiff(img, x, y) > 15 {
				continue
			}
			var cc *BinaryImage
			cc, visited = extractConnectedComponent(img, x, y, visited, dx, dy)
			if cc != nil{
				ret = append(ret, cc)
			}
		}
	}
	return ret
}
