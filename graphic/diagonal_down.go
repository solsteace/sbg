package graphic

import (
	"fmt"
	"strings"
)

const DIAGONAL_DOWN = "diagonal-down"

type DiagonalDown struct {
	ScaleX int
	ScaleY int
}

func (dd DiagonalDown) SVG(bm BrailleMap) string {
	var svgEl []string
	var longestLineLen int
	for _, row := range bm {
		if rowLen := len(row); longestLineLen < rowLen {
			longestLineLen = rowLen
		}
	}

	for initXIdx := -longestLineLen; initXIdx < longestLineLen; initXIdx++ {
		penDown := false
		xIdx := initXIdx
		var xBegin, yBegin int
		for yIdx, row := range bm {
			cellIsActive := xIdx >= 0 && len(row) > xIdx && row[xIdx]
			if cellIsActive && !penDown {
				xBegin = xIdx * int(dd.ScaleX)
				yBegin = yIdx * int(dd.ScaleX)
				penDown = true
			}

			if !cellIsActive && penDown {
				penDown = false
				svgEl = append(svgEl,
					fmt.Sprintf(
						`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`,
						xBegin, yBegin, (xIdx-1)*int(dd.ScaleX), (yIdx-1)*int(dd.ScaleY)))
			} else if yIdx == len(bm)-1 && penDown {
				svgEl = append(svgEl,
					fmt.Sprintf(
						`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`,
						xBegin, yBegin, xIdx*int(dd.ScaleX), yIdx*int(dd.ScaleY)))
			}
			xIdx++
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
		-dd.ScaleX,
		-dd.ScaleY,
		dd.ScaleX*(longestLineLen+1),
		dd.ScaleY*(len(bm)+1),
		strings.Join(svgEl, ""))
	return svg
}
