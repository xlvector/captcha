package main

import (
	"net/http"
	"captcha"
	"time"
	"log"
)

func main() {
	http.Handle("/captcha", captcha.NewCaptchaHandler())
	http.Handle("/", http.FileServer(http.Dir("./")))
	s := &http.Server{
		Addr:           ":8900",
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}