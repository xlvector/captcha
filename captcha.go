package captcha

import (
	"fmt"
	_ "log"
	"net/http"
	"strconv"
	_ "strings"
	"time"
	"encoding/json"
)

type CaptchaHandler struct {
	Masks *MaskIndex
}

func NewCaptchaHandler() *CaptchaHandler {
	ret := CaptchaHandler{}
	ret.Masks = LoadMasks("./masks/")
	return &ret
}

func (self *CaptchaHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	link := req.FormValue("link")
	chType := GetChType(req.FormValue("type"))
	format := req.FormValue("format")
	if format == "" {
		format = "html"
	}
	img := LoadImageFromURL(link, "", "")

	output := "./imgs/" + strconv.FormatInt(time.Now().Unix(), 10)
	SaveImage(img, output)

	predictor := NewMetaPredictor()
	results := predictor.Predict(img, self.Masks, chType)

	if format == "html" {
		html := "<html><body style=\"background-color:#000; color:#fff;\"><img src=\"" + output + ".png\"><br/>"
		for _, result := range results {
			for _, component := range result.Components{
				component.Save("./masks/" + component.Encode())
				html += "<img src=\"" + "./masks/" + component.Encode() + ".png\"/>&nbsp;"
			}
			html += "<br/>"
			html += result.Label + "&nbsp;:&nbsp;" + strconv.FormatFloat(result.GetWeight(), 'g', 5, 64) + "<br/>"
		}
		html += "</body></html>"
		fmt.Fprint(w, html)
	} else if format == "json" {
		output, _ := json.Marshal(results)
		fmt.Fprint(w, string(output))
	}
}
