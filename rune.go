package font

import (
	"image"

	"golang.org/x/image/font"
)

// Rune is a Face aware of UTF8 text encoding. It measures runes
// rather than bytes and draws runs of UTF8 text when used with
// StringBG
type Rune interface {
	Face
	canrune()
}

type runeface struct {
	Face
}

func NewRune(f font.Face) Face {
	if f == nil {
		panic("open: nil face")
	}
	return &runeface{Open(f)}
}

func (*runeface) canrune() {}

func (f *runeface) Dx(p []byte) (dx int) {
	for _, r := range string(p) {
		w, _ := f.Face.GlyphAdvance(r)
		dx += Fix(w)
	}
	return dx + f.Stride()*len(p)
}

func (f *runeface) Fits(p []byte, limitDx int) (n int) {
	var r rune
	stride := f.Stride()
	for n, r = range string(p) {
		w, _ := f.Face.GlyphAdvance(r)
		limitDx -= Fix(w) + stride
		if limitDx < 0 {
			return n
		}
	}
	return n
}

func newRuneCache(r Rune) Cache {
	cf := &cachedRuneFace{
		&cachedFace{
			Face:  r,
			cache: make(map[signature]*image.RGBA),
		},
	}
	for i := range cf.cachewidth {
		cf.cachewidth[i] = r.Dx([]byte{byte(i)})
	}
	return cf
}

type cachedRuneFace struct {
	*cachedFace
}

func (f *cachedRuneFace) canrune() {}

func (f *cachedRuneFace) Dx(p []byte) (dx int) {
	for _, c := range string(p) {
		if int(c) < len(f.cachewidth) && c > -1 {
			dx += f.cachewidth[c]
		} else {
			w, _ := f.Face.GlyphAdvance(c)
			dx += Fix(w)
		}
	}
	return dx
}

func (f *cachedRuneFace) Fits(p []byte, limitDx int) (n int) {
	var r rune
	stride := f.Stride()
	for n, r = range string(p) {
		if int(r) < len(f.cachewidth) && r > -1 {
			limitDx -= f.cachewidth[r]
		} else {
			w, _ := f.Face.GlyphAdvance(r)
			limitDx -= Fix(w) + stride
		}
		if limitDx < 0 {
			return n
		}
	}
	return n
}
