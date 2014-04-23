package cv

import (
	"math"
)

type Point struct {
	X, Y int
}

func EncodeXY(x, y int) int {
	return x*10000 + y
}

func DecodeXY(p int) (int, int) {
	x := p / 10000
	y := p % 10000
	return x, y
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Max3(a, b, c int) int {
	if a > b && a > c {
		return a
	}
	if b > a && b > c {
		return b
	}
	if c > a && c > b {
		return c
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Abs(a int) int {
	if a < 0 {
		return -1 * a
	} else {
		return a
	}
}

func Mean(h []float64) float64 {
	sum := 0.0
	for _, v := range h {
		sum += v
	}
	if len(h) == 0{
		return 0.0
	}
	return sum / float64(len(h))
}

func Var(h []float64) float64 {
	m := Mean(h)
	ret := 0.0
	for _, v := range h {
		ret += (v - m) * (v - m)
	}
	return math.Sqrt(ret / float64(len(h)))
}
