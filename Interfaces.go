package psvg

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/iamGreedy/canvas/psvg/seg"
	"fmt"
)

// refer by
// https://www.w3.org/TR/SVG/paths.html#DOMInterfaces
type (
	Elem interface {
		Type() seg.Type
	}
	// Kind Of Error
	UnknownError struct {
		From string
		Err error
	}
	UnknownCommand struct {
		Command string
	}

	// https://www.w3.org/TR/SVG/paths.html#Interfaceseg.TypeClosePath
	ClosePath struct {}

	// https://www.w3.org/TR/SVG/paths.html#Interfaceseg.TypeMovetoAbs
	MoveToAbs struct {
		To mgl32.Vec2
	}
	// https://www.w3.org/TR/SVG/paths.html#Interfaceseg.TypeMovetoRel
	MoveToRel struct {
		To mgl32.Vec2
	}

	// https://www.w3.org/TR/SVG/paths.html#Interfaceseg.TypeLinetoAbs
	LineToAbs struct {
		To mgl32.Vec2
	}
	// https://www.w3.org/TR/SVG/paths.html#Interfaceseg.TypeLinetoRel
	LineToRel struct {
		To mgl32.Vec2
	}

	// https://www.w3.org/TR/SVG/paths.html#Interfaceseg.TypeCurvetoCubicAbs
	CurveToCubicAbs struct {
		P0, P1, To mgl32.Vec2
	}
	// https://www.w3.org/TR/SVG/paths.html#Interfaceseg.TypeCurvetoCubicRel
	CurveToCubicRel struct {
		P0, P1, To mgl32.Vec2
	}

	// https://www.w3.org/TR/SVG/paths.html#Interfaceseg.TypeCurvetoQuadraticAbs
	CurveToQuadraticAbs struct {
		P0, To mgl32.Vec2
	}
	// https://www.w3.org/TR/SVG/paths.html#Interfaceseg.TypeCurvetoQuadraticRel
	CurveToQuadraticRel struct {
		P0, To mgl32.Vec2
	}

	// https://www.w3.org/TR/SVG/paths.html#Interfaceseg.TypeArcAbs
	ArcAbs struct {
		To       mgl32.Vec2
		Radius   mgl32.Vec2
		Angle    float32
		LargeArc bool
		Sweep    bool
	}
	// https://www.w3.org/TR/SVG/paths.html#Interfaceseg.TypeArcRel
	ArcRel struct {
		To       mgl32.Vec2
		Radius   mgl32.Vec2
		Angle    float32
		LargeArc bool
		Sweep    bool
	}

	// https://www.w3.org/TR/SVG/paths.html#Interfaceseg.TypeLinetoHorizontalAbs
	LineToHorizontalAbs struct {
		X float32
	}
	// https://www.w3.org/TR/SVG/paths.html#Interfaceseg.TypeLinetoHorizontalRel
	LineToHorizontalRel struct {
		X float32
	}

	// https://www.w3.org/TR/SVG/paths.html#Interfaceseg.TypeLinetoVerticalAbs
	LineToVerticalAbs struct {
		Y float32
	}
	// https://www.w3.org/TR/SVG/paths.html#Interfaceseg.TypeLinetoVerticalRel
	LineToVerticalRel struct {
		Y float32
	}

	// https://www.w3.org/TR/SVG/paths.html#Interfaceseg.TypeCurvetoCubicSmoothAbs
	CurveToCubicSmoothAbs struct {
		P0 mgl32.Vec2
		To mgl32.Vec2
	}
	// https://www.w3.org/TR/SVG/paths.html#Interfaceseg.TypeCurvetoCubicSmoothRel
	CurveToCubicSmoothRel struct {
		P0 mgl32.Vec2
		To mgl32.Vec2
	}

	// https://www.w3.org/TR/SVG/paths.html#Interfaceseg.TypeCurvetoQuadraticSmoothAbs
	CurveToQuadraticSmoothAbs struct {
		To mgl32.Vec2
	}
	// https://www.w3.org/TR/SVG/paths.html#Interfaceseg.TypeCurvetoQuadraticSmoothRel
	CurveToQuadraticSmoothRel struct {
		To mgl32.Vec2
	}
)

func (s UnknownError) Type() seg.Type {
	return seg.UNKNOWN
}
func (s UnknownError) Error() string {
	if len(s.From) > 0{
		return fmt.Sprintf("Unknown(%s, From : %s)", s.Err.Error(), s.From)
	}
	return fmt.Sprintf("Unknown(%s)", s.Err.Error())
}
func (s UnknownError) String() string {
	if len(s.From) > 0{
		return fmt.Sprintf("UnknownError(%s, From : %s)", s.Err.Error(), s.From)
	}
	return fmt.Sprintf("UnknownError(%s)", s.Err.Error())
}

func (s UnknownCommand) Type() seg.Type {
	return seg.UNKNOWN
}
func (s UnknownCommand) Error() string {
	return fmt.Sprintf("Unknown(%s)", s.Command)
}
func (s UnknownCommand) String() string {
	return fmt.Sprintf("UnknownCommand(%s)", s.Command)
}

func (s ClosePath) Type() seg.Type {
	return seg.CLOSEPATH
}
func (s ClosePath) String() string {
	return fmt.Sprintf("ClosePath()")
}

func (s MoveToAbs) Type() seg.Type {
	return seg.MOVETO_ABS
}
func (s MoveToAbs) String() string {
	return fmt.Sprintf("MoveToAbs((%f, %f))", s.To[0], s.To[1])
}

