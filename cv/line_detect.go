package cv

import (
	"container/list"
	"log"
	"math"
)

type PathNode struct {
	path []int
}

func NewPathNode() *PathNode {
	ret := PathNode{}
	ret.path = []int{}
	return &ret
}

func CopyPathNode(node *PathNode) *PathNode {
	ret := PathNode{}
	ret.path = []int{}
	for _, pxy := range node.path {
		ret.path = append(ret.path, pxy)
	}
	return &ret
}

func (self *PathNode) AddPoint(x, y int) {
	self.path = append(self.path, EncodeXY(x, y))
}

func (self *PathNode) GetLastPoint() (int, int) {
	pxy := self.path[len(self.path)-1]
	return DecodeXY(pxy)
}

func Neighbor8Sum(img *BinaryImage, x, y int) int {
	ret := 0
	for i := x - 1; i < x+1; i++ {
		for j := y - 1; j < y+1; j++ {
			ret += img.Get(i, j)
		}
	}
	return ret
}

func DetectLine(img *BinaryImage) {
	ps := []int{}
	for x := 0; x < img.Width; x++ {
		for y := 0; y < img.Height; y++{
			if img.IsOpen(x, y) && Neighbor8Sum(img, x, y) < 4 {
				ps = append(ps, EncodeXY(x, y))
			}
		}
	}

	bestStart := 0
	bestEnd := 0
	maxPointInLine := 0
	for i, start := range ps {
		xs, ys := DecodeXY(start)
		for j, end := range ps {
			if j <= i {
				continue
			}
			xe, ye := DecodeXY(end)
			length := math.Sqrt(float64((xs - xe) * (xs - xe) + (ys - ye) * (ys - ye)))
			pointInLine := 0
			for _, p := range ps {
				x, y := DecodeXY(p)
				l1 := math.Sqrt(float64((xs - x) * (xs - x) + (ys - y) * (ys - y)))
				l2 := math.Sqrt(float64((x - xe) * (x - xe) + (y - ye) * (y - ye)))
				if math.Abs(l1 + l2 - length) < 0.2 {
					pointInLine++
				}
			}
			if maxPointInLine < pointInLine {
				maxPointInLine = pointInLine
				bestStart = start
				bestEnd = end
			}
		}
	}

	ret := CopyBinaryImage(img)
	xs, ys := DecodeXY(bestStart)
	xe, ye := DecodeXY(bestEnd)
	length := math.Sqrt(float64((xs - xe) * (xs - xe) + (ys - ye) * (ys - ye)))
	for x := 0; x < img.Width; x++ {
		for y := 0; y < img.Height; y++{
			if img.IsOpen(x, y) && Neighbor8Sum(img, x, y) < 5 {
				l1 := math.Sqrt(float64((xs - x) * (xs - x) + (ys - y) * (ys - y)))
				l2 := math.Sqrt(float64((x - xe) * (x - xe) + (y - ye) * (y - ye)))
				log.Println(x, y, l1 + l2 - length)
				if math.Abs(l1 + l2 - length) < 0.2 {
					ret.Close(x, y)
				}
			}
		}
	}
	ret.Save("line")
	log.Println(bestStart, bestEnd, maxPointInLine)
}

