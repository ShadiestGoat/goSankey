package drawer

import (
	"os"
	"path"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

var fontNormal = &cFont{}

type cFontInfo struct {
	TTF  *truetype.Font
	SFNT *sfnt.Font
}

type cFont struct {
	fontName string
	fonts    []*cFontInfo
}

func initFont(fontName string, fontDir string, res *cFont) {
	baseDir := path.Join("resources", fontDir)
	dir, err := os.ReadDir(baseDir)
	if err != nil {
		return
	}

	res.fontName = fontName

	for _, f := range dir {
		name := f.Name()
		if len(name) < 4 {
			continue
		}
		path := path.Join("resources", fontDir, name)
		ext := name[len(name)-3:]
		b, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		
		if ext == "otf" {
			_sf, err := opentype.Parse(b)
			if err != nil {
				continue
			}
			res.fonts = append(res.fonts, &cFontInfo{
				SFNT: _sf,
			})
		} else if ext == "ttf" {
			_sf, err := truetype.Parse(b)
			if err != nil {
				continue
			}
			res.fonts = append(res.fonts, &cFontInfo{
				TTF: _sf,
			})
		} else {
			panic("Couldn't load a '" + ext + "'!")
		}
	}
}

func init() {
	initFont("normal font", "grotesk", fontNormal)
}