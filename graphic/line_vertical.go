package graphic

import (
	"fmt"
	"strings"
)

const LINE_VERTICAL = "line-vertical"

type LineVertical struct {
	ScaleX int
	ScaleY int
}

func (lv LineVertical) SVG(bm BrailleMap) string {
	var svgEl []string
	var longestLineLen int
	for _, row := range bm {
		if rowLen := len(row); longestLineLen < rowLen {
			longestLineLen = rowLen
		}
	}

	for xIdx := 0; xIdx < longestLineLen; xIdx++ {
		penDown := false
		xBegin, xEnd := xIdx*int(lv.ScaleY), xIdx*int(lv.ScaleY)
		var yBegin int
		for yIdx, row := range bm {
			cellIsActive := len(row) > xIdx && row[xIdx]
			if cellIsActive && !penDown {
				yBegin = yIdx * int(lv.ScaleX)
				penDown = true
			}

			if !cellIsActive && penDown {
				penDown = false
				svgEl = append(svgEl,
					fmt.Sprintf(
						`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`,
						xBegin, yBegin, xEnd, (yIdx-1)*int(lv.ScaleY)))
			} else if yIdx == len(bm)-1 && penDown {
				svgEl = append(svgEl,
					fmt.Sprintf(
						`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`,
						xBegin, yBegin, xEnd, (yIdx)*int(lv.ScaleY)))
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
		-lv.ScaleX,
		-lv.ScaleY,
		lv.ScaleX*(longestLineLen+1),
		lv.ScaleY*(len(bm)+1),
		strings.Join(svgEl, ""))
	return svg
}
