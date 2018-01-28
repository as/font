package gomono

import (
	"github.com/golang/freetype/truetype"
	. "golang.org/x/image/font/gofont/gomono"
)

var Font, _ = truetype.Parse(TTF)
