package psvg

import "github.com/go-gl/mathgl/mgl32"

type (
	Renderer struct {
		data  []Elem
		parser *Parser
		support Support
	}
	Support interface {
		MoveTo(to mgl32.Vec2)
		LineTo(to mgl32.Vec2)
		QuadTo(p0, to mgl32.Vec2)
		CubeTo(p0, p1, to mgl32.Vec2)
		CloseTo(to mgl32.Vec2)
	}
)

func NewRenderer(support Support, data []Elem) *Renderer {

	return &Renderer{
		support:support,
		data:nil,
		parser:nil,
	}
}

func NewRendererByParser(support Support, parser *Parser) *Renderer {
	return &Renderer{
		support:support,
		data:nil,
		parser:parser,
	}
}

func NewRendererByDataAndParser(support Support, data []Elem, parser *Parser) *Renderer {
	return &Renderer{
		support:support,
		data:data,
		parser:parser,
	}
}
