// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run gen.go

// Package vga provides fixed-size font faces.
package vga

import (
	"image"

	draw2 "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Range maps a contiguous range of runes to vertically adjacent sub-images of
// a Face's Mask image. The rune range is inclusive on the low end and
// exclusive on the high end.
//
// If Low <= r && r < High, then the rune r is mapped to the sub-image of
// Face.Mask whose bounds are image.Rect(0, y*h, Face.Width, (y+1)*h),
// where y = (int(r-Low) + Offset) and h = (Face.Ascent + Face.Descent).
type Range struct {
	Low, High rune
	Offset    int
}

var (
	kern     = draw2.Interpolator(draw2.NearestNeighbor)
	kern2    = draw2.Interpolator(draw2.CatmullRom)
	downkern = draw2.Interpolator(draw2.CatmullRom)
)

func Scale(img image.Image, r image.Rectangle) *image.Alpha {
	dst := image.NewAlpha(r)
	kern.Scale(dst, dst.Bounds(), img, img.Bounds(), draw2.Src, nil)
	return dst
}

// Face8x16 is a Face derived from the public domain X11 misc-fixed font files.
//
// At the moment, it holds the printable characters in ASCII starting with
// space, and the Unicode replacement character U+FFFD.
//
// Its data is entirely self-contained and does not require loading from
// separate files.
var Face8x16 = &Face{
	Advance: 8,
	Width:   8,
	Height:  16,
	Ascent:  16 - 4,
	Descent: 4,
	Mask:    mask8x16,
	Ranges: []Range{
		{'\u0020', '\u007f', 0},
		{'\ufffd', '\ufffe', 95},
	},
}

func NewFace(size int) *Face {
	r := mask8x16.Bounds()
	println(r.String())
	r.Max.Y = size * 96
	r.Max.X = size / 2
	println(r.String())

	mask := image.NewAlpha(r)
	kern.Scale(mask, mask.Bounds(), mask8x16, mask8x16.Bounds(), draw2.Src, nil)
	mask.Stride = size / 2

	f := &Face{
		Advance: size / 2,
		Width:   size / 2,
		Height:  size,
		Ascent:  size - (size / 4),
		Descent: (size / 4),
		Mask:    mask,
		Ranges: []Range{
			{'\u0020', '\u007f', 0},
			{'\ufffd', '\ufffe', 95},
		},
	}
	return f
}

// Face is a basic font face whose glyphs all have the same metrics.
//
// It is safe to use concurrently.
type Face struct {
	// Advance is the glyph advance, in pixels.
	Advance int
	// Width is the glyph width, in pixels.
	Width int
	// Height is the inter-line height, in pixels.
	Height int
	// Ascent is the glyph ascent, in pixels.
	Ascent int
	// Descent is the glyph descent, in pixels.
	Descent int
	// Left is the left side bearing, in pixels. A positive value means that
	// all of a glyph is to the right of the dot.
	Left int

	// Mask contains all of the glyph masks. Its width is typically the Face's
	// Width, and its height a multiple of the Face's Height.
	Mask image.Image
	// Ranges map runes to sub-images of Mask. The rune ranges must not
	// overlap, and must be in increasing rune order.
	Ranges []Range
}

func (f *Face) Close() error                   { return nil }
func (f *Face) Kern(r0, r1 rune) fixed.Int26_6 { return 0 }

func (f *Face) Metrics() font.Metrics {
	return font.Metrics{
		Height:  fixed.I(f.Height),
		Ascent:  fixed.I(f.Ascent),
		Descent: fixed.I(f.Descent),
	}
}

func (f *Face) Glyph(dot fixed.Point26_6, r rune) (
	dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {

loop:
	for _, rr := range [2]rune{r, '\ufffd'} {
		for _, rng := range f.Ranges {
			if rr < rng.Low || rng.High <= rr {
				continue
			}
			maskp.Y = (int(rr-rng.Low) + rng.Offset) * (f.Ascent + f.Descent)
			ok = true
			break loop
		}
	}
	if !ok {
		return image.Rectangle{}, nil, image.Point{}, 0, false
	}

	x := int(dot.X+32)>>6 + f.Left
	y := int(dot.Y+32) >> 6
	dr = image.Rectangle{
		Min: image.Point{
			X: x,
			Y: y - f.Ascent,
		},
		Max: image.Point{
			X: x + f.Width,
			Y: y + f.Descent,
		},
	}

	return dr, f.Mask, maskp, fixed.I(f.Advance), true
}

func (f *Face) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	return fixed.R(0, -f.Ascent, f.Width, +f.Descent), fixed.I(f.Advance), true
}

func (f *Face) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	return fixed.I(f.Advance), true
}
