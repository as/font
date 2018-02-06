package font

import (
	"image"
	"unicode"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Replacer struct {
	Face
	b  Face
	fn func(r rune) bool
}

func NewReplacer(a, b Face, cond func(r rune) bool) Face {
	if cond == nil {
		cond = func(r rune) bool {
			return r > 127 || !unicode.IsGraphic(r) || r == 0
		}
	}
	return NewCache(&Replacer{
		Face: a,
		b:    b,
		fn:   cond,
	})
}

func (f *Replacer) Dx(p []byte) (dx int) {
	for n, c := range p {
		if f.fn(rune(c)) {
			dx += f.b.Dx(p[n : n+1])
		} else {
			dx += f.Face.Dx(p[n : n+1])
		}
	}
	return dx
}

func (f *Replacer) Fits(p []byte, limitDx int) (n int) {
	var c byte
	for n, c = range p {
		if f.fn(rune(c)) {
			limitDx -= f.b.Dx(p[n : n+1])
		} else {
			limitDx -= f.Face.Dx(p[n : n+1])
		}
		if limitDx < 0 {
			return n
		}
	}
	return n
}
func (f *Replacer) Close() error {
	return nil
}
func (f *Replacer) Glyph(dot fixed.Point26_6, r rune) (dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {
	if f.fn(r) {
		return f.b.Glyph(dot, r)
	}
	return f.Face.Glyph(dot, r)
}
func (f *Replacer) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	if f.fn(r) {
		return f.b.GlyphBounds(r)
	}
	return f.Face.GlyphBounds(r)
}
func (f *Replacer) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	if f.fn(r) {
		return f.b.GlyphAdvance(r)
	}
	return f.Face.GlyphAdvance(r)
}
func (f *Replacer) Kern(r0, r1 rune) fixed.Int26_6 {
	if f.fn(r0) {
		return f.b.Kern(r0,r1)
	}
	return f.Face.Kern(r0,r1)
}

func (f *Replacer) Metrics() (m font.Metrics) {
	return f.Face.Metrics()
}
