package psvg

import (
	"bytes"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/pkg/errors"
	"strconv"
	"unicode"
)

type ArcArguments struct {
	To       mgl32.Vec2
	Radius   mgl32.Vec2
	Angle    float32
	LargeArc bool
	Sweep    bool
}

func vectors(bts []byte) (res []mgl32.Vec2, err error) {
	bts = bytes.TrimSpace(bts)
	var x *float32
	var from = 0
	for to, b := range bts {
		var temp float64
		switch b {
		case ' ':
			fallthrough
		case ',':
			temp, err = strconv.ParseFloat(string(bts[from:to]), 32)
			if err != nil {
				return nil, err
			}
			if x == nil {
				temp32 := float32(temp)
				x = &temp32
			} else {
				res = append(res, mgl32.Vec2{*x, float32(temp)})
				x = nil
			}
			from = to + 1
		}
	}
	temp, err := strconv.ParseFloat(string(bts[from:]), 32)
	if err != nil {
		return nil, err
	}
	if x == nil {
		temp32 := float32(temp)
		x = &temp32
	} else {
		res = append(res, mgl32.Vec2{*x, float32(temp)})
		x = nil
	}
	if x != nil {
		return nil, errors.New("There is remain float")
	}
	return
}
func floats(bts []byte) (res []float32, err error) {
	bts = bytes.TrimSpace(bts)
	var from = 0
	for to, b := range bts {
		var temp float64
		if unicode.IsSpace(rune(b)) || b == ',' {
			temp, err = strconv.ParseFloat(string(bts[from:to]), 32)
			if err != nil {
				return nil, err
			}
			res = append(res, float32(temp))
			from = to + 1
		}
	}
	temp, err := strconv.ParseFloat(string(bts[from:]), 32)
	if err != nil {
		return nil, err
	}
	res = append(res, float32(temp))
	return
}
func arcArgs(bts []byte) (res []ArcArguments, err error) {
	temp, err := floats(bts)
	if err != nil {
		return nil, err
	}
	if len(temp)%7 != 0 {
		return nil, errors.New("each arc argument have 7 arg(float, float, degree, flag, flag, float, float)")
	}
	res = make([]ArcArguments, len(temp)/7)
	for i, v := range temp {
		switch i % 7 {
		case 0:
			res[i/7].Radius[0] = v
		case 1:
			res[i/7].Radius[1] = v
		case 2:
			res[i/7].Angle = v
		case 3:
			res[i/7].LargeArc = v == 1.
		case 4:
			res[i/7].Sweep = v == 1.
		case 5:
			res[i/7].To[0] = v
		case 6:
			res[i/7].To[1] = v
		}
	}
	return res, nil
}
