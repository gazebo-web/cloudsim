package tools

import "strings"

// IsECR checks if the given image is a valid ECR image.
func IsECR(image string) bool {
	return strings.Contains(image, "dkr.ecr.") && strings.Contains(image, ".amazonaws.com")
}
