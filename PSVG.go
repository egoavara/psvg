// Spec SVG1.1
// https://www.w3.org/TR/SVG/paths.html
package psvg

import (
	"container/list"
	"github.com/pkg/errors"
	"io"
)

const bufferSize = 1024

type Parser struct {
	src io.Reader
	// Store for command
	temporal *list.List
	// buffer
	cmd byte
	buf []byte
	prv []byte
}

func NewParser(src io.Reader) *Parser {
	return &Parser{
		src:      src,
		cmd:      0,
		temporal: list.New(),
	}
}

// Elem can be error interface,
// Unknown* is Elem, also error
//
// If read all Elems from 'src',
// It return nil
// Return nil mean, 'src' faces io.EOF
func (s *Parser) Next() Elem {
	// if there is remain Elem
	if s.temporal.Len() > 0 {
		res := s.temporal.Front().Value
		if res == nil {
			return nil
		}
		s.temporal.Remove(s.temporal.Front())
		return res.(Elem)
	}
	// read data from src
	if len(s.buf) == 0 {
		s.buf = make([]byte, bufferSize)
		n, err := s.src.Read(s.buf)
		if err != nil {
			if err == io.EOF {
				for _, e := range convert(s.cmd, s.prv) {
					s.temporal.PushBack(e)
				}
				s.temporal.PushBack(nil)
				return s.Next()
			}
			return UnknownError{Err: err, From: string(s.buf[:n])}
		}
		s.buf = s.buf[:n]
	}
	//
	for i, v := range s.buf {
		if matchingSymbol(v) {
			if s.cmd == 0 {
				s.cmd = v
				s.buf = s.buf[i+1:]
				return s.Next()
			}
			data := append(s.prv, s.buf[:i]...)
			s.prv = nil
			for _, e := range convert(s.cmd, data) {
				s.temporal.PushBack(e)
			}
			s.cmd = v
			s.buf = s.buf[i+1:]
			return s.Next()
		}
	}
	s.prv = append(s.prv, s.buf...)
	s.buf = nil
	return s.Next()

}

