package psvg

var symbols = []byte{
	// https://www.w3.org/TR/SVG/paths.html#PathDataMovetoCommands
	'M', 'm', // MoveToAbs, MoveToRel

	// https://www.w3.org/TR/SVG/paths.html#PathDataClosePathCommand
	'Z', 'z', // ClosePath

	// https://www.w3.org/TR/SVG/paths.html#PathDataLinetoCommands
	'L', 'l', // LineToAbs, LineToRel
	'H', 'h', // LineToHorizontalAbs, LineToHorizontalRel
	'V', 'v', // LineToVerticalAbs, LineToVerticalRel

	// https://www.w3.org/TR/SVG/paths.html#PathDataCurveCommands
	'C', 'c', // CurveToCubicAbs, CurveToCubicRel]
	'S', 's', // CurveToCubicSmoothAbs, CurveToCubicSmoothRel

	// https://www.w3.org/TR/SVG/paths.html#PathDataCubicBezierCommands
	'Q', 'q', // CurveToQuadraticAbs, CurveToQuadraticRel
	'T', 't', // CurveToQuadraticSmoothAbs, CurveToQuadraticSmoothRel

	// https://www.w3.org/TR/SVG/paths.html#PathDataEllipticalArcCommands
	'A', 'a', // ArcAbs, ArcRel
}

func matchingSymbol(b byte) bool {
	for _, sym := range symbols {
		if b == sym{
			return true
		}
	}
	return false
}
