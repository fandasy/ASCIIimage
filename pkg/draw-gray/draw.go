package draw_gray

import (
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"unicode/utf8"
)

type Drawer struct {
	Dst  draw.Image
	Src  image.Image
	Face font.Face
	Dot  fixed.Point26_6
}

func Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	drawMask(dst, r, src, sp, nil, image.Point{})
}

func (d *Drawer) DrawBytes(s []byte) {
	prevC := rune(-1)
	for len(s) > 0 {
		c, size := utf8.DecodeRune(s)
		s = s[size:]
		if prevC >= 0 {
			d.Dot.X += d.Face.Kern(prevC, c)
		}
		dr, mask, maskp, advance, _ := d.Face.Glyph(d.Dot, c)
		if !dr.Empty() {
			drawMask(d.Dst, dr, d.Src, image.Point{}, mask, maskp)
		}
		d.Dot.X += advance
		prevC = c
	}
}

func clip(dst draw.Image, r *image.Rectangle, src image.Image, sp *image.Point, mask image.Image, mp *image.Point) {
	orig := r.Min
	*r = r.Intersect(dst.Bounds())
	*r = r.Intersect(src.Bounds().Add(orig.Sub(*sp)))
	if mask != nil {
		*r = r.Intersect(mask.Bounds().Add(orig.Sub(*mp)))
	}
	dx := r.Min.X - orig.X
	dy := r.Min.Y - orig.Y
	if dx == 0 && dy == 0 {
		return
	}
	sp.X += dx
	sp.Y += dy
	if mp != nil {
		mp.X += dx
		mp.Y += dy
	}
}

func drawMask(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point, mask image.Image, mp image.Point) {
	clip(dst, &r, src, &sp, mask, &mp)
	if r.Empty() {
		return
	}

	switch dst0 := dst.(type) {
	case *image.Gray:
		if mask == nil {
			switch src0 := src.(type) {
			case *image.Uniform:
				var Y uint8
				switch s := src0.C.(type) {
				case color.Gray:
					Y = s.Y
				case *color.Gray:
					Y = s.Y
				}
				if Y == 0xff {
					drawFillSrcGray(dst0, r, Y)
				} else {
					drawFillOverGray(dst0, r, Y)
				}
				return
			}
		} else if mask0, ok := mask.(*image.Alpha); ok {
			switch src0 := src.(type) {
			case *image.Uniform:
				var Y uint8
				switch s := src0.C.(type) {
				case color.Gray:
					Y = s.Y
				case *color.Gray:
					Y = s.Y
				}
				drawGlyphOverGray(dst0, r, Y, mask0, mp)
				return
			}
		}
	case *image.Gray16:
		if mask == nil {
			switch src0 := src.(type) {
			case *image.Uniform:
				var Y uint16
				switch s := src0.C.(type) {
				case color.Gray16:
					Y = s.Y
				case *color.Gray16:
					Y = s.Y
				}
				if Y == 0xffff {
					drawFillSrcGray16(dst0, r, Y)
				} else {
					drawFillOverGray16(dst0, r, Y)
				}
				return
			}
		} else if mask0, ok := mask.(*image.Alpha); ok {
			switch src0 := src.(type) {
			case *image.Uniform:
				var Y uint16
				switch s := src0.C.(type) {
				case color.Gray16:
					Y = s.Y
				case *color.Gray16:
					Y = s.Y
				}
				drawGlyphOverGray16(dst0, r, Y, mask0, mp)
				return
			}
		}
	}

	// Fallback to generic implementation
	// draw.DrawMask(dst, r, src, sp, mask, mp, draw.Over)
}

func drawFillSrcGray(dst *image.Gray, r image.Rectangle, y uint8) {
	width := r.Dx()
	height := r.Max.Y - r.Min.Y
	startOffset := dst.PixOffset(r.Min.X, r.Min.Y)

	firstRow := dst.Pix[startOffset : startOffset+width]
	for i := range firstRow {
		firstRow[i] = y
	}

	for dy := 1; dy < height; dy++ {
		dstOffset := startOffset + dy*dst.Stride
		copy(dst.Pix[dstOffset:dstOffset+width], firstRow)
	}
}

