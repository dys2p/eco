package image

import (
	"fmt"
	"path"
)

type Image struct {
	Store Store
	Dir   string // relative to Store.Dir, becomes part of URL
	Name  string
	Size  int64
}

func (img Image) Path(maxSide int) string {
	return fmt.Sprintf("%s?m=%d", path.Join(img.Dir, img.Name), maxSide)
}
