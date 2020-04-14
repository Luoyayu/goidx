package utils

import "github.com/fatih/color"

func COut(s string, a ...color.Attribute) string {
	return color.New(a...).SprintFunc()(s)
}
