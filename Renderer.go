package psvg

import (
	"github.com/go-gl/mathgl/mgl32"
	"io"
)

type (
	Renderer struct {
		data []Elem
	}
	Support interface {
		MoveTo(to mgl32.Vec2)
		LineTo(to mgl32.Vec2)
		QuadTo(p0, to mgl32.Vec2)
		CubeTo(p0, p1, to mgl32.Vec2)
		CloseTo()
	}
)

func NewRendererFromReader(src io.Reader) (*Renderer, error) {
	var data []Elem
	p := NewParser(src)
	for elem := p.Next(); elem != nil; elem = p.Next() {
		switch e := elem.(type) {
		case UnknownError:
			return nil, e
		default:
			data = append(data, e)
		}
	}
	return NewRenderer(data...), nil
}

func NewRenderer(data ...Elem) *Renderer {
	return &Renderer{
		data: data,
	}
}

func (s *Renderer) CheckError(onError func(unknown UnknownError)) (res bool) {
	res = false
	for _, d := range s.data {
		switch dt := d.(type) {
		case UnknownError:
			res = true
			if onError != nil {
				onError(dt)
			}
		}
	}
	return res
}
func (s *Renderer) CheckUnknown(onUnknown func(unknown UnknownCommand)) (res bool) {
	res = false
	for _, d := range s.data {
		switch dt := d.(type) {
		case UnknownCommand:
			res = true
			if onUnknown != nil {
				onUnknown(dt)
			}
		}
	}
	return res
}
func (s *Renderer) Render(support Support) {
	var last = mgl32.Vec2{0, 0}
	var lastQ *mgl32.Vec2
	var lastC *mgl32.Vec2
	for _, d := range s.data {
		switch dt := d.(type) {
		case ClosePath:
			support.CloseTo()
			lastQ = nil
			lastC = nil
		case MoveToAbs:
			last = dt.To
			support.MoveTo(dt.To)
			lastQ = nil
			lastC = nil
		case MoveToRel:
			last = last.Add(dt.To)
			support.MoveTo(last)
			lastQ = nil
			lastC = nil
		case LineToAbs:
			last = dt.To
			support.LineTo(dt.To)
			lastQ = nil
			lastC = nil
		case LineToRel:
			last = last.Add(dt.To)
			support.LineTo(last)
			lastQ = nil
			lastC = nil

		case CurveToCubicAbs:
			last = dt.To
			support.CubeTo(dt.P0, dt.P1, last)

			lastQ = nil
			lastC = &dt.P1
		case CurveToCubicRel:
			p0 := last.Add(dt.P0)
			p1 := last.Add(dt.P1)
			last = last.Add(dt.To)
			support.CubeTo(p0, p1, last)

			lastQ = nil
			lastC = &p1

		case CurveToQuadraticAbs:
			last = dt.To
			support.QuadTo(dt.P0, last)

			lastQ = &dt.P0
			lastC = nil
		case CurveToQuadraticRel:
			p0 := last.Add(dt.P0)
			last = last.Add(dt.To)
			support.QuadTo(p0, last)

			lastQ = &p0
			lastC = nil

		case ArcAbs:
			// TODO ArcRel support
			lastQ = nil
			lastC = nil
		case ArcRel:
			// TODO ArcRel support
			lastQ = nil
			lastC = nil

		case LineToHorizontalAbs:
			last = mgl32.Vec2{dt.X, last[1]}
			support.LineTo(last)
			lastQ = nil
			lastC = nil
		case LineToHorizontalRel:
			last = last.Add(mgl32.Vec2{dt.X, 0})
			support.LineTo(last)
			lastQ = nil
			lastC = nil

		case LineToVerticalAbs:
			last = mgl32.Vec2{last[0], dt.Y}
			support.LineTo(last)
			lastQ = nil
			lastC = nil
		case LineToVerticalRel:
			last = last.Add(mgl32.Vec2{0, dt.Y})
			support.LineTo(last)
			lastQ = nil
			lastC = nil

		case CurveToCubicSmoothAbs:
			var p0 mgl32.Vec2
			if lastC == nil {
				p0 = last
			} else {
				p0 = mirrorByPoint(*lastC, last)
			}
			p1 := dt.P1
			last = dt.To

			support.CubeTo(p0, p1, last)

			lastQ = nil
			lastC = &p1
		case CurveToCubicSmoothRel:
			var p0 mgl32.Vec2
			if lastC == nil {
				p0 = last
			} else {
				p0 = mirrorByPoint(*lastC, last)
			}
			p1 := last.Add(dt.P1)
			last = last.Add(dt.To)

			support.CubeTo(p0, p1, last)

			lastQ = nil
			lastC = &p1

		case CurveToQuadraticSmoothAbs:

			var p0 mgl32.Vec2
			if lastC == nil {
				p0 = last
			} else {
				p0 = mirrorByPoint(*lastQ, last)
			}
			last = dt.To

			support.QuadTo(p0, last)

			lastQ = &p0
			lastC = nil
		case CurveToQuadraticSmoothRel:

			var p0 mgl32.Vec2
			if lastC == nil {
				p0 = last
			} else {
				p0 = mirrorByPoint(*lastQ, last)
			}
			last = last.Add(dt.To)

			support.QuadTo(p0, last)

			lastQ = &p0
			lastC = nil

		}
	}
}

func mirrorByPoint(a, mirror mgl32.Vec2) mgl32.Vec2 {
	return mirror.Add(mirror.Sub(a))
}
