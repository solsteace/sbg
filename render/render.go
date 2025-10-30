package render

type BrailleMap = [][]bool

// Renders the braille into HTML SVG
type Renderer interface {
	Render(bm BrailleMap) string
}
