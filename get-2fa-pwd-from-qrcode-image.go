/*************************************************************************
	> File Name: get-2fa-pwd-from-qrcode-image.go
	> Author: xiaoyang.chen
	> Mail: xiaoyang.chen@???.com
	> Created Time: Fri 12/27 09:08:35 2024
 ************************************************************************/
/*
 * get-2fa-pwd-from-qrcode-image.go
 * go build -ldflags="-s -w" -o ./get-2fa-pwd-from-qrcode-image ./get-2fa-pwd-from-qrcode-image.go
 * example:
 * ./get-2fa-pwd-from-qrcode-image ../../Downloads/00019773.image
 * go run -C ~/scripts ~/scripts/get-2fa-pwd-from-qrcode-image.go ~/Downloads/00019773.image, -C set build and work dir, but param must be absolute path, or use relative dir for the workdir that -C set
 * go run -C ~/scripts ~/scripts/get-2fa-pwd-from-qrcode-image.go ../../Downloads/00019773.image
 * go run ~/scripts/get-2fa-pwd-from-qrcode-image.go ../../Downloads/00019773.image
 * go run ~/scripts/get-2fa-pwd-from-qrcode-image.go ~/Downloads/00019773.image
 * go run ~/scripts/get-2fa-pwd-from-qrcode-image.go "~/Downloads/00019773.image"
 * params:
 * 1. filepath of image having qrcode
 */

// Package main is the entry of script.
package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/url"
	"os"
	"time"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/pquerna/otp/totp"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
)

func main() {

	fmt.Println("vim-go")
	// get params
	var lenOsArgs = len(os.Args)
	fmt.Printf("command-args-cnt: %d, command-args: %v\n", lenOsArgs, os.Args)
	if lenOsArgs < 2 {
		fmt.Println("lenOsArgs < 2")
		return
	}
	// open and decode image file
	var file, err = os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("os.Open err: %s\n", err.Error())
		return
	}
	defer file.Close()
	decodeImg, format, err := image.Decode(file)
	if err != nil {
		fmt.Printf("image.Decode err: %s\n", err.Error())
		return
	}
	fmt.Println(format, err) // jpeg <nil>
	// prepare BinaryBitmap
	goZxingBmp, err := gozxing.NewBinaryBitmapFromImage(decodeImg)
	if err != nil {
		fmt.Printf("gozxing.NewBinaryBitmapFromImage err: %s\n", err.Error())
		return
	}
	// decode image
	result, err := qrcode.NewQRCodeReader().Decode(goZxingBmp, nil)
	if err != nil {
		fmt.Printf("qrcode.NewQRCodeReader().Decode err: %s\n", err.Error())
		return
	}
	// parse result and generate 2fa password
	parseUrl, err := url.Parse(result.GetText())
	if err != nil {
		fmt.Printf("url.Parse err: %s\n", err.Error())
		return
	}
	var strSecret = parseUrl.Query().Get("secret")
	if strSecret == "" {
		fmt.Println("strSecret == \"\" from qrcode in image")
		return
	}
	//fmt.Println(strSecret)
	//passcode, err := totp.GenerateCodeCustom(strSecret, time.Now().UTC(), totp.ValidateOpts{
	//	Period:    30,
	//	Skew:      1,
	//	Digits:    otp.DigitsSix,
	//	Algorithm: otp.AlgorithmSHA1,
	//})
	passcode, err := totp.GenerateCode(strSecret, time.Now().UTC())
	if err != nil {
		fmt.Printf("totp.GenerateCode err: %s\n", err.Error())
		return
	}
	fmt.Println("passcode: " + passcode)
	fmt.Println("get-2fa-pwd-from-qrcode-image.go finish")
}
