package drawer

import (
	"image"
	"image/color"
	"math"

	"github.com/AndreKR/multiface"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func (f cFont) fetchFace(size float64) *multiface.Face {
	mFace := multiface.Face{}

	for _, fnt := range f.fonts {
		if fnt.TTF != nil {
			fc := truetype.NewFace(fnt.TTF, &truetype.Options{
				Size:    size,
				DPI:     72,
				Hinting: font.HintingFull,
			})
			mFace.AddTruetypeFace(fc, fnt.TTF)
		} else if fnt.SFNT != nil {
			fc, err := opentype.NewFace(fnt.SFNT, &opentype.FaceOptions{
				Size:    size,
				DPI:     72,
				Hinting: font.HintingFull,
			})
			if err != nil {
				continue
			}
			mFace.AddFace(fc)
		}
	}
	return &mFace
}

// Code taken from the ebiten project! Big thanks to them! (Also Idk what half of this does)
func fixed26_6ToFloat64(x fixed.Int26_6) float64 {
	return float64(x>>6) + float64(x&((1<<6)-1))/float64(1<<6)
}

// Code taken from the ebiten project! Big thanks to them! (Also Idk what half of this does)
func BoundString(face *multiface.Face, text string) image.Rectangle {
	m := face.Metrics()
	faceHeight := m.Height

	fx, fy := fixed.I(0), fixed.I(0)
	prevR := rune(-1)

	var bounds fixed.Rectangle26_6

	for _, r := range text {
		if prevR >= 0 {
			fx += face.Kern(prevR, r)
		}
		if r == '\n' {
			fx = fixed.I(0)
			fy += faceHeight
			prevR = rune(-1)
			continue
		}

		b, a, _ := face.GlyphBounds(r)
		b.Min.X += fx
		b.Max.X += fx
		b.Min.Y += fy
		b.Max.Y += fy
		bounds = bounds.Union(b)
		fx += a

		prevR = r
	}

	return image.Rect(
		int(math.Floor(fixed26_6ToFloat64(bounds.Min.X))),
		int(math.Floor(fixed26_6ToFloat64(bounds.Min.Y))),
		int(math.Ceil(fixed26_6ToFloat64(bounds.Max.X))),
		int(math.Ceil(fixed26_6ToFloat64(bounds.Max.Y))),
	)
}

func addLabel(img draw.Image, face *multiface.Face, x, y int, label string, color color.Color) {
	point := fixed.P(x, y)

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color),
		Face: face,
		Dot:  point,
	}

	d.DrawString(label)
}
