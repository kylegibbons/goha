package goha

import (
	"context"
	"fmt"
)

type HASP struct {
	BasePath   string
	ActivePage int
	Brightness int
	Pages      []HASPPage
}
type HASPPage struct {
	Buttons map[string]HASPButton
}

type HASPButton struct {
	OnClick            func(value string) error
	Font               int
	BackgroundColorPri int
	BackgroundColorSec int
	ForegroundColorPri int
	ForegroundColorSec int
	HorizAlign         string
	VertAlign          string
	Text               string
	Wordwrap           bool
}

func (hasp *HASP) Init(ctx context.Context) {

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			}
		}
	}()

}

func (hasp *HASP) boot() {

}

func (hasp *HASP) click() {

}

func (hasp *HASP) ColorFromHex(red int, green int, blue int) (int, error) {

	if red > 255 || green > 255 || blue > 255 {
		return 0, fmt.Errorf("value cannot be greater than 255")
	}

	return ((red >> 3) << 11) + ((green >> 2) << 5) + (blue >> 3), nil
}

func (hasp *HASP) Update() error {
	return nil
}

func (hasp *HASP) Reboot() error {
	return nil
}

func (hasp *HASP) FactoryReset() error {
	return nil
}
