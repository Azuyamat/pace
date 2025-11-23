package logger

import (
	"strconv"
	"strings"
)

type Color string
type BackgroundColor string

const colorPrefix = "\033["
const colorSuffix = "m"

const (
	ColorReset  Color = colorPrefix + "0" + colorSuffix
	ColorRed    Color = colorPrefix + "31" + colorSuffix
	ColorGreen  Color = colorPrefix + "32" + colorSuffix
	ColorYellow Color = colorPrefix + "33" + colorSuffix
	ColorBlue   Color = colorPrefix + "34" + colorSuffix
	ColorPurple Color = colorPrefix + "35" + colorSuffix
	ColorCyan   Color = colorPrefix + "36" + colorSuffix
	ColorWhite  Color = colorPrefix + "37" + colorSuffix
	ColorGray   Color = colorPrefix + "90" + colorSuffix
)

func (c Color) Wrap(text string) string {
	return string(ColorReset) + string(c) + text + string(ColorReset)
}

func (c Color) Background() BackgroundColor {
	return BackgroundColor(c.transform(func(code int) int { return code + 10 }))
}

func (c Color) Bright() Color {
	return c.transform(func(code int) int {
		if code >= 30 && code <= 37 {
			return code + 60
		}
		return code
	})
}

func (c Color) Dark() Color {
	return c.transform(func(code int) int {
		if code >= 90 && code <= 97 {
			return code - 60
		}
		return code
	})
}

func (c Color) Dim() Color {
	if c == ColorReset {
		return ColorReset
	}
	s := string(c)
	s = strings.TrimPrefix(s, colorPrefix)
	s = strings.TrimSuffix(s, colorSuffix)
	return Color(colorPrefix + s + ";2" + colorSuffix)
}

func (c Color) transform(fn func(int) int) Color {
	if c == ColorReset {
		return ColorReset
	}

	s := string(c)
	s = strings.TrimPrefix(s, colorPrefix)
	s = strings.TrimSuffix(s, colorSuffix)

	code, err := strconv.Atoi(s)
	if err != nil {
		return ColorReset
	}

	return Color(colorPrefix + strconv.Itoa(fn(code)) + colorSuffix)
}

func (b BackgroundColor) Wrap(text string) string {
	return string(ColorReset) + string(b) + text + string(ColorReset)
}

func (b BackgroundColor) Bright() BackgroundColor {
	return BackgroundColor(transformBackgroundCode(string(b), func(code int) int {
		if code >= 40 && code <= 47 {
			return code + 60
		}
		return code
	}))
}

func (b BackgroundColor) Dark() BackgroundColor {
	return BackgroundColor(transformBackgroundCode(string(b), func(code int) int {
		if code >= 100 && code <= 107 {
			return code - 60
		}
		return code
	}))
}

func transformBackgroundCode(color string, fn func(int) int) string {
	s := strings.TrimPrefix(color, colorPrefix)
	s = strings.TrimSuffix(s, colorSuffix)

	code, err := strconv.Atoi(s)
	if err != nil {
		return string(ColorReset)
	}

	return colorPrefix + strconv.Itoa(fn(code)) + colorSuffix
}
