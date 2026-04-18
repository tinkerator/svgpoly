// Package svgpoly reads SVG files and represents them as a selection
// of polygon.Shapes.
package svgpoly

import (
	"fmt"
	"log"
	"os"

	"zappem.net/pub/graphics/polymark"
	"zappem.net/pub/graphics/svger"
	"zappem.net/pub/math/polygon"
)

// Debug causes extra logging for debugging purposes.
var Debug = false

// decodeSVG unravels the content of the SVG into sequences of svger.DrawingInstructions.
func decodeSVG(s *svger.Svg) (dis []*svger.DrawingInstruction, err error) {
	ins := s.ParseDrawingInstructions()
	for i := range ins {
		if i.Error != nil {
			err = i.Error
			return
		}
		dis = append(dis, i)
	}
	return
}

// disToPoints converts SVG points to polygon.Points.
func disToPoints(dis []*svger.DrawingInstruction) (pts []polygon.Point) {
	for _, a := range dis {
		pts = append(pts, polygon.Point{X: a.M[0], Y: a.M[1]})
	}
	return
}

// trace converts all of the solid SVG shapes (as determined by
// matching one of the supplied colors) and separates out contours
// lines as cuts. Lines not closed and with no width are silently
// dropped.  An absence of any specified colors, defaults to looking
// only for black = "#000000".
func trace(dis []*svger.DrawingInstruction, scribe float64, colors ...string) (shapes, cuts *polygon.Shapes, err error) {
	cols := make(map[string]bool)
	if len(colors) == 0 {
		cols["#000000"] = true
	} else {
		for _, color := range colors {
			cols[color] = true
		}
	}
	pen := polymark.Pen{
		Scribe: scribe,
	}
	for base, i := 0, 0; i < len(dis); i++ {
		a := dis[i]
		if a.Kind == svger.PaintInstruction {
			switch dis[base].Kind {
			case svger.MoveInstruction:
				if f := a.Fill; f != nil && *f != "none" {
					if dis[i-1].Kind != svger.CloseInstruction {
						err = fmt.Errorf("unexpected filled object is not closed: %v %v", dis[i-1], *f)
						return
					}
					if cols[*f] {
						shapes = shapes.Builder(disToPoints(dis[base : i-1])...)
					} else if *f == "none" {
						cuts = cuts.Builder(disToPoints(dis[base : i-1])...)
					} // ignore other color choices (as painted holes)
					break
				}
				// No substance
				if a.StrokeWidth == nil || *a.StrokeWidth == 0 {
					break
				}
				// "i != base" because the .Kind values are different...
				pts := dis[base:i]
				if dis[i-1].Kind == svger.CloseInstruction {
					pts = append([]*svger.DrawingInstruction{dis[i-2]}, dis[base:i-1]...)
				}
				shapes = pen.Line(shapes, disToPoints(pts), *a.StrokeWidth, a.StrokeLineJoin != nil && *a.StrokeLineJoin == "round", a.StrokeLineCap != nil && *a.StrokeLineCap == "round")
			case svger.CircleInstruction:
				b := dis[base]
				if f := a.Fill; f != nil {
					pt := polygon.Point{
						X: b.M[0],
						Y: b.M[1],
					}
					if *f == "none" {
						cuts = pen.Circle(cuts, pt, *b.Radius)
					} else if cols[*f] {
						shapes = pen.Circle(shapes, pt, *b.Radius)
					} // ignore other colors
				}
			default:
				log.Printf("uncharacterized dis[%d:%d]", base, i)
			}
			base = i + 1
		}
	}
	shapes.Reorder()
	cuts.Reorder()
	return
}

// LoadSVG parses an svg file into memory and returns a set of polygon
// Shapes, and a separate set of unfilled cuts.
func LoadSVG(path string, scribe float64, colors ...string) (shapes, cuts *polygon.Shapes, err error) {
	if Debug {
		log.Printf("parsing %q", path)
	}
	f, err2 := os.Open(path)
	if err2 != nil {
		err = err2
		return
	}
	s, err2 := svger.ParseSvgFromReader(f, path, 1)
	f.Close()
	if err2 != nil {
		err = fmt.Errorf("failed to parse %q: %v", path, err2)
		return
	}
	if Debug {
		log.Printf("SVG: %#v", s)
	}
	dis, err2 := decodeSVG(s)
	if err2 != nil {
		err = fmt.Errorf("failed to fully decode SVG [%q]: %v", path, err2)
	}
	if Debug {
		log.Printf("decoded %q", path)
	}
	// Convert the SVG shapes into polygon.Shapes.
	shapes, cuts, err = trace(dis, scribe)
	return
}
