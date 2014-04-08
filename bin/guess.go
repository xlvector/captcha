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
	decoder := &captcha.Decoder{
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
	/*
	decoder := &captcha.Decoder{
		BiColorProcessor: &cv.PeakAverageBasedBiColor{},
		BinaryImageProcessors: []cv.BinaryImageProcessor{
			&cv.RemoveBinaryImageBorder{},
			&cv.RemoveIsolatePoints{},
			&cv.BoundBinaryImage{XMinOpen: 1, YMinOpen : 3},
			&cv.RemoveXAxis{K : 1},
			&cv.ScaleBinaryImage{Height : captcha.SCALE_HEIGHT},
		},
		BinaryImagePredictor: &captcha.FastCutBasedPredictor{},
	}
	*/
	masks := captcha.LoadMasks("./masks")
	log.Println(decoder.Predict(img, masks, captcha.MIX))
}

func main() {
	input := flag.String("input", "", "input")
	flag.Parse()

	Guess(*input)
}
