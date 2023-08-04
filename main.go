package main

import (
	"github.com/shadiestgoat/goSankey/drawer"
	"github.com/shadiestgoat/goSankey/parser"
)

func main() {
	c := parser.Parse()
	drawer.Draw(c)
}