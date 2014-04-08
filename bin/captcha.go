package main

import (
	"captcha"
	"flag"
	"fmt"
	_ "image"
	"log"
)

func TestAll(folder string) {
	imgs := captcha.LoadTestImages(folder)
	log.Println("test images", len(imgs))
	masks := captcha.LoadMasks("./masks")
	hit := 0.0
	total := 0.0
	for _, img := range imgs {
		predictor := captcha.NewMetaPredictor()
		log.Println(img.Path)
		ok := false
		for {
			copyImg := captcha.CopyImage(img.Img)
			results := predictor.Predict(copyImg, masks, img.ChType)
			if results == nil {
				break
			}
			for _, result := range results {
				log.Println(">>>>>>>>>", result.Label, img.Label)
				if result.Label == img.Label {
					log.Println("--------- ok")
					hit += 1.0
					ok = true
					break
				}
			}
			if ok {
				break
			}
		}
		if !ok {
			captcha.SaveImage(img.Img, "./failed/" + img.Label)
		}
		total += 1.0
		fmt.Println()
	}
	log.Println("hit ratio", hit/total)
}

func main() {
	input := flag.String("input", "", "input")
	flag.Parse()

	TestAll(*input)
}
