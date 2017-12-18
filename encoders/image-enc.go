package encoders

import (
	"fmt"
	"image"
	"image/color"
	"io"
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

type ImageEncoder interface {
	Init(string)
	Run()
	Encode(image.Image)
	Close()
}
