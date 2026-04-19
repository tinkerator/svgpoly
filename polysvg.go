package svgpoly

import (
	"errors"
	"fmt"
	"io"

	"zappem.net/pub/graphics/svgof"
	"zappem.net/pub/math/polygon"
)

// ErrNoData is an error returned by (*Shapes).SVG() when p is nil, or
// empty.
var ErrNoData = errors.New("no shapes data")

func min(as ...float64) (ans float64) {
	for i, a := range as {
		if i == 0 || a < ans {
			ans = a
		}
	}
	return
}

func max(as ...float64) (ans float64) {
	for i, a := range as {
		if i == 0 || a > ans {
			ans = a
		}
	}
	return
}

// SVG generates a stylized SVG from the *polygon.Shapes, p. The color
// choice is "blue" outlines polygons that are filled with "cyan", and
// holes are outlined in "red" and filled with "white".
func SVG(p *polygon.Shapes, out io.Writer, scribe float64, lines []polygon.Line) error {
	if p == nil || len(p.P) == 0 {
		return ErrNoData
	}
	ll, tr := p.BB()
	ll.X -= 1
	ll.Y -= 1
	tr.X += 1
	tr.Y += 1

	canvas := svgof.New(out)
	canvas.Decimals = 3

	canvas.StartviewUnit(tr.X-ll.X, tr.Y-ll.Y, "mm", ll.X, ll.Y, tr.X-ll.X, tr.Y-ll.Y)
	canvas.Rect(ll.X, ll.Y, tr.X-ll.X, tr.Y-ll.Y, `fill="white"`)

	for _, s := range p.P {
		xs := []float64{s.PS[len(s.PS)-1].X}
		ys := []float64{s.PS[len(s.PS)-1].Y}
		for _, pt := range s.PS {
			xs = append(xs, pt.X)
			ys = append(ys, pt.Y)
		}
		if s.Hole {
			canvas.Polyline(xs, ys, fmt.Sprintf(`fill="white" stroke="red" stroke-width="%.3f"`, scribe))
		} else {
			canvas.Polyline(xs, ys, fmt.Sprintf(`fill="cyan" stroke="blue" stroke-width="%.3f"`, scribe))
		}
	}
	for _, line := range lines {
		xs := []float64{line.From.X, line.To.X}
		ys := []float64{line.From.Y, line.To.Y}
		canvas.Polyline(xs, ys, fmt.Sprintf(`fill="none" stroke="green" stroke-width="%.3f"`, scribe))
	}
	canvas.End()

	return nil
}
