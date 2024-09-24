package renderers

import (
	"gioui/anime/canvas"
	"io"
)

const mmPerPt = 25.4 / 72.0
const ptPerMm = 72.0 / 25.4
const mmPerPx = 25.4 / 96.0

func errorWriter(err error) canvas.Writer {
	return func(w io.Writer, c *canvas.Canvas) error {
		return err
	}
}
