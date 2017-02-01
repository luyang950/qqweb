package controllers

import (
	"fmt"
	"qqweb/libraries/log"
	"net/http"
	"io/ioutil"
	"github.com/davecgh/go-spew/spew"
)

const(
	// get login QR code
	qrCodeUrl = "https://ssl.ptlogin2.qq.com/ptqrshow?appid=501004106&e=0&l=M&s=5&d=72&v=4&t=0.1"

	// check if QR code is available
)

var logger *log.Log

func init() {
	logger = log.DLog
}

func Login() {
	fmt.Println("----->trying to login")
	fmt.Println("----->downloading QR code")
	getQRCode()

}

func getQRCode() {
	req, reqErr := http.NewRequest("GET", qrCodeUrl, nil)
	logger.CheckErr(reqErr, log.ERROR, false)

	res, doErr := http.DefaultClient.Do(req)
	logger.CheckErr(doErr, log.ERROR, false)

	defer res.Body.Close()
	body, readErr := ioutil.ReadAll(res.Body)
	logger.CheckErr(readErr, log.ERROR, false)

	spew.Dump(body)
}
