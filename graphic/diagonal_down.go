package graphic

import (
	"fmt"
	"sync"
)

const DIAGONAL_DOWN = "diagonal-down"

type DiagonalDown struct {
	ScaleX int
	ScaleY int
}

func (dd DiagonalDown) SVG(bm BrailleMap) string {
	var nPainters int
	var longestLineLen int
	for _, row := range bm {
		if nPainters < 4 { // cap
			nPainters++
		}

		if rowLen := len(row); longestLineLen < rowLen {
			longestLineLen = rowLen
		}
	}

	painters := sync.WaitGroup{}
	svgEl := make(chan string, nPainters)
	diagonals := make(chan int)
	for pIdx := 0; pIdx < nPainters; pIdx++ {
		painters.Add(1)
		go func(id int) {
			defer painters.Done()

			for startingIdx := range diagonals {
				penDown := false
				xIdx := startingIdx
				var xBegin, yBegin int
				for yIdx, row := range bm {
					cellIsActive := xIdx >= 0 && len(row) > xIdx && row[xIdx]
					if cellIsActive && !penDown {
						xBegin = xIdx * int(dd.ScaleX)
						yBegin = yIdx * int(dd.ScaleY)
						penDown = true
					}

					if !cellIsActive && penDown {
						penDown = false
						svgEl <- fmt.Sprintf(
							`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`,
							xBegin, yBegin, (xIdx-1)*int(dd.ScaleX), (yIdx-1)*int(dd.ScaleY))
					} else if yIdx == len(bm)-1 && penDown {
						svgEl <- fmt.Sprintf(
							`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`,
							xBegin, yBegin, xIdx*int(dd.ScaleX), yIdx*int(dd.ScaleY))
					}
					xIdx++
				}
			}
		}(pIdx)
	}

	go func() {
		painters.Wait()
		close(svgEl)
	}()

	go func() {
		for startingIdx := -longestLineLen; startingIdx < longestLineLen; startingIdx++ {
			diagonals <- startingIdx
		}
		close(diagonals)
	}()

	svg := fmt.Sprintf(`
		<svg 
			height="320px"
			width="320px"
			viewBox="%d %d %d %d" 
			xmlns="http://www.w3.org/2000/svg" 
			style="background: black; border:1px solid black; stroke-linecap:round; stroke: white; stroke-width: 3;"
		>`,
		-dd.ScaleX,
		-dd.ScaleY,
		dd.ScaleX*(longestLineLen+1),
		dd.ScaleY*(len(bm)+1))
	for el := range svgEl {
		svg += el
	}
	return fmt.Sprintf("%s </svg>", svg)
}
