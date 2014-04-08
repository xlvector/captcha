package captcha

import (
	"image"
	_ "image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "code.google.com/p/go.image/bmp"
	_ "code.google.com/p/go.image/tiff"
	"image/png"
	_ "log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func dialTimeout(network, addr string) (net.Conn, error) {
	timeout := time.Duration(10) * time.Second
	deadline := time.Now().Add(timeout)
	c, err := net.DialTimeout(network, addr, timeout)
	if err != nil {
		return nil, err
	}
	c.SetDeadline(deadline)
	return c, nil
}

func CopyImage(img image.Image) image.Image {
	ret := image.NewRGBA(img.Bounds())
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			ret.Set(x, y, img.At(x, y))
		}
	}
	return ret
}

func LoadImageFromURL(link string, proxy string, pageLink string) image.Image {
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
				ResponseHeaderTimeout: 10 * time.Second,
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
	}
	req.Header.Add("Referer", pageLink)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")
	resp, err := client.Do(req)

	if err != nil || resp == nil || resp.Body == nil {
		return nil
	} else {
		defer resp.Body.Close()
		img, _, err := image.Decode(resp.Body)
		if err != nil {
			panic(err)
		}
		return img
	}
}

func LoadImageFromFile(fname string) image.Image {
	reader, err := os.Open(fname)
	if err != nil {
		return nil
	}
	defer reader.Close()
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil
	}
	return img
}

func LoadImage(fname string) image.Image {
	if strings.HasPrefix(fname, "http://") || strings.HasPrefix(fname, "https://") {
		return LoadImageFromURL(fname, "", "")
	} else {
		return LoadImageFromFile(fname)
	}
}

func SaveImage(img image.Image, name string) {
	f, _ := os.Create(name + ".png")
	defer f.Close()
	png.Encode(f, img)
}

type TestImage struct {
	Label  string
	ChType int
	Img    image.Image
	Path   string
}

func LoadTestImages(folder string) []*TestImage {
	ret := []*TestImage{}
	filepath.Walk(folder, func(root string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		testImg := TestImage{}
		tks := strings.Split(f.Name(), ".")
		tks = strings.Split(tks[0], "_")
		testImg.Label = tks[0]
		testImg.ChType = MIX
		testImg.Path = root
		if len(tks) > 1 {
			if tks[1] == "num" {
				testImg.ChType = NUMBER
			} else if tks[1] == "alpha" {
				testImg.ChType = ALPHA
			}
		}
		testImg.Img = LoadImage(root)
		if testImg.Img != nil {
			ret = append(ret, &testImg)
		}
		return nil
	})
	return ret
}
