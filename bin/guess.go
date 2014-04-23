package main

import (
	"captcha"
	"captcha/cv"
	"flag"
	_ "image"
	"log"
)

func Guess(link string) {
	img := captcha.LoadImage(link)
	/*
	decoder := &captcha.Decoder{
		ImageProcessors: []cv.ImageProcessor {
				&cv.MeanShift{K:1},
			},
		BiColorProcessor: &cv.PeakAverageBasedBiColor{},
		BinaryImageProcessors: []cv.BinaryImageProcessor{
			&cv.RemoveBinaryImageBorder{},
			&cv.RemoveIsolatePoints{},
			&cv.BoundBinaryImage{XMinOpen: 2, YMinOpen : 5},
			&cv.RemoveXAxis{K : 2},
			&cv.ScaleBinaryImage{Height : captcha.SCALE_HEIGHT},
		},
		BinaryImagePredictor: &captcha.BinaryImageConnectedComponentPredictor{Dx : 1, Dy : 1},
	}
	*/
	decoder := &captcha.Decoder{
		BiColorProcessor: &cv.PeakAverageBasedBiColor{},
		BinaryImageProcessors: []cv.BinaryImageProcessor{
			&cv.RemoveBinaryImageBorder{},
			&cv.RemoveIsolatePoints{},
			//&cv.Thining{},
			&cv.BoundBinaryImage{XMinOpen: 0, YMinOpen : 0},
			//&cv.RemoveXAxis{K : 1},
			&cv.ScaleBinaryImage{Height : captcha.SCALE_HEIGHT},
		},
		BinaryImagePredictor: &captcha.FastCutBasedPredictor{},
	}
	
	masks := captcha.LoadMasks("./masks")
	results := decoder.Predict(img, masks, captcha.NUMBER)
	
}

func main() {
	input := flag.String("input", "", "input")
	flag.Parse()

	Guess(*input)
}
