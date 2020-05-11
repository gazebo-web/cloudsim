package tools

import "strings"

func IsECR(image string) bool {
	return strings.Contains(image, "dkr.ecr.") && strings.Contains(image, ".amazonaws.com")
}
