package drawer

import (
	"embed"
	"fmt"
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
	fonts    []*cFontInfo
}

func Load(res *embed.FS) error {
	dir, err := os.ReadDir("resources")
	useEmbed := false
	if err != nil || len(dir) == 0 {
		dir, err = res.ReadDir("resources")
		if err != nil {
			return err
		}
		useEmbed = true
	}

	for _, f := range dir {
		name := f.Name()
		if len(name) < 4 {
			continue
		}
		path := path.Join("resources", name)
		ext := name[len(name)-3:]

		var b []byte
		if useEmbed {
			b, err = res.ReadFile(path)
			if err != nil {
				return err
			}
		} else {
			b, err = os.ReadFile(path)
			if err != nil {
				return err
			}
		}
		
		if ext == "otf" {
			_sf, err := opentype.Parse(b)
			if err != nil {
				continue
			}
			fontNormal.fonts = append(fontNormal.fonts, &cFontInfo{
				SFNT: _sf,
			})
		} else if ext == "ttf" {
			_sf, err := truetype.Parse(b)
			if err != nil {
				continue
			}
			fontNormal.fonts = append(fontNormal.fonts, &cFontInfo{
				TTF: _sf,
			})
		} else {
			continue
		}
	}

	if len(fontNormal.fonts) == 0 {
		return fmt.Errorf("couldn't load any fonts")
	}

	return nil
}