package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/solsteace/sbg/render"
)

const (
	fLAG_SOURCE      = "--source"
	fLAG_DESTINATION = "--destination"
	fLAG_VARIATION   = "--variation"
)

type config struct {
	source      string
	destination string
	variation   string
}

func configFromArgs(osArgs []string) config {
	config := config{}

	var lastFlag string
	flags := []string{
		fLAG_SOURCE,
		fLAG_DESTINATION,
		fLAG_VARIATION}
	for _, arg := range os.Args {
		if slices.Index(flags, arg) != -1 {
			lastFlag = arg
			continue
		}

		switch lastFlag {
		case fLAG_SOURCE:
			config.source = arg
		case fLAG_DESTINATION:
			config.destination = arg
		case fLAG_VARIATION:
			config.variation = arg
		}
	}

	return config
}

func (c config) getSource() (*os.File, error) {
	source := c.source
	if source == "" {
		return os.Stdin, nil
	}

	f, err := os.Open(source)
	if err != nil {
		return nil, fmt.Errorf("config.getDestination: %w", err)
	}
	return f, nil
}

func (c config) getDestination() (*os.File, error) {
	destination := c.destination
	if destination == "" {
		return os.Stdout, nil
	}

	f, err := os.Open(destination)
	if err != nil {
		return nil, fmt.Errorf("config.getDestination: %w", err)
	}
	return f, nil
}

func (c config) getRenderer() (render.Renderer, error) {
	switch c.variation {
	case "", "line-horizontal":
		return render.LineHorizontal{ScaleX: 5, ScaleY: 5}, nil
	case "line-vertical":
		return render.LineVertical{ScaleX: 5, ScaleY: 5}, nil
	default:
		return nil, fmt.Errorf(
			"config.getRenderer: %s is not available yet. Try one from the list below:\n%s",
			c.variation,
			strings.Join(
				[]string{"line-horizontal", "line-vertical"},
				"\n"))
	}
}
