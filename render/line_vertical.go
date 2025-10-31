package render

import (
	"fmt"
	"strings"
)

type LineVertical struct {
	ScaleX int
	ScaleY int
}

func (lh LineVertical) Render(bm BrailleMap) string {
	var svgEl []string
	var longestLineLen int
	for _, row := range bm {
		if rowLen := len(row); longestLineLen < rowLen {
			longestLineLen = rowLen
		}
	}

	penDown := false
	for xIdx := 0; xIdx < longestLineLen; xIdx++ {
		yBegin := 0
		xBegin, xEnd := xIdx*int(lh.ScaleY), xIdx*int(lh.ScaleY)
		for yIdx, row := range bm {
			cellIsActive := len(row) > xIdx && row[xIdx]
			strokeDone := (!cellIsActive && penDown) || (yIdx == len(bm)-1 && penDown)
			if strokeDone {
				penDown = false
				svgEl = append(svgEl,
					fmt.Sprintf(
						`<line x1="%d" y1="%d" x2="%d" y2="%d" style="stroke-linecap:round;stroke:white;stroke-width:3"/>`,
						xBegin, yBegin, xEnd, (yIdx-1)*int(lh.ScaleY)))
			} else if cellIsActive && !penDown {
				yBegin = yIdx * int(lh.ScaleX)
				penDown = true
			}
		}
	}

	svg := fmt.Sprintf(`
		<svg 
			height="320px"
			width=320px
			viewBox="%d %d %d %d" 
			xmlns="http://www.w3.org/2000/svg" 
			style="background: black; border:1px solid black"
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
