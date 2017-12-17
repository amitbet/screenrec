package main

import (
	"image"
	"image/jpeg"
	"os"
	"strings"
)

type MJpegEncoder struct {
	File *os.File
}

func NewMJpegEncoder(fileName string) *MJpegEncoder {
	if !strings.HasSuffix(fileName, ".mjpeg") {
		fileName = fileName + ".mjpeg"
	}
	file, err := os.Create(fileName)
	enc := &MJpegEncoder{File: file}
	if err != nil {
		panic(err)
	}

	return enc
}

func (m *MJpegEncoder) Close(fileName string) {
	m.File.Close()
}

func (m *MJpegEncoder) EncodeFrame(img image.Image) {
	jpeg.Encode(m.File, img, &jpeg.Options{Quality: 9})
}
