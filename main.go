package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"os"
	"os/exec"
	"screenrec/screenshot"
	"strings"
	"time"
	"vnc2webm/logger"

	"github.com/icza/mjpeg"
)

func encodePPMFast(w io.Writer, m image.Image) error {
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
func encodePPM(w io.Writer, m image.Image) error {
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

type DV8ImageEncoder struct {
	cmd   *exec.Cmd
	input io.WriteCloser
}

func (enc *DV8ImageEncoder) Init(videoFileName string) {
	fileExt := ".webm"
	if !strings.HasSuffix(videoFileName, fileExt) {
		videoFileName = videoFileName + fileExt
	}
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

		videoFileName,
	)
	//cmd := exec.Command("/bin/echo")

	//io.Copy(cmd.Stdout, os.Stdout)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	encInput, err := cmd.StdinPipe()
	enc.input = encInput
	if err != nil {
		logger.Error("can't get ffmpeg input pipe")
	}
	enc.cmd = cmd
}
func (enc *DV8ImageEncoder) Run() {
	logger.Debugf("launching binary: %v", enc.cmd.Args)
	err := enc.cmd.Run()
	if err != nil {
		logger.Error("error while launching ffmpeg:", err)
	}
}
func (enc *DV8ImageEncoder) Encode(img image.Image) {
	err := encodePPMFast(enc.input, img)
	if err != nil {
		logger.Error("error while encoding image:", err)
	}
}
func (enc *DV8ImageEncoder) Close() {

}

type X264ImageEncoder struct {
	cmd   *exec.Cmd
	input io.WriteCloser
}

func (enc *X264ImageEncoder) Init(videoFileName string) {
	fileExt := ".mp4"
	if !strings.HasSuffix(videoFileName, fileExt) {
		videoFileName = videoFileName + fileExt
	}
	binary := "./ffmpeg"
	cmd := exec.Command(binary,
		"-f", "image2pipe",
		"-vcodec", "ppm",
		//"-r", strconv.Itoa(framerate),
		"-r", "4",
		//"-i", "pipe:0",
		"-i", "-",
		"-vcodec", "libx264", //"libvpx",//"libvpx-vp9"//"libx264"
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

		videoFileName,
	)
	//cmd := exec.Command("/bin/echo")

	//io.Copy(cmd.Stdout, os.Stdout)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	encInput, err := cmd.StdinPipe()
	enc.input = encInput
	if err != nil {
		logger.Error("can't get ffmpeg input pipe")
	}
	enc.cmd = cmd
}
func (enc *X264ImageEncoder) Run() {
	logger.Debugf("launching binary: %v", enc.cmd.Args)
	err := enc.cmd.Run()
	if err != nil {
		logger.Error("error while launching ffmpeg:", err)
	}
}
func (enc *X264ImageEncoder) Encode(img image.Image) {
	err := encodePPMFast(enc.input, img)
	if err != nil {
		logger.Error("error while encoding image:", err)
	}
}
func (enc *X264ImageEncoder) Close() {

}

type DV9ImageEncoder struct {
	cmd   *exec.Cmd
	input io.WriteCloser
}

func (enc *DV9ImageEncoder) Init(videoFileName string) {
	fileExt := ".webm"
	if !strings.HasSuffix(videoFileName, fileExt) {
		videoFileName = videoFileName + fileExt
	}
	binary := "./ffmpeg"
	cmd := exec.Command(binary,
		"-f", "image2pipe",
		"-vcodec", "ppm",
		//"-r", strconv.Itoa(framerate),
		"-r", "3",
		//"-i", "pipe:0",
		"-i", "-",
		"-vcodec", "libvpx-vp9", //"libvpx",//"libvpx-vp9"//"libx264"
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

		videoFileName,
	)
	//cmd := exec.Command("/bin/echo")

	//io.Copy(cmd.Stdout, os.Stdout)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	encInput, err := cmd.StdinPipe()
	enc.input = encInput
	if err != nil {
		logger.Error("can't get ffmpeg input pipe")
	}
	enc.cmd = cmd
}
func (enc *DV9ImageEncoder) Run() {
	logger.Debugf("launching binary: %v", enc.cmd.Args)
	err := enc.cmd.Run()
	if err != nil {
		logger.Error("error while launching ffmpeg:", err)
	}
}
func (enc *DV9ImageEncoder) Encode(img image.Image) {
	err := encodePPMFast(enc.input, img)
	if err != nil {
		logger.Error("error while encoding image:", err)
	}
}
func (enc *DV9ImageEncoder) Close() {

}

type MJPegImageEncoder struct {
	avWriter  mjpeg.AviWriter
	Quality   int
	Framerate int32
}

func (enc *MJPegImageEncoder) Init(videoFileName string) {
	fileExt := ".avi"
	if !strings.HasSuffix(videoFileName, fileExt) {
		videoFileName = videoFileName + fileExt
	}
	if enc.Framerate <= 0 {
		enc.Framerate = 5
	}
	avWriter, err := mjpeg.New(videoFileName, 1024, 768, enc.Framerate)
	if err != nil {
		logger.Error("Error during mjpeg init: ", err)
	}
	enc.avWriter = avWriter
}
func (enc *MJPegImageEncoder) Run() {
}
func (enc *MJPegImageEncoder) Encode(img image.Image) {
	buf := &bytes.Buffer{}
	jOpts := &jpeg.Options{Quality: enc.Quality}
	if enc.Quality <= 0 {
		jOpts = nil
	}
	err := jpeg.Encode(buf, img, jOpts)
	if err != nil {
		logger.Error("Error while creating jpeg: ", err)
	}
	err = enc.avWriter.AddFrame(buf.Bytes())
	if err != nil {
		logger.Error("Error while adding frame to mjpeg: ", err)
	}

}
func (enc *MJPegImageEncoder) Close() {
	err := enc.avWriter.Close()
	if err != nil {
		logger.Error("Error while closing mjpeg: ", err)
	}
}

type ImageEncoder interface {
	Init(string)
	Run()
	Encode(image.Image)
	Close()
}

func runCaptureLoop(videoFileName string, codec ImageEncoder) {

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

	//enc := &MJPegImageEncoder{Quality: 60, Framerate: 6}
	enc := &X264ImageEncoder{}
	runCaptureLoop("stam", enc)
}
