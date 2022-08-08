package color

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"os"
)

var lightGreyCode string = "\u001B[38;5;242m%s\u001B[0m"
var redCode string = "\u001B[38;5;160m%s\u001B[0m"
var greenCode string = "\u001B[38;5;2m%s\u001B[0m"
var blueCode string = "\u001B[38;5;27m%s\u001B[0m"

const NL = "\n"

var (
	// NoColor defines if the output is colorized or not. It's dynamically set to
	// false or true based on the stdout's file descriptor referring to a terminal
	// or not. This is a global option and affects all colors. For more control
	// over each color block use the methods DisableColor() individually.
	NoColor = os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))
)

func onColor(txt string, colorCode string) string {
	if NoColor {
		return txt
	} else {
		return fmt.Sprintf(colorCode, txt)
	}
}

func Green(txt string) string {
	return onColor(txt, greenCode)
}
func Red(txt string) string {
	return onColor(txt, redCode)
}
func LightGrey(txt string) string {
	return onColor(txt, lightGreyCode)
}
func Blue(txt string) string {
	return onColor(txt, blueCode)
}