// Allways return at least 1 args
func convert(command byte, data []byte) (res []Elem) {
	switch command {
	case 'Z':
		fallthrough
	case 'z':
		res = []Elem{ClosePath{}}

	case 'M':
		vecs, err := vectors(data)
		if err != nil {
			return []Elem{UnknownError{Err: err, From: string(command) + string(data)}}
		}
		res = make([]Elem, len(vecs))
		for i, vec := range vecs {
			res[i] = MoveToAbs{
				To: vec,
			}
		}
	case 'm':
		vecs, err := vectors(data)
		if err != nil {
			return []Elem{UnknownError{Err: err, From: string(command) + string(data)}}
		}
		res = make([]Elem, len(vecs))
		for i, vec := range vecs {
			res[i] = MoveToRel{
				To: vec,
			}
		}

	case 'L':
		vecs, err := vectors(data)
		if err != nil {
			return []Elem{UnknownError{Err: err, From: string(command) + string(data)}}
		}
		res = make([]Elem, len(vecs))
		for i, vec := range vecs {
			res[i] = LineToAbs{
				To: vec,
			}
		}
	case 'l':
		vecs, err := vectors(data)
		if err != nil {
			return []Elem{UnknownError{Err: err, From: string(command) + string(data)}}
		}
		res = make([]Elem, len(vecs))
		for i, vec := range vecs {
			res[i] = LineToRel{
				To: vec,
			}
		}
	case 'H':
		f32s, err := floats(data)
		if err != nil {
			return []Elem{UnknownError{Err: err, From: string(command) + string(data)}}
		}
		res = make([]Elem, len(f32s))
		for i, f32 := range f32s {
			res[i] = LineToHorizontalAbs{
				X: f32,
			}
		}
	case 'h':
		f32s, err := floats(data)
		if err != nil {
			return []Elem{UnknownError{Err: err, From: string(command) + string(data)}}
		}
		res = make([]Elem, len(f32s))
		for i, f32 := range f32s {
			res[i] = LineToHorizontalRel{
				X: f32,
			}
		}
	case 'V':
		f32s, err := floats(data)
		if err != nil {
			return []Elem{UnknownError{Err: err, From: string(command) + string(data)}}
		}
		res = make([]Elem, len(f32s))
		for i, f32 := range f32s {
			res[i] = LineToVerticalAbs{
				Y: f32,
			}
		}
	case 'v':
		f32s, err := floats(data)
		if err != nil {
			return []Elem{UnknownError{Err: err, From: string(command) + string(data)}}
		}
		res = make([]Elem, len(f32s))
		for i, f32 := range f32s {
			res[i] = LineToVerticalRel{
				Y: f32,
			}
		}

	case 'C':
		vecs, err := vectors(data)
		if err != nil {
			return []Elem{UnknownError{Err: err, From: string(command) + string(data)}}
		}
		if len(vecs)%3 != 0 {
			return []Elem{UnknownError{Err: errors.New("CurveToCubicAbs has a multiple of 3 vectors"), From: string(command) + string(data)}}
		}
		res = make([]Elem, len(vecs)/3)
		for i := 0; i < len(vecs); i += 3 {
			res[i/3] = CurveToCubicAbs{
				P0: vecs[i+0],
				P1: vecs[i+1],
				To: vecs[i+2],
			}
		}
	case 'c':
		vecs, err := vectors(data)
		if err != nil {
			return []Elem{UnknownError{Err: err, From: string(command) + string(data)}}
		}
		if len(vecs)%3 != 0 {
			return []Elem{UnknownError{Err: errors.New("CurveToCubicRel has a multiple of 3 vectors"), From: string(command) + string(data)}}
		}
		res = make([]Elem, len(vecs)/3)
		for i := 0; i < len(vecs); i += 3 {
			res[i/3] = CurveToCubicRel{
				P0: vecs[i+0],
				P1: vecs[i+1],
				To: vecs[i+2],
			}
		}
	case 'S':
		vecs, err := vectors(data)
		if err != nil {
			return []Elem{UnknownError{Err: err, From: string(command) + string(data)}}
		}
		if len(vecs)%2 != 0 {
			return []Elem{UnknownError{Err: errors.New("CurveToCubicSmoothAbs has a multiple of 2 vectors"), From: string(command) + string(data)}}
		}
		res = make([]Elem, len(vecs)/2)
		for i := 0; i < len(vecs); i += 2 {
			res[i/2] = CurveToCubicSmoothAbs{
				P1: vecs[i+0],
				To: vecs[i+1],
			}
		}
	case 's':
		vecs, err := vectors(data)
		if err != nil {
			return []Elem{UnknownError{Err: err, From: string(command) + string(data)}}
		}
		if len(vecs)%2 != 0 {
			return []Elem{UnknownError{Err: errors.New("CurveToCubicSmoothRel has a multiple of 2 vectors"), From: string(command) + string(data)}}
		}
		res = make([]Elem, len(vecs)/2)
		for i := 0; i < len(vecs); i += 2 {
			res[i/2] = CurveToCubicSmoothRel{
				P1: vecs[i+0],
				To: vecs[i+1],
			}
		}

	case 'Q':
		vecs, err := vectors(data)
		if err != nil {
			return []Elem{UnknownError{Err: err, From: string(command) + string(data)}}
		}
		if len(vecs)%2 != 0 {
			return []Elem{UnknownError{Err: errors.New("CurveToQuadraticAbs has a multiple of 2 vectors"), From: string(command) + string(data)}}
		}
		res = make([]Elem, len(vecs)/2)
		for i := 0; i < len(vecs); i += 2 {
			res[i/2] = CurveToQuadraticAbs{
				P0: vecs[i+0],
				To: vecs[i+1],
			}
		}
	case 'q':
		vecs, err := vectors(data)
		if err != nil {
			return []Elem{UnknownError{Err: err, From: string(command) + string(data)}}
		}
		if len(vecs)%2 != 0 {
			return []Elem{UnknownError{Err: errors.New("CurveToQuadraticRel has a multiple of 2 vectors"), From: string(command) + string(data)}}
		}
		res = make([]Elem, len(vecs)/2)
		for i := 0; i < len(vecs); i += 2 {
			res[i/2] = CurveToQuadraticRel{
				P0: vecs[i+0],
				To: vecs[i+1],
			}
		}
	case 'T':
		vecs, err := vectors(data)
		if err != nil {
			return []Elem{UnknownError{Err: err, From: string(command) + string(data)}}
		}
		res = make([]Elem, len(vecs))
		for i, vec := range vecs {
			res[i] = CurveToQuadraticSmoothAbs{
				To: vec,
			}
		}
	case 't':
		vecs, err := vectors(data)
		if err != nil {
			return []Elem{UnknownError{Err: err, From: string(command) + string(data)}}
		}
		res = make([]Elem, len(vecs))
		for i, vec := range vecs {
			res[i] = CurveToQuadraticSmoothRel{
				To: vec,
			}
		}

	case 'A':
		args, err := arcArgs(data)
		if err != nil {
			return []Elem{UnknownError{Err: err, From: string(command) + string(data)}}
		}
		res = make([]Elem, len(args))
		for i, arg := range args {
			res[i] = ArcAbs{
				To:       arg.To,
				Radius:   arg.Radius,
				Angle:    arg.Angle,
				LargeArc: arg.LargeArc,
				Sweep:    arg.Sweep,
			}
		}
	case 'a':
		args, err := arcArgs(data)
		if err != nil {
			return []Elem{UnknownError{Err: err, From: string(command) + string(data)}}
		}
		res = make([]Elem, len(args))
		for i, arg := range args {
			res[i] = ArcRel{
				To:       arg.To,
				Radius:   arg.Radius,
				Angle:    arg.Angle,
				LargeArc: arg.LargeArc,
				Sweep:    arg.Sweep,
			}
		}
	default:
		res = append(res, UnknownCommand{Command: string(command) + string(data)})
	}
	return
}
