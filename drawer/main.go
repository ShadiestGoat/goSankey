package drawer

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/AndreKR/multiface"
	"golang.org/x/image/draw"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/tools/bezier"
	"gonum.org/v1/plot/vg"

	"shadygoat.eu/goSankey/common"
)

const GENERAL_PAD float64 = 0.02
const NODE_WIDTH = 0.01
const PAD_LEFT = 0.01
const VERT_SPACE_NODES = 0.85

const HORZ_TEXT_PAD = 15
const VERT_TEXT_PAD = 5

func drawRect(img draw.Image, c color.Color, rect image.Rectangle) {
	draw.Draw(img, rect, image.NewUniform(c), image.Point{}, draw.Over)
}

func drawEmptyRect(img draw.Image, c color.Color, rect image.Rectangle, lw int) {
	drawRect(img, c, image.Rect(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Min.Y + lw))
	drawRect(img, c, image.Rect(rect.Min.X, rect.Max.Y, rect.Max.X, rect.Max.Y - lw))

	drawRect(img, c, image.Rect(rect.Min.X, rect.Min.Y, rect.Min.X + lw, rect.Max.Y))
	drawRect(img, c, image.Rect(rect.Max.X, rect.Min.Y, rect.Max.X - lw, rect.Max.Y))
}

func blend(cA, cB uint8, aA float64) uint8 {
	return uint8(round(float64(cA) * aA + float64(cB)*(1-aA)))
}

func round(f float64) int {
	return int(math.Round(f))
}

// Calculates padding, withe assumption of side padding being 0.5 of regular pad
func calcPad(available int, taken int, elements int) int {
	return int(math.Floor(float64(available - taken)/float64(elements)))
}

func nColor(c *common.Color) color.Color {
	return &color.NRGBA{
		R: c.R,
		G: c.G,
		B: c.B,
		A: 255,
	}
}

type JustifyMode int

const (
	JUSTIFY_LEFT JustifyMode = iota
	JUSTIFY_CENTER
	JUSTIFY_RIGHT
)

// Draw text within a specified width, where y0 is the top of textbox, nothing goes above it.
// returns the bottom y level
func textBox(img draw.Image, face *multiface.Face, c color.Color, str string, y0 int, x0 int, x1 int, just JustifyMode) int {
	lines := []string{}
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	
	x0 += HORZ_TEXT_PAD
	x1 -= HORZ_TEXT_PAD
	width := x1 - x0
	
	linesRaw := strings.Split(str, "\n")
	for _, text := range linesRaw {
		curL := ""
		words := strings.Split(text, " ")
	
		for _, w := range words {
			tmp := ""
			if curL == "" {
				tmp = w 
			} else {
				tmp = curL + " " + w
			}
	
			b := BoundString(face, tmp)
			if b.Dx() > width {
				if curL == "" {
					panic("Too big words :(")
				}
	
				lines = append(lines, curL)
				curL = w
			} else {
				curL = tmp
			}
		}

		lines = append(lines, curL)
	}
	

	y := y0

	for _, l := range lines {
		b := BoundString(face, l)
		var pad int
		switch just {
		case JUSTIFY_RIGHT:
			pad = width - b.Dx()
		case JUSTIFY_CENTER:
			pad = round(float64(width - b.Dx())/2)
		default:

		}

		y += b.Dy()
		addLabel(img, face, x0 + pad, y, l, c)
		y += VERT_TEXT_PAD
	}

	return y
}

type labelInfo struct {
	Txt string
	Justify JustifyMode
	X0 int
	X1 int
	Y int
	// Make Step - 0.5 to indicate this is a node label
	Step float64
	Color color.Color
}

