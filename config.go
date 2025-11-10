package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/solsteace/sbg/graphic"
)

type painterConfig struct {
	ScaleX int `json:"scaleX"`
	ScaleY int `json:"scaleY"`
}

type config struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Variation   string `json:"variation"`

	Painter painterConfig `json:"painter"`
}

func (c *config) fromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("config.fromFile: %w", err)
	}
	defer f.Close()

	confString, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("config.fromFile: %w", err)
	}

	if err := json.Unmarshal(confString, c); err != nil {
		return fmt.Errorf("config.fromFile: %w", err)
	}
	return nil
}

func (c config) getSource() (*os.File, error) {
	source := c.Source
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
	destination := c.Destination
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
	if c.Painter.ScaleX <= 0 {
		return nil, fmt.Errorf("config.getPainter: `painter.ScaleX` should be a positive number")
	} else if c.Painter.ScaleY <= 0 {
		return nil, fmt.Errorf("config.getPainter: `painter.ScaleY` should be a positive number")
	}

	switch c.Variation {
	case graphic.LINE_HORIZONTAL:
		return graphic.LineHorizontal{
			ScaleX: c.Painter.ScaleX,
			ScaleY: c.Painter.ScaleY}, nil
	case graphic.LINE_VERTICAL:
		return graphic.LineVertical{
			ScaleX: c.Painter.ScaleX,
			ScaleY: c.Painter.ScaleY}, nil
	case graphic.DIAGONAL_UP:
		return graphic.DiagonalUp{
			ScaleX: c.Painter.ScaleX,
			ScaleY: c.Painter.ScaleY}, nil
	case graphic.DIAGONAL_DOWN:
		return graphic.DiagonalDown{
			ScaleX: c.Painter.ScaleX,
			ScaleY: c.Painter.ScaleY}, nil
	default:
		return nil, fmt.Errorf(
			"config.getPainter: %s is not available yet. Try one from the list below:\n%s",
			c.Variation,
			strings.Join([]string{
				graphic.LINE_HORIZONTAL,
				graphic.LINE_VERTICAL,
				graphic.DIAGONAL_UP,
				graphic.DIAGONAL_DOWN}, "\n"))
	}
}
