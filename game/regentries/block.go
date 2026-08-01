package regentries

import "image/color"

type BlockId uint16
type BlockType struct {
	IsSolid bool
	Color   color.RGBA
}
