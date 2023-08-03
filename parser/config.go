package parser

import (
	"math"
	"strconv"

	"github.com/shadiestgoat/colorutils"
	"shadygoat.eu/goSankey/common"
)

func config(inp []string) *common.Config {
	c := &common.Config{}

	bgLightIsSet := false

	for _, l := range inp {
		opt, v := parseOpt(l)

		switch opt {
		case "connectionopacity", "connection_opacity":
			out, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			if out > 100 || out <= 0 {
				continue
			}
			c.ConnectionOpacity = float64(out)/100
		case "background":
			c.Background = parseColor(v)
		case "backgroundislight", "background_is_light":
			out := parseBool(v)
			if out != nil {
				c.BackgroundIsLight = *out
				bgLightIsSet = true
			}
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
		bgLightIsSet = true
	}

	if !bgLightIsSet {
		_, _, l := colorutils.RGBToHSL(c.Background.R, c.Background.G, c.Background.B)
		c.BackgroundIsLight = l > 0.5
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
	if c.ConnectionOpacity == 0 {
		c.ConnectionOpacity = 0.4
	}

	return c
}