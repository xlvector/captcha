package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func Color2Double(color color.Color) float64 {
	r, g, b, _ := color.RGBA()
	c := (r + g + b) / (255 * 3)
	if c > 100 {
		return 1.0
	} else {
		return -1.0
	}
}

func EncodeXY(x int, y int) int {
	return x*10000 + y
}

func MinusXY(xy int, px int, py int) int {
	x := xy / 10000
	y := xy % 10000
	return EncodeXY(x-px, y-py)
}

func AddXY(xy int, px int, py int) int {
	x := xy / 10000
	y := xy % 10000
	return EncodeXY(x+px, y+py)
}

func Contains(m map[int]bool, k int) bool {
	v, ok := m[k]
	return v && ok
}

func BiColor(img image.Image) map[int]bool {
	out := make(map[int]bool)
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			r, g, b, a := img.At(x, y).RGBA()
			c := (r + g + b) / (255 * 3)
			if c < 200 && a > 20 {
				out[EncodeXY(x, y)] = true
			}
		}
	}
	return out
}

func GetImageXWeight(src map[int]bool) map[int]int {
	ret := make(map[int]int)
	for xy, _ := range src {
		x := xy / 10000
		_, ok := ret[x]
		if !ok {
			ret[x] = 0
		}
		ret[x] += 1
	}
	return ret
}

func MatchMask(src map[int]bool, srcLeft int, srcTop int, srcXWeight map[int]int, mask map[int]bool) float64 {
	mw, _ := getSize(mask)
	sim := 0.0
	r1 := 0.0
	for i := srcLeft; i <= srcLeft+mw; i++ {
		v, ok := srcXWeight[i]
		if ok {
			r1 += float64(v)
		}
	}
	r2 := (float64)(len(mask))
	for xy, _ := range mask {
		newXY := AddXY(xy, srcLeft, srcTop)
		if Contains(src, newXY) {
			sim += 1.0
		}
	}
	s1 := sim / r1
	s2 := sim / r2
	return s1*0.3 + s2*0.7
}

func dialTimeout(network, addr string) (net.Conn, error) {
	timeout := time.Duration(1) * time.Second
	deadline := time.Now().Add(timeout)
	c, err := net.DialTimeout(network, addr, timeout)
	if err != nil {
		return nil, err
	}
	c.SetDeadline(deadline)
	return c, nil
}

