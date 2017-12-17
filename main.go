package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"os"
	"os/exec"
	"screenrec/screenshot"
	"vnc2webm/logger"
)

func encodePPM1(w io.Writer, m image.Image) error {
	maxvalue := 255
	b := m.Bounds()
	// write header
	_, err := fmt.Fprintf(w, "P6\n%d %d\n%d\n", b.Dx(), b.Dy(), maxvalue)
	if err != nil {
		return err
	}

	// write raster
	cm := color.RGBAModel
	row := make([]uint8, b.Dx()*3)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		i := 0
		for x := b.Min.X; x < b.Max.X; x++ {
			c := cm.Convert(m.At(x, y)).(color.RGBA)
			row[i] = c.R
			row[i+1] = c.G
			row[i+2] = c.B
			i += 3
		}
		if _, err := w.Write(row); err != nil {
			return err
		}
	}
	return nil
}

// Encode an image.Image as a PPM image and write the result to w. The resulting
// image always has a MaxVal of 255. Alpha values are treated as if the image
// was displayed on a white background.
func EncodePPM(w io.Writer, m image.Image) error {
	logger.Debug("start encode ppm")
	bounds := m.Bounds().Canon()

	_, err := fmt.Fprintf(w, "P6\n%d %d\n255\n", bounds.Dx(), bounds.Dy())
	if err != nil {
		return err
	}

	triple := make([]byte, 3)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()

			triple[0] = byte(r >> 8)
			triple[1] = byte(g >> 8)
			triple[2] = byte(b >> 8)

			n, err := w.Write(triple)
			if n != 3 {
				return err
			}
			//logger.Debug("one ppm pixel: ", x)
		}
		//logger.Debug("one ppm line:", y)
	}
	logger.Debug("done writing ppm!")
	return nil
}
func saveImage(myImg image.Image) error {
	f, err := os.Create("img.jpg")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	return jpeg.Encode(f, myImg, nil)
}

func loadImage(path string) image.Image {
	logger.Debug("loading image: ", path)
	jpg, err := os.Open(path)
	if err != nil {
		logger.Error("error while loading image: ", err)
	}
	logger.Debug("image file opened: ", path)
	img, err := jpeg.Decode(jpg)
	if err != nil {
		logger.Error("error while decoding image: ", err)
	}
	logger.Debug("done loading image: ", path)
	return img
}

func main() {

	// if err != nil {
	// 	// replace this with real error handling
	// 	panic(err)
	// }

	const framerate int = 10
	//cmd := exec.Command("./ffmpeg", "-f", "image2pipe", "-vcodec", "ppm", "-r", "20", "-i", "-", "-r", "10", "-c:v", "libx264", "-preset", "slow", "-crf", "22", "-c:a", "copy", "video.avi")
	println("starting")
	//////////////
	// binary, lookErr := exec.LookPath("ffmpeg")
	// if lookErr != nil {
	// 	panic(lookErr)
	// }
	binary := "./ffmpeg"
	cmd := exec.Command(binary,
		"-f", "image2pipe",
		"-vcodec", "ppm",
		//"-r", strconv.Itoa(framerate),
		"-r", "3",
		//"-i", "pipe:0",
		"-i", "-",
		"-vcodec", "libvpx", //"libvpx",//"libvpx-vp9"//"libx264"
		"-b:v", "2M",
		"-threads", "8",
		//"-speed", "0",
		//"-lossless", "1", //for vpx
		// "-tile-columns", "6",
		//"-frame-parallel", "1",
		// "-an", "-f", "webm",
		"-cpu-used", "-16",

		"-preset", "ultrafast",
		"-deadline", "realtime",
		//"-cpu-used", "-5",
		"-maxrate", "2.5M",
		"-bufsize", "10M",
		"-g", "6",

		//"-rc_lookahead", "16",
		//"-profile", "0",
		"-qmax", "51",
		"-qmin", "11",
		//"-slices", "4",
		//"-vb", "2M",

		"./stam.webm",
	)
	//cmd := exec.Command("/bin/echo")

	encInput, err := cmd.StdinPipe()
	if err != nil {
		logger.Error("can't get ffmpeg input pipe")
	}
	//io.Copy(cmd.Stdout, os.Stdout)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	//os.Stdout.Write(out)
	//create screenshot frames and send to the ffmpeg process via stdin in ppm format
	logger.Infof("launching screencap routine")
	go func() {
		// infile, err := os.Open("./img.jpg")
		// if err != nil {
		// 	panic(err)
		// }
		// defer infile.Close()
		// myImg, _ := jpeg.Decode(infile)

		for {
			logger.Infof(">>>inside screencap routine")
			//defer encInput.Close()
			//----------------------

			img, err := screenshot.CaptureScreen()

			if err != nil {
				println("err %v", err)
			}
			myImg := image.Image(img)
			//-----------------------

			logger.Infof(">>sending screen in ppm, size: %v", myImg.Bounds().Max)
			//png.Encode(encInput, myImg)
			encodePPM1(encInput, myImg)

			// mjenc := NewMJpegEncoder("./stam")
			// mjenc.EncodeFrame(myImg)
			if err != nil {
				logger.Errorf("err while writing ppm to ffmpeg %v", err)
			}
			//time.Sleep(time.Duration(1000/framerate) * time.Millisecond)
		}

	}()

	logger.Debugf("launching binary: %v", cmd.Args)
	err = cmd.Run()
	//time.Sleep(1* time.Second)
	if err != nil {
		logger.Error(err)
	}
	//cmd.Wait()

	/////////////

	// encOutput, err := cmd.StdoutPipe()
	// if err != nil {
	// 	logger.Errorf("can't get ffmpeg output pipe %v", err)
	// }

	// //defer encInput.Close()
	// err = cmd.Run()
	// if err != nil {
	// 	logger.Error("error while running ffmpeg", err)
	// }

	// io.Copy(os.Stdout, encOutput)
	//cmd.Wait()
}
