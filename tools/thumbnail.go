package tools

import (
	"fmt"
	"html"
)

// GenerateThumbnailURI generates an URI of a thumbnail for a certain robot.
func GenerateThumbnailURI(url, owner, robotName string, thumbnailNo int) string {
	robotName = html.EscapeString(robotName)
	return fmt.Sprintf("%s/%s/models/%s/tip/files/thumbnails/%d.png", url, owner, robotName, thumbnailNo)
}