func (s MoveToRel) Type() seg.Type {
	return seg.MOVETO_REL
}
func (s MoveToRel) String() string {
	return fmt.Sprintf("MoveToRel((%f, %f))", s.To[0], s.To[1])
}

func (s LineToAbs) Type() seg.Type {
	return seg.LINETO_ABS
}
func (s LineToAbs) String() string {
	return fmt.Sprintf("LineToAbs((%f, %f))", s.To[0], s.To[1])
}

func (s LineToRel) Type() seg.Type {
	return seg.LINETO_REL
}
func (s LineToRel) String() string {
	return fmt.Sprintf("LineToRel((%f, %f))", s.To[0], s.To[1])
}

func (s CurveToCubicAbs) Type() seg.Type {
	return seg.CURVETO_CUBIC_ABS
}
func (s CurveToCubicAbs) String() string {
	return fmt.Sprintf("CurveToCubicAbs((%f, %f), (%f, %f), (%f, %f))", s.P0[0], s.P0[1], s.P1[0], s.P1[1], s.To[0], s.To[1])
}

func (s CurveToCubicRel) Type() seg.Type {
	return seg.CURVETO_CUBIC_REL
}
func (s CurveToCubicRel) String() string {
	return fmt.Sprintf("CurveToCubicRel((%f, %f), (%f, %f), (%f, %f))", s.P0[0], s.P0[1], s.P1[0], s.P1[1], s.To[0], s.To[1])
}

func (s CurveToQuadraticAbs) Type() seg.Type {
	return seg.CURVETO_QUADRATIC_ABS
}
func (s CurveToQuadraticAbs) String() string {
	return fmt.Sprintf("CurveToQuadraticAbs((%f, %f), (%f, %f))", s.P0[0], s.P0[1], s.To[0], s.To[1])
}

func (s CurveToQuadraticRel) Type() seg.Type {
	return seg.CURVETO_QUADRATIC_REL
}
func (s CurveToQuadraticRel) String() string {
	return fmt.Sprintf("CurveToQuadraticRel((%f, %f), (%f, %f))", s.P0[0], s.P0[1], s.To[0], s.To[1])
}

func (s ArcAbs) Type() seg.Type {
	return seg.ARC_ABS
}
func (s ArcAbs) String() string {
	return fmt.Sprintf("ArcAbs(To : (%f, %f), Radius : (%f, %f), Angle : %f, LargeArc : %t, Sweep : %t)", s.To[0], s.To[1], s.Radius[0], s.Radius[1], s.Angle, s.LargeArc, s.Sweep)
}

func (s ArcRel) Type() seg.Type {
	return seg.ARC_REL
}
func (s ArcRel) String() string {
	return fmt.Sprintf("ArcRel(To : (%f, %f), Radius : (%f, %f), Angle : %f, LargeArc : %t, Sweep : %t)", s.To[0], s.To[1], s.Radius[0], s.Radius[1], s.Angle, s.LargeArc, s.Sweep)
}

func (s LineToHorizontalAbs) Type() seg.Type {
	return seg.LINETO_HORIZONTAL_ABS
}
func (s LineToHorizontalAbs) String() string {
	return fmt.Sprintf("LineToHorizontalAbs(%f)", s.X)
}

func (s LineToHorizontalRel) Type() seg.Type {
	return seg.LINETO_HORIZONTAL_REL
}
func (s LineToHorizontalRel) String() string {
	return fmt.Sprintf("LineToHorizontalRel(%f)", s.X)
}

func (s LineToVerticalAbs) Type() seg.Type {
	return seg.LINETO_VERTICAL_ABS
}
func (s LineToVerticalAbs) String() string {
	return fmt.Sprintf("LineToVerticalAbs(%f)", s.Y)
}

func (s LineToVerticalRel) Type() seg.Type {
	return seg.LINETO_VERTICAL_REL
}
func (s LineToVerticalRel) String() string {
	return fmt.Sprintf("LineToVerticalRel(%f)", s.Y)
}

func (s CurveToCubicSmoothAbs) Type() seg.Type {
	return seg.CURVETO_CUBIC_SMOOTH_ABS
}
func (s CurveToCubicSmoothAbs) String() string {
	return fmt.Sprintf("CurveToCubicSmoothAbs((%f, %f), (%f, %f))", s.P0[0], s.P0[1], s.To[0], s.To[1])
}

func (s CurveToCubicSmoothRel) Type() seg.Type {
	return seg.CURVETO_CUBIC_SMOOTH_REL
}
func (s CurveToCubicSmoothRel) String() string {
	return fmt.Sprintf("CurveToCubicSmoothRel((%f, %f), (%f, %f))", s.P0[0], s.P0[1], s.To[0], s.To[1])
}

func (s CurveToQuadraticSmoothAbs) Type() seg.Type {
	return seg.CURVETO_QUADRATIC_SMOOTH_ABS
}
func (s CurveToQuadraticSmoothAbs) String() string {
	return fmt.Sprintf("CurveToQuadraticSmoothAbs((%f, %f))", s.To[0], s.To[1])
}
func (s CurveToQuadraticSmoothRel) Type() seg.Type {
	return seg.CURVETO_QUADRATIC_SMOOTH_REL
}
func (s CurveToQuadraticSmoothRel) String() string {
	return fmt.Sprintf("CurveToQuadraticSmoothRel((%f, %f))", s.To[0], s.To[1])
}
