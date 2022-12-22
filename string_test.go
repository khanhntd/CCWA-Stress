package main

import (
	"fmt"
	"strings"
	"testing"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func BenchmarkStringOps(b *testing.B) {
	foldCaser := cases.Fold()
	lowerCaser := cases.Lower(language.English)

	tests := []struct {
		description   string
		first, second string
	}{
		{
			description: "both strings equal",
			first:       "aaaa",
			second:      "aaaa",
		},
	}

	for _, tt := range tests {
		b.Run(fmt.Sprintf("%s::equality op", tt.description), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchmarkStringEqualsOperation(tt.first, tt.second, b)
			}
		})

		b.Run(fmt.Sprintf("%s::strings equal fold", tt.description), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchmarkStringsEqualFold(tt.first, tt.second, b)
			}
		})

		b.Run(fmt.Sprintf("%s::fold caser", tt.description), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchmarkStringsFoldCaser(tt.first, tt.second, foldCaser, b)
			}
		})

		b.Run(fmt.Sprintf("%s::lower caser", tt.description), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchmarkStringsLowerCaser(tt.first, tt.second, lowerCaser, b)
			}
		})
	}
}

func benchmarkStringEqualsOperation(first, second string, b *testing.B) bool {
	return first == second
}

func benchmarkStringsEqualFold(first, second string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = strings.EqualFold(first, second)
	}
}

func benchmarkStringsFoldCaser(first, second string, caser cases.Caser, b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = caser.String(first) == caser.String(second)
	}
}

func benchmarkStringsLowerCaser(first, second string, caser cases.Caser, b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = caser.String(first) == caser.String(second)
	}
}
