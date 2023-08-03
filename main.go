package main

import (
	"shadygoat.eu/goSankey/drawer"
	"shadygoat.eu/goSankey/parser"
)

func main() {
	c := parser.Parse()
	drawer.Draw(c)
}