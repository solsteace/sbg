package graphic

// Renders the braille into HTML SVG
type Painter interface {
	SVG(bm BrailleMap) string
}
