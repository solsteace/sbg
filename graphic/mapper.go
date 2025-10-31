package graphic

import "fmt"

type BrailleMap = [][]bool

// One line produces Nx4 grid with N is the length of the braille text
func Map(braille string) (BrailleMap, error) {
	lane := BrailleMap{
		[]bool{},
		[]bool{},
		[]bool{},
		[]bool{}}
	for brailleIdx, b := range braille {
		for laneIdx, _ := range lane {
			lane[laneIdx] = append(lane[laneIdx], false, false)
		}

		dotValue := b - '⠀'
		if dotValue < 0 || dotValue > 256 {
			return [][]bool{}, fmt.Errorf(
				"processBraille: Illegal character found `%c`(%d)", b, b)
		}

		offset := (brailleIdx / 3) * 2 // Braille characters are 3 bytes long
		for reducer := int32(128); reducer > 0 && dotValue > 0; reducer /= 2 {
			if dotValue < reducer {
				continue
			}
			dotValue -= reducer

			//	 1    8
			//	 2   16
			//	 4   32
			//	64  128
			switch reducer {
			case 128:
				lane[3][offset+1] = true
			case 64:
				lane[3][offset] = true
			case 32:
				lane[2][offset+1] = true
			case 16:
				lane[1][offset+1] = true
			case 8:
				lane[0][offset+1] = true
			case 4:
				lane[2][offset] = true
			case 2:
				lane[1][offset] = true
			case 1:
				lane[0][offset] = true
			}
		}
	}
	return lane, nil
}
