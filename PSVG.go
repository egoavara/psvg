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
	buf []byte
	// previous remain bytes
	cmd byte
	prv []byte
	// buffer size
	bsz int
	// buffer offset
	off int
}

func NewParser(src io.Reader) *Parser {
	return &Parser{
		src:      src,
		temporal: list.New(),
		buf:      make([]byte, bufferSize),
		prv:      nil,
		bsz:      0,
		off:      0,
	}
}
// Elem can be error interface,
// Unknown* is Elem, also error
//
// If read all Elems from 'src',
// It return nil
// Return nil mean, 'src' faces io.EOF
func (s *Parser) Next() Elem {
	// If remain previous parsed Elem
	if s.temporal.Len() > 0 {
		// return that
		temp := s.temporal.Front()
		s.temporal.Remove(temp)
		return temp.Value.(Elem)
	}
	// Else there is no Elems remain
	var err error
	// If there is no Remain buffer
	if s.bsz <= s.off {
		// Read buffer from src
		s.bsz, err = s.src.Read(s.buf)
		if err == io.EOF {
			if len(s.prv) > 0{
				elem := convert(s.prv[0], s.prv[1:])
				if len(elem) > 1 {
					// If there is remain, Push it all to temporal
					for _, e := range elem[1:] {
						s.temporal.PushBack(e)
					}
				}
				s.cmd = 0
				s.prv = nil
				// return first elem
				return elem[0]
			}
			// Read all Elem
			return nil
		}
		if err != nil {
			return UnknownError{Err: err}
		}
		s.off = 0
	}
	var to = s.off
	if s.cmd == 0{
		s.cmd = s.buf[to]
		to += 1
	}
	for ; to < s.bsz; to++ {
		if matchingSymbol(s.buf[to]) {
			// parsing bytes
			temp := append(s.prv, s.buf[s.off:to]...)
			s.prv = nil
			elem := convert(temp[0], temp[1:])
			if len(elem) > 1 {
				// If there is remain, Push it all to temporal
				for _, e := range elem[1:] {
					s.temporal.PushBack(e)
				}
			}
			// setup offset for next
			s.cmd = 0
			s.off = to
			// return first elem
			return elem[0]
		}
	}
	// Out of remaining buffer data
	if to >= s.bsz {
		s.prv = append(s.prv, s.buf[s.off:s.bsz]...)
		s.off = s.bsz
		return s.Next()
	}
	return UnknownError{Err: errors.New("Unexpected Parsing Fail")}
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
				P0: vecs[i+0],
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
				P0: vecs[i+0],
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
		res = append(res, UnknownCommand{Command:string(command) + string(data)})
	}
	return
}