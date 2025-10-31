package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/solsteace/sbg/graphic"
)

type config struct {
	source      string
	destination string
	variation   string
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

	f, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return nil, fmt.Errorf("config.getDestination: %w", err)
	}
	return f, nil
}

func (c config) getPainter() (graphic.Painter, error) {
	switch c.variation {
	case "", graphic.LINE_HORIZONTAL:
		return graphic.LineHorizontal{ScaleX: 5, ScaleY: 5}, nil
	case graphic.LINE_VERTICAL:
		return graphic.LineVertical{ScaleX: 5, ScaleY: 5}, nil
	case graphic.DIAGONAL_UP:
		return graphic.DiagonalUp{ScaleX: 5, ScaleY: 5}, nil
	case graphic.DIAGONAL_DOWN:
		return graphic.DiagonalDown{ScaleX: 5, ScaleY: 5}, nil
	default:
		return nil, fmt.Errorf(
			"config.getRenderer: %s is not available yet. Try one from the list below:\n%s",
			c.variation,
			strings.Join(
				[]string{"line-horizontal", "line-vertical"},
				"\n"))
	}
}