func drawFillOverGray(dst *image.Gray, r image.Rectangle, y uint8) {
	srcY := uint32(y) * 0x101
	width := r.Dx()
	for dy := r.Min.Y; dy < r.Max.Y; dy++ {
		i := dst.PixOffset(r.Min.X, dy)
		for dx := 0; dx < width; dx++ {
			dstY := uint32(dst.Pix[i])
			dst.Pix[i] = uint8((dstY*(0xFFFF-srcY>>8) + srcY) >> 8)
			i++
		}
	}
}

func drawGlyphOverGray(dst *image.Gray, r image.Rectangle, y uint8, mask *image.Alpha, mp image.Point) {
	srcY := uint32(y) * 0x101
	maskPix := mask.Pix[mask.PixOffset(mp.X, mp.Y):]
	maskStride := mask.Stride

	dstWidth := r.Dx()
	dstHeight := r.Dy()
	dstPix := dst.Pix
	dstOffset := dst.PixOffset(r.Min.X, r.Min.Y)

	for dy := 0; dy < dstHeight; dy++ {
		dstI := dstOffset + dy*dst.Stride
		maskI := dy * maskStride
		maskRow := maskPix[maskI : maskI+dstWidth]

		for dx := 0; dx < dstWidth; dx++ {
			if ma := uint32(maskRow[dx]) * 0x101; ma != 0 {
				dstY := uint32(dstPix[dstI])
				dstPix[dstI] = uint8((dstY*(0xFFFF-ma) + srcY*ma) >> 16)
			}
			dstI++
		}
	}
}

func drawFillSrcGray16(dst *image.Gray16, r image.Rectangle, y uint16) {
	width := r.Dx()
	y0 := uint8(y >> 8)
	y1 := uint8(y)
	for dy := r.Min.Y; dy < r.Max.Y; dy++ {
		i := dst.PixOffset(r.Min.X, dy)
		for dx := 0; dx < width; dx++ {
			dst.Pix[i] = y0
			dst.Pix[i+1] = y1
			i += 2
		}
	}
}

func drawFillOverGray16(dst *image.Gray16, r image.Rectangle, y uint16) {
	srcY := uint32(y)
	width := r.Dx()
	for dy := r.Min.Y; dy < r.Max.Y; dy++ {
		i := dst.PixOffset(r.Min.X, dy)
		for dx := 0; dx < width; dx++ {
			dstY := uint32(dst.Pix[i])<<8 | uint32(dst.Pix[i+1])
			res := (dstY*(0xFFFF-srcY) + srcY*0xFFFF) / 0xFFFF
			dst.Pix[i] = uint8(res >> 8)
			dst.Pix[i+1] = uint8(res)
			i += 2
		}
	}
}

func drawGlyphOverGray16(dst *image.Gray16, r image.Rectangle, y uint16, mask *image.Alpha, mp image.Point) {
	srcY := uint32(y)
	maskStride := mask.Stride
	maskPix := mask.Pix[mask.PixOffset(mp.X, mp.Y):]

	dstWidth := r.Dx()
	dstHeight := r.Dy()

	for dy := 0; dy < dstHeight; dy++ {
		dstI := dst.PixOffset(r.Min.X, r.Min.Y+dy)
		maskI := dy * maskStride

		for dx := 0; dx < dstWidth; dx++ {
			ma := uint32(maskPix[maskI+dx]) * 0x101
			if ma == 0 {
				dstI += 2
				continue
			}

			dstY := uint32(dst.Pix[dstI])<<8 | uint32(dst.Pix[dstI+1])
			res := (dstY*(0xFFFF-ma) + srcY*ma) / 0xFFFF
			dst.Pix[dstI] = uint8(res >> 8)
			dst.Pix[dstI+1] = uint8(res)
			dstI += 2
		}
	}
}
