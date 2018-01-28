package font

import (
	"github.com/as/font/gomono"
	"github.com/as/font/goregular"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func NewGoMono(size int) font.Face {
	return truetype.NewFace(gomono.Font, &truetype.Options{
		SubPixelsX: 64,
		SubPixelsY: 64,
		Size:       float64(size),
	})
}

func NewGoRegular(size int) font.Face {
	return truetype.NewFace(goregular.Font, &truetype.Options{
		SubPixelsX: 64,
		SubPixelsY: 64,
		Size:       float64(size),
	})
}
