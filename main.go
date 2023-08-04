package main

import (
	"image/png"
	"os"
	"path"

	"embed"

	"github.com/shadiestgoat/goSankey/drawer"
	"github.com/shadiestgoat/goSankey/parser"
	"github.com/shadiestgoat/log"
)

//go:embed config.sankey resources
var res embed.FS

func writeDir(dir string) {
	out, _ := res.ReadDir(dir)
	for _, o := range out {
		p := path.Join(dir, o.Name())
		if o.IsDir() {
			os.Mkdir(p, 0755)
			writeDir(p)
		} else {
			b, _ := res.ReadFile(p)
			os.WriteFile(p, b, 0755)
		}
	}
}

func main() {
	log.Init(log.NewLoggerPrint())

	args := os.Args[1:]
	if len(args) == 0 {
		log.Error("No command given!")
	}

	switch args[0] {
	case "init":
		writeDir(".")
	case "create":
		confName := "config.sankey"

		if len(args) >= 2 {
			confName = args[1]
		}

		c, err := parser.Parse(confName)
		log.FatalIfErr(err, "parsing the config file '%s'", confName)
		log.Success("Loaded config!")
		if c.Config.OutputName == "" {
			log.Success("Dry run complete!")
			return
		}

		log.FatalIfErr(drawer.Load(&res), "loading fonts")
		
		img := drawer.Draw(c)
		
		f, err := os.Create(c.Config.OutputName)
		log.FatalIfErr(err, "creating output file")
		png.Encode(f,img)
	}
}