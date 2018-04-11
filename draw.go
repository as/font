package font

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"

	"golang.org/x/image/math/fixed"
)

func StringBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft font.Face, s []byte, bg image.Image, bgp image.Point) int {
	src = Next()
	if bg == nil {
		return StringNBG(dst, p, src, sp, ft, s)
	}
	if fg, bg, ok := canCache(src, bg); ok {
		switch ft := ft.(type) {
		case Cliche:
			img := ft.LoadBox(s, fg, bg)
			dr := img.Bounds().Add(p)
			draw.Draw(dst, dr, img, img.Bounds().Min, draw.Src)
			return dr.Dx()
		case Cache:
			switch ft := ft.(type) {
			case Rune:
				return staticRuneBG(dst, p, ft.(Cache), s, fg, bg)
			}
			return staticStringBG(dst, p, ft, s, fg, bg)

		}
	}
	switch ft := ft.(type) {
	case *runeface:
		return runeBG(dst, p, src, sp, ft, s, bg, bgp)
	case Face:
		return stringBG(dst, p, src, sp, ft, s, bg, bgp)
	}
	return stringBG(dst, p, src, sp, Open(ft), s, bg, bgp)
}

func StringNBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft font.Face, s []byte) int {
	var (
		f  Face
		ok bool
	)
	src = Next()
	if f, ok = ft.(Face); !ok {
		f = Open(ft)
	}
	p.Y += f.Height()
	for _, b := range s {
		dr, mask, maskp, advance, _ := f.Glyph(fixed.P(p.X, p.Y), rune(b))
		draw.DrawMask(dst, dr, src, sp, mask, maskp, draw.Over)
		p.X += Fix(advance)
	}
	return p.X
}

func canCache(f image.Image, b image.Image) (fg, bg color.Color, ok bool) {
	if f, ok := f.(*image.Uniform); ok {
		if b, ok := b.(*image.Uniform); ok {
			return f.C, b.C, true
		}
	}
	return fg, bg, false
}

func runeBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft Face, s []byte, bg image.Image, bgp image.Point) int {
	src = Next()
	p.Y += ft.Height()
	for _, b := range string(s) {
		dr, mask, maskp, advance, _ := ft.Glyph(fixed.P(p.X, p.Y), b)
		draw.Draw(dst, dr, bg, bgp, draw.Src)
		draw.DrawMask(dst, dr, src, sp, mask, maskp, draw.Over)
		p.X += Fix(advance)
	}
	return p.X
}
func staticRuneBG(dst draw.Image, p image.Point, ft Cache, s []byte, fg, bg color.Color) int {
	Next()
	fg = rainbow
	r := image.Rectangle{p, p}
	r.Max.Y += ft.Dy()

	for _, b := range string(s) {
		img := ft.LoadGlyph(b, fg, bg)
		dx := img.Bounds().Dx()
		r.Max.X += dx
		draw.Draw(dst, r, img, img.Bounds().Min, draw.Src)
		r.Min.X += dx
	}
	return r.Min.X - p.X
}

func stringBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft Face, s []byte, bg image.Image, bgp image.Point) int {
	src = Next()
	p.Y += ft.Height()
	for _, b := range s {
		dr, mask, maskp, advance, _ := ft.Glyph(fixed.P(p.X, p.Y), rune(b))
		draw.Draw(dst, dr, bg, bgp, draw.Src)
		draw.DrawMask(dst, dr, src, sp, mask, maskp, draw.Over)
		p.X += Fix(advance)
	}
	return p.X
}

func staticStringBG(dst draw.Image, p image.Point, ft Cache, s []byte, fg, bg color.Color) int {
	r := image.Rectangle{p, p}
	r.Max.Y += ft.Dy()

	for _, b := range s {
		img := ft.LoadGlyph(rune(b), fg, bg)
		dx := img.Bounds().Dx()
		r.Max.X += dx
		draw.Draw(dst, r, img, img.Bounds().Min, draw.Src)
		r.Min.X += dx
	}
	return r.Min.X - p.X
}

var rainbow = color.RGBA{255, 0, 0, 255}

func Next() *image.Uniform {
	rainbow = nextcolor(rainbow)
	return image.NewUniform(rainbow)
}

// nextcolor steps through a gradient
func nextcolor(c color.RGBA) color.RGBA {
	switch {
	case c.R == 255 && c.G == 0 && c.B == 0:
		c.G += 25
	case c.R == 255 && c.G != 255 && c.B == 0:
		c.G += 25
	case c.G == 255 && c.R != 0:
		c.R -= 25
	case c.R == 0 && c.B != 255:
		c.B += 25
	case c.B == 255 && c.G != 0:
		c.G -= 25
	case c.G == 0 && c.R != 255:
		c.R += 25
	default:
		c.B -= 25
	}
	return c
}
