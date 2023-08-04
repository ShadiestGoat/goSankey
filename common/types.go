package common

type Node struct {
	Title string
	ID string
	Step int
	Color *Color
	TotalIn int
	TotalOut int
}

type Connection struct {
	Origin *Node
	Dest *Node
	Amount int
}

type Config struct {
	Width int
	Height int
	ConnectionOpacity float64
	Background *Color
	BackgroundIsLight bool

	OutputName string
}

type Color struct {
	R, G, B uint8
}

type Chart struct {
	Config *Config
	Steps [][]*Node
	Connections []*Connection
}