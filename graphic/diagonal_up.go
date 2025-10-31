package graphic

import (
	"fmt"
	"strings"
)

const DIAGONAL_UP = "diagonal-up"

type DiagonalUp struct {
	ScaleX int
	ScaleY int
}

func (du DiagonalUp) SVG(bm BrailleMap) string {
	var svgEl []string
	var longestLineLen int
	for _, row := range bm {
		if rowLen := len(row); longestLineLen < rowLen {
			longestLineLen = rowLen
		}
	}

	for initXIdx := longestLineLen*2 - 1; initXIdx > 0; initXIdx-- {
		penDown := false
		xIdx := initXIdx
		var xBegin, yBegin int
		for yIdx, row := range bm {
			cellIsActive := xIdx >= 0 && len(row) > xIdx && row[xIdx]
			if cellIsActive && !penDown {
				xBegin = xIdx * int(du.ScaleX)
				yBegin = yIdx * int(du.ScaleX)
				penDown = true
			}

			if !cellIsActive && penDown {
				penDown = false
				svgEl = append(svgEl,
					fmt.Sprintf(
						`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`,
						xBegin, yBegin, (xIdx+1)*int(du.ScaleX), (yIdx-1)*int(du.ScaleY)))
			} else if yIdx == len(bm)-1 && penDown {
				svgEl = append(svgEl,
					fmt.Sprintf(
						`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`,
						xBegin, yBegin, xIdx*int(du.ScaleX), yIdx*int(du.ScaleY)))
			}
			xIdx--
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
		-du.ScaleX,
		-du.ScaleY,
		du.ScaleX*(longestLineLen+1),
		du.ScaleY*(len(bm)+1),
		strings.Join(svgEl, ""))
	return svg
}
