package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/solsteace/sbg/graphic"
)

func main() {
	// https://stackoverflow.com/questions/22744443/check-if-there-is-something-to-read-on-stdin-in-golang
	// Someone says it doesn't work on Git Bash. I think it's alright for now since we're focusing on Unix environment
	var input *os.File
	inputPiped := false
	if inStat, err := os.Stdin.Stat(); err != nil {
		log.Fatal(err)
	} else if (inStat.Mode() & os.ModeCharDevice) == 0 {
		inputPiped = true
		input = os.Stdin
	}

	state := "app"
	args := os.Args
	endEarly := false
	var config config
	var lastFlag string
	for idx := 1; idx < len(args) && !endEarly; idx++ {
		switch arg := args[idx]; state {
		case "app":
			switch arg {
			case fLAG_HELP,
				fLAG_HELP_SHORT:
				state = "instaEndFlag"
				lastFlag = arg

			case fLAG_SOURCE,
				fLAG_VARIATION,
				fLAG_DESTINATION,
				fLAG_SOURCE_SHORT,
				fLAG_VARIATION_SHORT,
				fLAG_DESTINATION_SHORT:
				state = "flagOptParams"
				lastFlag = arg
			default:
				endEarly = true

				// No "flagWithParams" yet...
				// No "flagNoParams" yet...
			}

		case "instaEndFlag":
			endEarly = true

		case "flagOptParams":
			// Is it a flag?
			switch arg {
			case fLAG_SOURCE,
				fLAG_VARIATION,
				fLAG_DESTINATION,
				fLAG_SOURCE_SHORT,
				fLAG_VARIATION_SHORT,
				fLAG_DESTINATION_SHORT:

				state = "flagOptParams"
				lastFlag = arg
				continue

				// No "flagWithParam" yet...
				// No "flagNoParam" yet...
				// No "paramsInsufficient" yet...
			}

			// Perhaps its an parameter for...
			switch lastFlag {
			case fLAG_SOURCE, fLAG_SOURCE_SHORT:
				config.source = arg
				state = "paramsSufficient"
			case fLAG_DESTINATION, fLAG_DESTINATION_SHORT:
				config.destination = arg
				state = "paramsSufficient"
			case fLAG_VARIATION, fLAG_VARIATION_SHORT:
				switch arg {
				case graphic.LINE_HORIZONTAL,
					graphic.LINE_VERTICAL,
					graphic.DIAGONAL_UP,
					graphic.DIAGONAL_DOWN:

					config.variation = arg
					state = "paramsSufficient"
				default:
					log.Fatalf("Illegal value for variation: `%s`", arg)
				}
			}
		case "flagWithParams":
			continue // Nothing here yet
		case "flagNoParams":
			continue // Nothing here yet
		case "paramsSufficient":
			switch arg {
			case fLAG_SOURCE,
				fLAG_VARIATION,
				fLAG_DESTINATION,
				fLAG_SOURCE_SHORT,
				fLAG_VARIATION_SHORT,
				fLAG_DESTINATION_SHORT:

				state = "flagOptParams"
				lastFlag = arg

				// No "flagWithParam" yet...
				// No "flagNoParam" yet...
			}
		case "paramsInsufficient":
			continue // Nothing here yet
		}
	}

	switch {
	case lastFlag == "" && !inputPiped,
		state == "instaEndFlag" && lastFlag == fLAG_HELP_SHORT:
		fmt.Println(tEXT_HELP_SHORT)
		os.Exit(0)
	case state == "instaEndFlag" && lastFlag == fLAG_HELP:
		fmt.Println(tEXT_HELP)
		os.Exit(0)
	case state == "paramsInsufficient":
		log.Fatalf("Insufficient parameter for `%s`", lastFlag)
	case state == "flagWithParam":
		log.Fatalf("No parameter was supplied for `%s`", lastFlag)
	}

	if !inputPiped {
		in, err := config.getSource()
		if err != nil {
			log.Fatal(err)
		}

		if statRes, err := in.Stat(); err != nil {
			log.Fatal(err)
		} else if statRes.Size() == 0 {
			log.Fatalf("There's nothing in the stdin")
			os.Exit(0)
		}
		input = in
	}
	defer input.Close()

	bufR := bufio.NewReader(input)
	var braille graphic.BrailleMap
	for EOF := false; !EOF; {
		content, err := bufR.ReadString('\n')
		if err == io.EOF {
			EOF = true
		} else if err != nil {
			log.Fatal(err)
		}

		rows, err := graphic.Map(strings.TrimSpace(content))
		if err != nil {
			log.Fatal(err)
		}
		braille = append(braille, rows...)
	}

	renderer, err := config.getPainter()
	if err != nil {
		log.Fatal(err)
	}
	svg := renderer.SVG(braille)

	out, err := config.getDestination()
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	out.WriteString(svg)
}
