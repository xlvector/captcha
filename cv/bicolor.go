package cv

import (
	"image"
	"log"
	"math"
)

func Convert2GrayImage(img image.Image) *GrayImage {
	dx := img.Bounds().Dx()
	dy := img.Bounds().Dy()
	gray := NewGrayImage(dx, dy)

	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			c := (r + g + b) / (255 * 3)
			if c > 255 {
				c = 255
			}
			if c < 0 {
				c = 0
			}
			gray.Set(x, y, byte(c))
		}
	}
	return gray
}

func Convert2BinaryImage(img image.Image) *BinaryImage {
	dx := img.Bounds().Dx()
	dy := img.Bounds().Dy()

	ret := NewBinaryImage(dx, dy)
	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			c := (r + g + b) / (255 * 3)
			if c > 200 {
				ret.Open(x, y)
			}
		}
	}
	return ret
}

type ColorCount struct {
	Color int
	Count float64
}

type PeakAverageBasedBiColor struct {

}

func (self *PeakAverageBasedBiColor) Process(img image.Image) *BinaryImage {
	gray := Convert2GrayImage(img)

	his := [256]float64{}
	for x := 0; x < gray.Width; x++ {
		for y := 0; y < gray.Height; y++ {
			his[gray.Get(x, y)] += 1.0
		}
	}
	avePeakColor := 0
	nPeak := 0
	for i := 0; i < len(his); i++ {
		peak := true
		for j := Max(0, i-25); j < Min(i+25, len(his)); j++ {
			if i == j {
				continue
			}
			if his[i] <= his[j] {
				peak = false
				break
			}
		}

		if peak {
			avePeakColor += i
			nPeak += 1
		}
	}
	if nPeak > 1 {
		avePeakColor /= nPeak
	}
	left := 0.0
	right := 0.0
	for i, c := range his {
		if i < avePeakColor {
			left += c
		} else {
			right += c
		}
	}
	log.Println("PeakAverageBasedBiColor", avePeakColor)
	ret := NewBinaryImage(gray.Width, gray.Height)
	for x := 0; x < gray.Width; x++ {
		for y := 0; y < gray.Height; y++ {
			if left < right && gray.Get(x, y) < avePeakColor {
				ret.Open(x, y)
			} else if left >= right && gray.Get(x, y) >= avePeakColor {
				ret.Open(x, y)
			}
		}
	}
	return ret
}

type MaxColorBasedBiColor struct {

}

func (self *MaxColorBasedBiColor) Process(img image.Image) *BinaryImage {
	gray := Convert2GrayImage(img)

	his := [256]float64{}
	for x := 0; x < gray.Width; x++ {
		for y := 0; y < gray.Height; y++ {
			his[gray.Get(x, y)] += 1.0
		}
	}
	left := 0
	maxSum := 0.0
	for i := 0; i < 256; i++ {
		sum := 0.0
		for j := i; j < i+50; j++ {
			sum += his[i]
		}
		if maxSum < sum {
			maxSum = sum
			left = i
		}
	}
	log.Println("MaxColorBasedBiColor", left)
	ret := NewBinaryImage(gray.Width, gray.Height)
	for x := 0; x < gray.Width; x++ {
		for y := 0; y < gray.Height; y++ {
			if gray.Get(x, y) < left || gray.Get(x, y) > left+50 {
				ret.Open(x, y)
			}
		}
	}
	return ret
}

type EntropyBasedBiColor struct {

}

func (self *EntropyBasedBiColor) Process(img image.Image) *BinaryImage {
	gray := Convert2GrayImage(img)

	his := [256]float64{}
	p := [256]float64{}
	total := 0.0
	for x := 0; x < gray.Width; x++ {
		for y := 0; y < gray.Height; y++ {
			his[gray.Get(x, y)] += 1.0
			p[gray.Get(x, y)] += 1.0
			total += 1.0
		}
	}
	for i := 0; i < 256; i++ {
		his[i] /= total
		p[i] /= total
	}
	for i := 1; i < len(p); i++ {
		p[i] += p[i-1]
	}

	maxHab := 0.0
	bestSplit := 0
	for k := 1; k < 255; k++ {
		if p[k-1] < 0.001 || p[k-1] > 0.999 {
			continue
		}
		ha := 0.0
		for i := 0; i < k; i++ {
			if his[i] < 1e-8 {
				continue
			}
			ha -= (his[i] / p[k-1]) * math.Log(his[i]/p[k-1])
		}
		hb := 0.0
		for i := k; i < 256; i++ {
			if his[i] < 1e-8 {
				continue
			}
			hb -= (his[i] / (1.0 - p[k-1])) * math.Log(his[i]/(1.0-p[k-1]))
		}
		if maxHab < ha+hb {
			maxHab = ha + hb
			bestSplit = k
		}
	}
	log.Println("EntropyBasedBiColor", maxHab, bestSplit)
	ret := NewBinaryImage(gray.Width, gray.Height)
	for x := 0; x < gray.Width; x++ {
		for y := 0; y < gray.Height; y++ {
			if p[bestSplit] < 0.5 {
				if gray.Get(x, y) < bestSplit {
					ret.Open(x, y)
				}
			} else {
				if gray.Get(x, y) >= bestSplit {
					ret.Open(x, y)
				}
			}
		}
	}
	return ret
}