func Draw(c *common.Chart) draw.Image {
	img := image.NewNRGBA(image.Rect(0, 0, c.Config.Width, c.Config.Height))

	drawRect(img, nColor(c.Config.Background), img.Rect)
	drawEmptyRect(img, color.White, img.Rect, 2)

	labels := []*labelInfo{}
	
	sidePadV := round(float64(c.Config.Width) * GENERAL_PAD)
	vertPadV := round(float64(c.Config.Height) * GENERAL_PAD)
	width :=  float64(c.Config.Width -  sidePadV * 2)
	height := float64(c.Config.Height - vertPadV * 2)
	leftPadV := round(width * PAD_LEFT)
	nodeW := round(width*NODE_WIDTH)
	nodeAvailableH := height * VERT_SPACE_NODES

	horzPad := calcPad(round(width - float64(leftPadV)), nodeW*len(c.Steps), len(c.Steps))
	
	x := sidePadV + leftPadV + round(float64(horzPad) * 0.5)

	nodeLocations := map[*common.Node]image.Rectangle{}

	n0Total := 0

	for i, step := range c.Steps {
		total := 0
		totalIn := 0
		
		for _, n := range step {
			total += n.TotalIn + n.TotalOut
			totalIn += n.TotalIn
		}

		availableHeight := nodeAvailableH

		if i != 0 {
			availableHeight *= float64(totalIn)/float64(n0Total)
		}

		vertPad := calcPad(round(height), round(availableHeight), len(step))
		x1 := x + nodeW
		
		y := vertPadV + round(float64(vertPad) * 0.5)
		

		for _, n := range step {
			col := nColor(n.Color)
			if i == 0 {
				n0Total += n.TotalOut
			}
			h := round(availableHeight * float64(n.TotalIn + n.TotalOut)/float64(total))
			nodeLocations[n] = image.Rect(x, y, x1, y + h)
			drawRect(img, col, nodeLocations[n])
			y += vertPad + h
		}

		x = x1 + horzPad
	}

	nodeToUsedYOut := map[*common.Node]int{}
	nodeToUsedYIn := map[*common.Node]int{}

	face := fontNormal.fetchFace(28)

	for _, conn := range c.Connections {
		ogLoc := nodeLocations[conn.Origin]
		destLoc := nodeLocations[conn.Dest]

		ogHeight := round(float64(ogLoc.Dy())*(float64(conn.Amount)/float64(conn.Origin.TotalOut)))
		ogMin := ogLoc.Min.Y + nodeToUsedYOut[conn.Origin]
		nodeToUsedYOut[conn.Origin] += ogHeight

		destHeight := round(float64(destLoc.Dy())*(float64(conn.Amount)/float64(conn.Dest.TotalIn)))
		destMin := destLoc.Min.Y + nodeToUsedYIn[conn.Dest]
		nodeToUsedYIn[conn.Dest] += destHeight

		connWidth := destLoc.Min.X - ogLoc.Max.X
		midPointOg := ogMin + round(float64(ogHeight)/2)
		midPointDest := destMin + round(float64(destHeight)/2)
		midPointX := ogLoc.Max.X + round(float64(connWidth)/2)

		cPoints := []vg.Point{{
			X: font.Length(ogLoc.Max.X),
			Y: font.Length(midPointOg),
		},{
			X: font.Length(midPointX),
			Y: font.Length(midPointOg),
		}, {
			X: font.Length(midPointX),
			Y: font.Length(midPointDest),
		}, {
			X: font.Length(destLoc.Min.X),
			Y: font.Length(midPointDest),
		}}

		dH := float64(ogHeight - destHeight)/float64(connWidth)

		curve := bezier.New(cPoints...)

		points := curve.Curve(make([]vg.Point, connWidth * 2))
		
		usedXs := map[int]bool{}

		for _, p := range points {
			x := round(float64(p.X))
			if usedXs[x] {
				continue
			}
			if x == destLoc.Min.X {
				continue
			}
			
			usedXs[x] = true
			i := ogLoc.Max.X - x
			height := int(math.RoundToEven((float64(ogHeight) + dH * float64(i))))

			minY := round(float64(p.Y)) - height/2
			
			for y := 0; y < height; y++ {
				cBg := img.At(x, minY + y)
				r, g, b, _ := cBg.RGBA()
				
				img.Set(x, minY + y, color.RGBA{
					R: blend(conn.Dest.Color.R, uint8(r), c.Config.ConnectionOpacity),
					G: blend(conn.Dest.Color.G, uint8(g), c.Config.ConnectionOpacity),
					B: blend(conn.Dest.Color.B, uint8(b), c.Config.ConnectionOpacity),
					A: 255,
				})
			} 
		}

		labels = append(labels, &labelInfo{
			Txt:     conn.Dest.Title + "\n" + fmt.Sprint(conn.Amount),
			Justify: JUSTIFY_LEFT,
			X0:      ogLoc.Max.X,
			X1:      destLoc.Min.X,
			Y:       ogMin,
			Step:    float64(conn.Origin.Step),
			Color:   nColor(conn.Dest.Color),
		})
	}

	for i, step := range c.Steps {
		for _, n := range step {
			loc := nodeLocations[n]
			prevX := sidePadV + leftPadV
			if i != 0 {
				prevX = nodeLocations[c.Steps[i-1][0]].Max.X
			}

			tot := n.TotalIn
			if n.TotalIn == 0 {
				tot = n.TotalOut
			}
			labels = append(labels, &labelInfo{
				Txt:     n.Title + "\n" + fmt.Sprint(tot),
				Justify: JUSTIFY_RIGHT,
				X0:      prevX,
				X1:      loc.Min.X,
				Y:       loc.Min.Y,
				Step:    float64(i) - 0.5,
				Color:   nColor(n.Color),
			})
		}
	}

	sort.Slice(labels, func(i, j int) bool {
		if labels[i].Step == labels[j].Step {
			return labels[i].Y < labels[j].Y
		}
		return labels[i].Step < labels[j].Step
	})

	lastStep := 0.0
	minY := 0

	for _, l := range labels {
		if lastStep != l.Step {
			minY = 0
			lastStep = l.Step
		}

		if l.Y < minY {
			l.Y = minY
		}

		minY = textBox(img, face, l.Color, l.Txt, l.Y, l.X0, l.X1, l.Justify)
	}

	f, _ := os.Create("test.png")
	png.Encode(f, img)

	return img
}