package graphic

import (
	"fmt"
	"sync"
)

const LINE_HORIZONTAL = "line-horizontal"

type LineHorizontal struct {
	ScaleX int
	ScaleY int
}

func (lh LineHorizontal) SVG(bm BrailleMap) string {
	var nPainters int
	var longestLineLen int
	for _, row := range bm {
		if nPainters < 4 { // cap
			nPainters++
		}

		if rowLen := len(row); rowLen > longestLineLen {
			longestLineLen = rowLen
		}
	}

	// Idea: each worker handles one line at a time
	painters := sync.WaitGroup{}
	svgEl := make(chan string, nPainters)
	rows := make(chan struct {
		nRow    int
		content []bool
	})
	for pIdx := 0; pIdx < nPainters; pIdx++ {
		painters.Add(1)
		go func(id int) {
			defer painters.Done()

			for row := range rows {
				penDown := false
				yBegin, yEnd := row.nRow*int(lh.ScaleY), row.nRow*int(lh.ScaleY)
				var xBegin int
				for xIdx, cellIsActive := range row.content {
					if cellIsActive && !penDown {
						xBegin = xIdx * int(lh.ScaleX)
						penDown = true
					}

					if !cellIsActive && penDown {
						penDown = false
						svgEl <- fmt.Sprintf(
							`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`,
							xBegin, yBegin, (xIdx-1)*int(lh.ScaleX), yEnd)
					} else if xIdx == len(bm)-1 && penDown {
						svgEl <- fmt.Sprintf(
							`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`,
							xBegin, yBegin, (xIdx)*int(lh.ScaleX), yEnd)
					}
				}
			}
		}(pIdx)
	}
	go func() {
		painters.Wait()
		close(svgEl)
	}()

	go func() {
		for idx, row := range bm {
			rows <- struct {
				nRow    int
				content []bool
			}{
				idx, row}
		}
		close(rows)
	}()

	svg := fmt.Sprintf(`
		<svg 
			height="320px"
			width="320px"
			viewBox="%d %d %d %d" 
			xmlns="http://www.w3.org/2000/svg" 
			style="background: black; border:1px solid black; stroke-linecap:round; stroke: white; stroke-width: 3;"
		>`,
		-lh.ScaleX,
		-lh.ScaleY,
		lh.ScaleX*(longestLineLen+1),
		lh.ScaleY*(len(bm)+1))
	for el := range svgEl {
		svg += el
	}
	return fmt.Sprintf("%s </svg>", svg)
}