/*
func ExtractLines(img *BinaryImage) *BinaryImage {
	ends := []int{}
	for x := 0; x < img.Width; x++ {
		for y := 0; y < img.Height; y++ {
			p := EncodeXY(x, y)
			if Neighbor8OpenCount(img, p) == 1 {
				ends = append(ends, p)
			}
		}
	}
	for i, p1 := range ends {
		for j, p2 := range ends {
			if j <= i {
				continue
			}
			if HasLine(img, p1, p2) {

			}
		}
	}
}

func AddLine(img *BinaryImage, p1, p2) *BinaryImage{
	for
}

func HasLine(img *BinaryImage, p1, p2 int) bool {
	x1, y1 := p1
	x2, y2 := p2
	if Abs(x1-x2) < 10 && Abs(y1-y2) < 10 {
		return false
	}
	if Abs(x1-x2) > Abs(y1, y2) {
		var xl, xr, yl, yr int
		if x1 < x2 {
			xl, xr = x1, x2
			yl, yr = y1, y2
		} else {
			xl, xr = x2, x1
			yl, yr = y2, y1
		}
		open := 0.0
		total := 0.0
		for x := xl + 1; x < xr; x++ {
			y = yl + (x-xl)*(yr-yl)/(xr-xl)
			if Neighbor8Sum(x, y) > 0 {
				open += 1.0
			}
			total += 1.0
		}
		if open/total > 0.9 {
			return true
		} else {
			return false
		}
	} else {
		var xb, xt, yb, yt int
		if y1 < y2 {
			yb, yt = y1, y2
			xb, xt = x1, x2
		} else {
			yb, yt = y2, y1
			xt, xb = x2, x1
		}
		open := 0.0
		total := 0.0
		for y := yb + 1; y < yt; y++ {
			x = xb + (xt-xb)*(y-yb)/(yt-yb)
			if Neighbor8Sum(x, y) > 0 {
				open += 1.0
			}
			total += 1.0
		}
		if open/total > 0.9 {
			return true
		} else {
			return false
		}
	}
	return false
}

*/
func ExtractLines(img *BinaryImage) *BinaryImage {
	starts := make(map[int]bool)
	for x := 0; x < img.Width; x++ {
		for y := 0; y < img.Height; y++ {
			p := EncodeXY(x, y)
			if Neighbor8Sum(img, x, y) == 1 {
				starts[p] = true
			}
		}
	}
	ret := CopyBinaryImage(img)
	//ret := NewBinaryImage(img.Width, img.Height)
	for p, _ := range starts {
		x0, y0 := DecodeXY(p)
		paths := BFSPathExtract(img, x0, y0)
		bestPath := -1
		maxProb := 0.0
		maxLength := 0
		for k, path := range paths {
			if len(path.path) < 30 {
				continue
			}
			prob := LineProb(path.path)
			if maxProb < prob {
				maxProb = prob
				bestPath = k
				maxLength = len(path.path)
			}
			if maxProb == prob && maxLength < len(path.path) {
				maxLength = len(path.path)
				bestPath = k
			}
		}
		if bestPath >= 0 && maxProb > 0.98 {
			log.Println(maxProb, len(paths[bestPath].path))
			for _, pxy := range paths[bestPath].path {
				x, y := DecodeXY(pxy)
				if Neighbor8Sum(img, x, y) < 4 {
					ret.Close(x, y)
				}
			}
		}
	}
	return ret
}

func BFSPathExtract(img *BinaryImage, x0, y0 int) []*PathNode {
	visited := make(map[int]bool)
	stack := list.New()
	root := NewPathNode()
	root.AddPoint(x0, y0)
	stack.PushBack(root)
	paths := []*PathNode{}
	for {
		if stack.Len() == 0 {
			break
		}
		node := stack.Front().Value.(*PathNode)
		stack.Remove(stack.Front())
		x, y := node.GetLastPoint()
		pxy := EncodeXY(x, y)
		_, ok := visited[pxy]
		if ok {
			continue
		}
		visited[pxy] = true
		end := true
		for i := Max(0, x-1); i <= Min(x+1, img.Width-1); i++ {
			for j := Max(0, y-1); j <= Min(y+1, img.Height-1); j++ {
				if i == x && j == y {
					continue
				}

				if i != x && j != y {
					continue
				}

				if !img.IsOpen(i, j) {
					continue
				}
				_, sok := visited[EncodeXY(i, j)]
				if !sok {
					nextNode := CopyPathNode(node)
					nextNode.AddPoint(i, j)
					stack.PushBack(nextNode)
					end = false
				}
			}
		}
		for i := Max(0, x-1); i <= Min(x+1, img.Width-1); i++ {
			for j := Max(0, y-1); j <= Min(y+1, img.Height-1); j++ {
				if i == x && j == y {
					continue
				}

				if i == x || j == y {
					continue
				}

				if img.Get(i, j) <= 0 {
					continue
				}
				_, sok := visited[EncodeXY(i, j)]
				if !sok {
					nextNode := CopyPathNode(node)
					nextNode.AddPoint(i, j)
					stack.PushBack(nextNode)
					end = false
				}
			}
		}

		if end {
			paths = append(paths, node)
		}
	}
	return paths
}

func LineProb(path []int) float64 {
	if len(path) < 3 {
		return 0.0
	}
	x0, y0 := DecodeXY(path[0])
	x2, y2 := DecodeXY(path[len(path)-1])
	dx02 := float64(x2 - x0)
	dy02 := float64(y2 - y0)
	minRatio := 1.0
	l3 := math.Sqrt(dx02*dx02 + dy02*dy02)
	for i, p := range path {
		if i == 0 || i == len(path)-1 {
			continue
		}
		x1, y1 := DecodeXY(p)
		dx01 := float64(x1 - x0)
		dy01 := float64(y1 - y0)
		dx12 := float64(x2 - x1)
		dy12 := float64(y2 - y1)

		l1 := math.Sqrt(dx01*dx01 + dy01*dy01)
		l2 := math.Sqrt(dx12*dx12 + dy12*dy12)
		ratio := l3 / (l1 + l2 + 0.01)
		if minRatio > ratio {
			minRatio = ratio
		}
	}
	return minRatio
}
