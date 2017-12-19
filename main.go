package main

import (
	"image"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"screenrec/encoders"
	"screenrec/logger"
	"screenrec/screenshot"
	"syscall"
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
		iter := 0
		startTime := time.Now()
		for {
			img, err := screenshot.CaptureScreen()
			if err != nil {
				println("err %v", err)
			}
			myImg := image.Image(img)
			timeSinceStart := time.Since(startTime).Seconds()
			logger.Infof(">>sending screen in ppm, fps: %f size: %v", float64(iter)/timeSinceStart, myImg.Bounds().Max)
			iter++

			codec.Encode(myImg)
		}

	}()
	codec.Run()
	time.Sleep(30 * time.Second)
	codec.Close()

}

func main() {
	argsWithoutProg := os.Args[1:]
	runWithProfiler := len(argsWithoutProg) > 0
	if runWithProfiler {
		profFile := argsWithoutProg[0]
		f, err := os.Create(profFile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		s := <-sigc
		println("got signal:", s.String())
		if s != nil && runWithProfiler {
			pprof.StopCPUProfile()
		}
		os.Exit(1)
		// ... do something ...
	}()

	//enc := &encoders.MJPegImageEncoder{Quality: 60, Framerate: 6}
	enc := &encoders.X264ImageEncoder{}
	runCaptureLoop("stam", enc)
}
