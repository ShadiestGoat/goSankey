package parser

import (
	"math"
	"strconv"
	"strings"

	"github.com/shadiestgoat/colorutils"
	"github.com/shadiestgoat/goSankey/common"
)

func config(inp []string) *common.Config {
	c := &common.Config{
		DrawBorder: true,
		BorderSize: 3,
		BorderPadding: 0.02,
		NodeWidth: 0.02,
		PadLeft: 0.01,
		VertSpaceNodes: 0.85,
		HorzTextPad: 15,
		TextLinePad: 5,
		ConnectionOpacity: 0.2,
	}

	for _, l := range inp {
		opt, v := parseOpt(l)

		switch opt {
		case "connectionopacity", "connection_opacity", "connection opacity":
			out, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			if out > 100 || out <= 0 {
				continue
			}
			c.ConnectionOpacity = float64(out)/100
		case "outputname", "output_name", "output name":
			c.OutputName = strings.TrimSpace(v)
		case "background":
			c.Background = parseColor(v)
		case "width":
			out, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			c.Width = out
		case "height":
			out, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			c.Height = out
		case "drawborder", "draw_border", "draw border":
			o := parseBool(v)
			if o != nil {
				c.DrawBorder = *o
			}
		case "bordercolor", "border_color", "border color":
			c.BorderColor = parseColor(v)
		case "bordersize", "border_size", "border size":
			out, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			c.BorderSize = out
		case "borderpadding", "border_padding", "border padding":
			out, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			c.BorderPadding = float64(out)/100
		case "nodewidth", "node_width", "node width":
			out, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			c.NodeWidth = float64(out)/100
		case "padleft", "pad_left", "pad left":
			out, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			c.PadLeft = float64(out)/100
		case "vertspacenodes", "vert_space_nodes", "vert space nodes":
			out, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			c.VertSpaceNodes = float64(out)/100
		case "horizontaltextpad", "horizontal_text_pad", "horizontal text pad":
			out, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			c.HorzTextPad = out
		case "textlinepad", "text_line_pad", "text line pad":
			out, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			c.TextLinePad = out
		default:
			continue
		}
	}

	if c.Background == nil {
		c.Background = &common.Color{
			R: 246,
			G: 248,
			B: 250,
		}
		c.BackgroundIsLight = true
	} else {
		_, _, l := colorutils.RGBToHSL(c.Background.R, c.Background.G, c.Background.B)
		c.BackgroundIsLight = l > 0.5
	}

	r, g, b := colorutils.NewContrastColor(8, c.BackgroundIsLight, colorutils.RelativeLuminosity(c.Background.R, c.Background.G, c.Background.B))
	c.BorderColor = &common.Color{
		R: r,
		G: g,
		B: b,
	}


	if c.Width == 0 {
		if c.Height == 0 {
			c.Width = 1920
		} else {
			c.Width = int(math.Round(16 * float64(c.Height)/9))
		}
	}
	if c.Height == 0 {
		if c.Width == 0 {
			c.Height = 1080
		} else {
			c.Height = int(math.Round(9 * float64(c.Width)/16))
		}
	}

	return c
}