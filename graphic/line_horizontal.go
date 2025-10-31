package graphic

import (
	"fmt"
	"strings"
)

type LineHorizontal struct {
	ScaleX int
	ScaleY int
}

func (lh LineHorizontal) SVG(bm BrailleMap) string {
	var svgEl []string
	var longestLineLen int
	for idx, row := range bm {
		penDown := false
		xBegin := 0
		yBegin, yEnd := idx*int(lh.ScaleY), idx*int(lh.ScaleY)

		if rowLen := len(row); rowLen > longestLineLen {
			longestLineLen = rowLen
		}

		for idx, cellIsActive := range row {
			if cellIsActive && !penDown {
				xBegin = idx * int(lh.ScaleX)
				penDown = true
			}

			if !cellIsActive && penDown {
				penDown = false
				svgEl = append(svgEl,
					fmt.Sprintf(
						`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`,
						xBegin, yBegin, (idx-1)*int(lh.ScaleX), yEnd))
			} else if idx == len(bm)-1 && penDown {
				svgEl = append(svgEl,
					fmt.Sprintf(
						`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`,
						xBegin, yBegin, (idx)*int(lh.ScaleX), yEnd))
			}
		}
	}

	svg := fmt.Sprintf(`
		<svg 
			height="320px"
			width=320px
			viewBox="%d %d %d %d" 
			xmlns="http://www.w3.org/2000/svg" 
			style="background: black; border:1px solid black; stroke-linecap:round; stroke: white; stroke-width: 3;"
		>
			%s
		</svg>`,
		-lh.ScaleX,
		-lh.ScaleY,
		lh.ScaleX*(longestLineLen+1),
		lh.ScaleY*(len(bm)+1),
		strings.Join(svgEl, ""))
	return svg
}
