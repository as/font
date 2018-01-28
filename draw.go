package font

import (
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"unicode/utf8"
)

func StringBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft font.Face, s []byte, bg image.Image, bgp image.Point) int{
	cfg, bfg, ok := canCache(src, bg)
	if !ok{
		return stringBG(dst, p, src, sp, ft, s, bg, bgp)
	}
	switch ft := ft.(type){
	case Cliche:
			img := ft.LoadBox(s, fg, bg)
			dr := img.Bounds().Add(p)
			draw.Draw(dst, dr, img, img.Bounds().Min, draw.Src)
			return dr.Dx()
	case Cache:
		
	}
}

func StringNBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft Face, s []byte, bg image.Image, bgp image.Point) int{
	p.Y += ft.Height() 
	for _, b := range s {
		dr, mask, maskp, advance, _ := ft.Glyph(fixed.P(p.X, p.Y), rune(b))
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

func stringBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft Font, s []byte, bg image.Image, bgp image.Point) int{
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
	r.Max.Y += Dy(ft)
	for _, b := range s {
		img := ft.LoadGlyph(b, fg, bg)
		dx := img.Bounds().Dx()
		r.Max.X += dx
		draw.Draw(dst, r, img, img.Bounds().Min, draw.Src)
		r.Min.X += dx 
	}
	return r.Min.X-p.X
}

/*
func StringBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft *Font, s []byte, bg image.Image, bgp image.Point) int {
	for _, b := range s {
		mask := ft.Char(b)
		if mask == nil {
			panic("StringBG")
		}
		r := mask.Bounds()
		//draw.Draw(dst, r.Add(p), bg, bgp, draw.Src)
		draw.DrawMask(dst, r.Add(p), src, sp, mask, mask.Bounds().Min, draw.Over)
		p.X += r.Dx() + ft.stride
	}
	return p.X
}

func StringNBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft *Font, s []byte) int {
	for _, b := range s {
		mask := ft.Char(b)
		if mask == nil {
			panic("StringBG")
		}
		r := mask.Bounds()
		draw.DrawMask(dst, r.Add(p), src, sp, mask, mask.Bounds().Min, draw.Over)
		p.X += r.Dx() + ft.stride
	}
	return p.X
}

func RuneBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft *Font, s []byte, bg image.Image, bgp image.Point) int {
	p.Y += ft.Size()
	for {
		b, size := utf8.DecodeRune(s)
		dr, mask, maskp, advance, ok := ft.Glyph(fixed.P(p.X, p.Y), b)
		if !ok {
			panic("RuneBG")
		}
		//draw.Draw(dst, dr, bg, bgp, draw.Src)
		draw.DrawMask(dst, dr, src, sp, mask, maskp, draw.Over)
		p.X += Fix(advance)
		if len(s)-size == 0 {
			break
		}
		s = s[size:]
	}
	return p.X
}

func RuneNBG(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft *Font, s []byte) int {
	p.Y += ft.Size()
	for {
		b, size := utf8.DecodeRune(s)
		dr, mask, maskp, advance, ok := ft.Glyph(fixed.P(p.X, p.Y), b)
		if !ok {
			panic("RuneBG")
		}
		draw.DrawMask(dst, dr, src, sp, mask, maskp, draw.Over)
		p.X += Fix(advance)
		if len(s)-size == 0 {
			break
		}
		s = s[size:]
	}
	return p.X
}
*/