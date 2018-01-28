package goregular

import (
	"github.com/golang/freetype/truetype"
	. "golang.org/x/image/font/gofont/goregular"
)

var Font, _ = truetype.Parse(TTF)