func LoadImageFromURL(link string, proxy string) image.Image {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil || req == nil || req.Header == nil {
		return nil
	}
	client := &http.Client{}
	if proxy != "" {
		proxyUrl, _ := url.Parse(proxy)
		client = &http.Client{
			Transport: &http.Transport{
				Dial:                  dialTimeout,
				DisableKeepAlives:     true,
				ResponseHeaderTimeout: 1 * time.Second,
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
	}
	resp, err := client.Do(req)

	if err != nil || resp == nil || resp.Body == nil {
		return nil
	} else {
		defer resp.Body.Close()
		format := strings.Trim(resp.Header.Get("Content-Type"), "; ,")
		log.Println(format)
		var img image.Image
		if format == "image/png" {
			img, err = png.Decode(resp.Body)
		} else if format == "image/jpeg" {
			img, err = jpeg.Decode(resp.Body)
		} else if format == "image/gif" {
			img, err = gif.Decode(resp.Body)
		} else {
			log.Println(link)
		}
		if err != nil {
			log.Println(err)
			return nil
		}
		return img
	}
}

func LoadImage(fname string) image.Image {
	tks := strings.Split(fname, ".")
	format := tks[len(tks)-1]
	reader, err := os.Open(fname)
	if err != nil {
		log.Println(err)
	}
	defer reader.Close()
	var img image.Image
	if format == "png" {
		img, err = png.Decode(reader)
	} else if format == "jpeg" {
		img, err = jpeg.Decode(reader)
	} else if format == "gif" {
		img, err = gif.Decode(reader)
	}
	if err != nil {
		log.Println(err)
	}
	return img
}

type Mask struct {
	Img   map[int]bool
	Label string
}

func LoadMasks(folder string) map[string][]*Mask {
	masks := make(map[string][]*Mask)
	err := filepath.Walk(folder, func(root string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		tks := strings.Split(f.Name(), "_")
		if len(tks) == 2 {
			img := LoadImage(root)
			if img != nil {
				mask := Mask{}
				mask.Img = BiColor(img)
				mask.Label = tks[0]
				_, ok := masks[mask.Label]
				if !ok {
					masks[mask.Label] = []*Mask{}
				}
				masks[mask.Label] = append(masks[mask.Label], &mask)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return masks
}

func FindBeginPos(img map[int]bool, beginX int) int {
	width, height := getSize(img)
	for x := beginX; x < width; x++ {
		for y := 0; y < height/3; y++ {
			if Contains(img, EncodeXY(x, y)) {
				return x
			}
		}
		for y := 2 * height / 3; y <= height; y++ {
			if Contains(img, EncodeXY(x, y)) {
				return x
			}
		}
	}
	return width
}

func getSize(img map[int]bool) (int, int) {
	width := 0
	height := 0
	for xy, _ := range img {
		x := xy / 10000
		y := xy % 10000
		if width < x {
			width = x
		}
		if height < y {
			height = y
		}
	}
	return width, height
}

func Recognize(srcImg map[int]bool, masks map[string][]*Mask) string {
	width, height := getSize(srcImg)
	srcXWeight := GetImageXWeight(srcImg)
	curX := FindBeginPos(srcImg, 0)
	log.Println("begin pos", curX)
	ret := ""
	lastMaskWidth := 0
	for {
		if curX >= width {
			break
		}
		maxSim := 0.0
		nextX := curX
		//nextY := 0
		var bestMask *Mask
		label := ""
		beginX := FindBeginPos(srcImg, curX+2)
		endX := beginX + 2
		if beginX == curX+2 {
			beginX = curX - lastMaskWidth/3
			endX = curX + 2
		}
		if beginX < 0 {
			beginX = 0
		}
		for _, labelMasks := range masks {
			localMaxSim := 0.0
			var localBestMask *Mask
			localNextX := curX
			localLabel := ""
			for k, mask := range labelMasks {
				for x := beginX; x <= endX; x++ {
					for y := 0; y <= height/3; y++ {
						sim := MatchMask(srcImg, x, y, srcXWeight, mask.Img)
						if sim > localMaxSim {
							localMaxSim = sim
							localBestMask = mask
							localNextX = x
							localLabel = mask.Label
						}
					}
				}
				if k > 30 && localMaxSim < 0.7 {
					break
				}
			}
			if localMaxSim > maxSim {
				maxSim = localMaxSim
				nextX = localNextX
				label = localLabel
				bestMask = localBestMask
			}
		}
		if bestMask == nil {
			break
		}

		maskWidth, _ := getSize(bestMask.Img)
		curX = nextX + maskWidth
		lastMaskWidth = maskWidth
		log.Println(maxSim, label)
		if maxSim < 0.8 {
			break
		}

		ret += label
	}
	return ret
}

func ExtractDomain(path string) string {
	tks := strings.Split(path, "/")
	if len(tks) < 3 {
		return ""
	}
	ret := tks[0] + "//" + tks[2]
	return ret
}

func ConcatLink(root0 string, link0 string) string {
	if len(link0) == 0 {
		return root0
	}
	root := root0
	if link0[0] == '/' {
		root = ExtractDomain(root0)
	}
	link := strings.ToLower(link0)
	if strings.Index(link, "http://") == 0 || strings.Index(link, "https://") == 0 {
		return link0
	} else {
		srcTks := strings.Split(root, "/")
		dstTks := strings.Split(link0, "/")
		n := 0
		k := -1
		for i, tk := range dstTks {
			if tk == ".." {
				n += 1
			} else if tk == "." {
				n += 0
			} else {
				k = i
				break
			}
		}
		ret := ""
		if len(srcTks) < n {
			return ""
		}
		for _, tk := range srcTks[:len(srcTks)-n] {
			ret += tk
			ret += "/"
		}
		if k >= 0 {
			for _, tk := range dstTks[k:] {
				if len(tk) == 0 {
					continue
				}
				ret += tk
				ret += "/"
			}
		}
		if link0[len(link0)-1] != '/' {
			return strings.TrimRight(ret, "/")
		} else {
			return ret
		}
	}
}

func MD5(buf string) string {
	m := md5.New()
	m.Write([]byte(buf))
	return hex.EncodeToString(m.Sum(nil))
}

func SaveImage(img image.Image, name string) {
	f, _ := os.Create(name + ".png")
	defer f.Close()
	png.Encode(f, img)
}

func LoadProcessedImages() map[string]string {
	f, _ := os.Open("label.tsv")
	defer f.Close()
	r := bufio.NewReader(f)
	ret := make(map[string]string)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		tks := strings.Split(line, "\t")
		ret[tks[0]] = strings.Trim(tks[1], "\n")
	}
	return ret
}

func LoadProxyList() []string {
	ret := []string{}
	f, _ := os.Open("proxy.list")
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		ret = append(ret, strings.Trim(line, "\n"))
	}
	return ret
}

func ScanList(fname string, output string, masks map[string][]*Mask) {
	proxys := LoadProxyList()
	proced := LoadProcessedImages()

	fin, _ := os.Open(fname)
	defer fin.Close()
	fout, _ := os.Create(output)
	defer fout.Close()

	for link, phone := range proced {
		fout.WriteString(link + "\t" + phone + "\n")
	}

	pattern := regexp.MustCompile("[0-9]+")
	r := bufio.NewReader(fin)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.Trim(line, "\n")
		tks := strings.Split(line, "\t")
		if pattern.FindString(strings.Replace(tks[1], " ", "", -1)) == strings.Replace(tks[1], " ", "", -1) {
			continue
		}
		link := ConcatLink(tks[0], tks[1])
		phone, ok := proced[link]
		if ok {
			fmt.Println("skip", link, phone)
			continue
		}
		if !strings.Contains(link, ".ganji.com") {
			continue
		}
		img := LoadImageFromURL(link, proxys[rand.Intn(len(proxys))])
		if img == nil {
			continue
		}
		eimg := BiColor(img)
		phone = Recognize(eimg, masks)
		if len(phone) == 11 && phone[0] == '1' {
			fmt.Println(link, phone)
			fout.WriteString(link + "\t" + phone + "\n")
			//SaveImage(img, phone)
		}
	}
}

func main() {
	imgFile := flag.String("img", "", "image file name")
	maskFolder := flag.String("mask", "masks", "mask folder name")
	listFile := flag.String("list", "", "list file name")
	flag.Parse()
	masks := LoadMasks(*maskFolder)
	log.Println("load", len(masks), "masks")

	if *listFile != "" {
		ScanList(*listFile, "label.tsv", masks)
	} else {
		var img image.Image
		if strings.Contains(*imgFile, "http://") {
			img = LoadImageFromURL(*imgFile, "")
		} else {
			img = LoadImage(*imgFile)
		}
		if img == nil {
			log.Println("load image failed")
			os.Exit(1)
		}
		SaveImage(img, "origin")
		eimg := BiColor(img)
		log.Println(Recognize(eimg, masks))
	}
}
