package captcha

import (
	"fmt"
	_ "log"
	"net/http"
	"math/rand"
	"strconv"
	_ "strings"
	"time"
	"encoding/json"
)

type CaptchaHandler struct {
	Masks *MaskIndex
	Links []string
}

func NewCaptchaHandler() *CaptchaHandler {
	ret := CaptchaHandler{}
	ret.Masks = LoadMasks("./masks/")
	ret.Links = []string{
		"http://qyxy.baic.gov.cn/checkCodeServlet?currentTimeMillis=1396966478094",
		"http://www.tjaic.gov.cn/info/inc/checkcode.asp",
		"http://www.hebscztxyxx.gov.cn/notice/img.captcha?ra=0.420229553245008",
		"http://218.26.1.108/Manage/VerifyCode.aspx?",
		"http://www.nmgs.gov.cn:7001/aiccips/verify.html",
		"http://222.171.175.16:9080/ECPS/validateCode.jspx",
		"http://www.jsgsj.gov.cn:58888/ecipplatform/rand_img.jsp?type=2",
		"http://gsxt.zjaic.gov.cn/common/captcha/doReadKaptcha.do",
		"http://www.ahcredit.gov.cn:9090/ECPS/validateCode.jspx",
		"http://wsgs.fjaic.gov.cn/creditpub/govRandValidate",
		"http://gsxt.jxaic.gov.cn/ECPS/rand.action",
		"http://218.57.139.24/shandong/secimg",
		"http://211.141.74.198:8081/aiccips/verify.html",
		"http://gsxt.lngs.gov.cn/saicpub/jcaptcha",
		"http://gsxt.gdgs.gov.cn/verify.html",
		"http://221.232.141.250:9080/ECPS/rand.action",
		"http://gsxt.hnaic.gov.cn/VerifierImage.jpg?noCache=1396969773230",
		"http://gsxt.cqgs.gov.cn/makeCertPic.jsp",
		"http://gsxt.gzgs.gov.cn/validCode?validTag=searchImageCode&1396969875371",
		"http://220.163.86.132:88/admin/_SYSTEM/CheckCode.aspx",
		"http://117.22.252.219:8002/servlet/validateCodeServlet",
		"http://gsxt.ngsh.gov.cn/ECPS/rand.action",
		"http://www.360doc.com/VerifyCode.aspx?",
	}
	return &ret
}

func (self *CaptchaHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	link := req.FormValue("link")
	if len(link) == 0 {
		link = self.Links[rand.Intn(len(self.Links))]
	}
	chType := GetChType(req.FormValue("type"))
	format := req.FormValue("format")
	if format == "" {
		format = "html"
	}
	img := LoadImageFromURL(link, "", "")

	output := "./imgs/" + strconv.FormatInt(time.Now().UnixNano(), 10)
	SaveImage(img, output)

	predictor := NewMetaPredictor()
	results := predictor.Predict(img, self.Masks, chType)

	if format == "html" {
		html := "<html><body style=\"background-color:#000; color:#fff;\"><img src=\"" + output + ".png\"><br/>"
		for _, result := range results {
			for _, component := range result.Components{
				component.Save("./masks/" + component.Encode())
				html += "<img style=\"border:solid 1px #fff;\" src=\"" + "./masks/" + component.Encode() + ".png\"/>&nbsp;"
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
