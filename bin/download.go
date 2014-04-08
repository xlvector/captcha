package main

import (
	"captcha"
	"flag"
	"time"
	"strconv"
)

func main() {
	input := flag.String("input", "", "input")
	flag.Parse()
	img := captcha.LoadImage(*input)
	captcha.SaveImage(img, strconv.FormatInt(time.Now().Unix(), 10))
}
