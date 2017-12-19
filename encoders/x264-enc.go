package encoders

import (
	"image"
	"io"
	"os"
	"os/exec"
	"strings"
	"screenrec/logger"
)

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
	err := encodePPM(enc.input, img)
	if err != nil {
		logger.Error("error while encoding image:", err)
	}
}
func (enc *X264ImageEncoder) Close() {

}
