package main

import (
	"image"
	"screenrec/encoders"

	"screenrec/logger"
	"screenrec/screenshot"
	"time"
)

// func saveImage(myImg image.Image) error {
// 	f, err := os.Create("img.jpg")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()

// 	return jpeg.Encode(f, myImg, nil)
// }

// func loadImage(path string) image.Image {
// 	logger.Debug("loading image: ", path)
// 	jpg, err := os.Open(path)
// 	if err != nil {
// 		logger.Error("error while loading image: ", err)
// 	}
// 	logger.Debug("image file opened: ", path)
// 	img, err := jpeg.Decode(jpg)
// 	if err != nil {
// 		logger.Error("error while decoding image: ", err)
// 	}
// 	logger.Debug("done loading image: ", path)
// 	return img
// }

func runCaptureLoop(videoFileName string, codec encoders.ImageEncoder) {

	codec.Init(videoFileName)
	go func() {
		for {
			img, err := screenshot.CaptureScreen()
			if err != nil {
				println("err %v", err)
			}
			myImg := image.Image(img)

			logger.Infof(">>sending screen in ppm, size: %v", myImg.Bounds().Max)
			codec.Encode(myImg)
		}

	}()
	codec.Run()
	time.Sleep(30 * time.Second)
	codec.Close()
}

func main() {

	enc := &encoders.MJPegImageEncoder{Quality: 60, Framerate: 6}
	//enc := &encoders.X264ImageEncoder{}
	runCaptureLoop("stam", enc)
}
