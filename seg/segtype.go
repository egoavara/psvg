package seg

type Type uint16

// https://www.w3.org/TR/SVG/paths.html#InterfaceType
const (
	UNKNOWN                      Type = 0
	CLOSEPATH                    Type = 1
	MOVETO_ABS                   Type = 2
	MOVETO_REL                   Type = 3
	LINETO_ABS                   Type = 4
	LINETO_REL                   Type = 5
	CURVETO_CUBIC_ABS            Type = 6
	CURVETO_CUBIC_REL            Type = 7
	CURVETO_QUADRATIC_ABS        Type = 8
	CURVETO_QUADRATIC_REL        Type = 9
	ARC_ABS                      Type = 10
	ARC_REL                      Type = 11
	LINETO_HORIZONTAL_ABS        Type = 12
	LINETO_HORIZONTAL_REL        Type = 13
	LINETO_VERTICAL_ABS          Type = 14
	LINETO_VERTICAL_REL          Type = 15
	CURVETO_CUBIC_SMOOTH_ABS     Type = 16
	CURVETO_CUBIC_SMOOTH_REL     Type = 17
	CURVETO_QUADRATIC_SMOOTH_ABS Type = 18
	CURVETO_QUADRATIC_SMOOTH_REL Type = 19
)
