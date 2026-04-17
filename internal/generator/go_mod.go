package generator

import (
	"os"
	"strings"
)

func GetModulePath() string {
	data, _ := os.ReadFile("go.mod")
	lines := strings.Split(string(data), "\n")

	for _, l := range lines {
		if strings.HasPrefix(l, "module ") {
			return strings.TrimSpace(strings.Replace(l, "module", "", 1))
		}
	}
	return ""
}
