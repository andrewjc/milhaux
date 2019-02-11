package common

import (
	"strings"
	"text/scanner"
)

func ParseQuoteAwareFields(rawString string, len int) []string {
	var s scanner.Scanner
	s.Init(strings.NewReader(rawString))
	slice := make([]string, 0, len)
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		slice = append(slice, s.TokenText())
	}
	return slice
}
