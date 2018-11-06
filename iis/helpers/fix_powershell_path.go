package helpers

import "strings"

func FixPowerShellPath(input string) string {
	output := input
	output = strings.Replace(output, "\\\\", "\\", -1)
	return strings.TrimSpace(output)
}
