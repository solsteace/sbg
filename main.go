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

	state := sTATE_READY
	var lastFlag string
	var config config
	var configSet bool
	for idx := 1; idx < len(os.Args); idx++ {
		arg := os.Args[idx]
		switch state {
		case sTATE_INSTA_END:
			switch lastFlag {
			case fLAG_HELP:
				fmt.Println(tEXT_HELP)
			case fLAG_HELP_SHORT:
				fmt.Println(tEXT_HELP_SHORT)
			}
			os.Exit(0)
		case sTATE_READY:
			switch arg {
			case fLAG_HELP, fLAG_HELP_SHORT:
				state = sTATE_INSTA_END
				lastFlag = arg
			case fLAG_CONFIG, fLAG_CONFIG_SHORT:
				state = sTATE_NEED_PARAMS
				lastFlag = arg
			}
		case sTATE_NEED_PARAMS:
			switch lastFlag {
			case fLAG_CONFIG, fLAG_CONFIG_SHORT:
				if err := config.fromFile(arg); err != nil {
					log.Fatal(err)
				}
				configSet = true
				state = sTATE_READY
			}
		}
	}

	if lastFlag == "" && !inputPiped {
		fmt.Println(tEXT_HELP_SHORT)
		os.Exit(0)
	} else if !configSet {
		if err := config.fromFile("./config.json"); err != nil {
			log.Fatal(err)
		}
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
