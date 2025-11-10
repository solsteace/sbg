package graphic

import (
	"fmt"
	"sync"
)

const LINE_VERTICAL = "line-vertical"

type LineVertical struct {
	ScaleX int
	ScaleY int
}

func (lv LineVertical) SVG(bm BrailleMap) string {
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
	columns := make(chan int)
	for pIdx := 0; pIdx < nPainters; pIdx++ {
		painters.Add(1)
		go func(id int) {
			defer painters.Done()

			for cIdx := range columns {
				penDown := false
				xBegin, xEnd := cIdx*int(lv.ScaleX), cIdx*int(lv.ScaleX)
				var yBegin int
				for yIdx, row := range bm {
					cellIsActive := len(row) > cIdx && row[cIdx]
					if cellIsActive && !penDown {
						yBegin = yIdx * int(lv.ScaleY)
						penDown = true
					}

					if !cellIsActive && penDown {
						penDown = false
						svgEl <- fmt.Sprintf(
							`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`,
							xBegin, yBegin, xEnd, (yIdx-1)*int(lv.ScaleY))
					} else if yIdx == len(bm)-1 && penDown {
						svgEl <- fmt.Sprintf(
							`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`,
							xBegin, yBegin, xEnd, (yIdx)*int(lv.ScaleY))
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
		for idx := 0; idx < longestLineLen; idx++ {
			columns <- idx
		}
		close(columns)
	}()

	svg := fmt.Sprintf(`
		<svg 
			height="320px"
			width="320px"
			viewBox="%d %d %d %d" 
			xmlns="http://www.w3.org/2000/svg" 
			style="background: black; border:1px solid black; stroke-linecap:round; stroke: white; stroke-width: 3;"
		>`,
		-lv.ScaleX,
		-lv.ScaleY,
		lv.ScaleX*(longestLineLen+1),
		lv.ScaleY*(len(bm)+1))
	for el := range svgEl {
		svg += el
	}
	return fmt.Sprintf("%s </svg>", svg)
}
