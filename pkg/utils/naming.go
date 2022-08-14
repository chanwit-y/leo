package utils

import (
	"strings"

	"github.com/samber/lo"
)

func toCamelCase(text string) string {
	spName := strings.Split(text, "_")
	tabelName := lo.Reduce(spName, func(r string, t string, i int) string {
		return r + t[0:1] + strings.ToLower(t[1:])
	}, "")

	return tabelName
}
