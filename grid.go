package gsheets

import (
	"time"
)

type Grid struct {
	colHdr, rowHdr map[string]int
	data [][]interface{}
}

func NewGrid(data [][]interface{}) (g Grid, ok bool) {
	if data == nil || len(data) == 0 || len(data[0]) == 0 {
		return g, false
	}
	g.data = data
	g.colHdr = make(map[string]int)
	g.rowHdr = make(map[string]int)
	for i, line := range data {
		if key, ok := line[0].(string); ok {
			g.rowHdr[key] = i
		}
	}
	for j, hdr := range data[0] {
	 if key, ok := hdr.(string); ok {
			g.colHdr[key] = j
		}
	}
	return g, true
}

func (g Grid) Value(rowHdrStr, colHdrStr string) (value interface{}, ok bool) {
	if g.data == nil {
		return value, false
	}
	if i, present := g.rowHdr[rowHdrStr]; present {
		if j, present := g.colHdr[colHdrStr]; present {
			if j < len(g.data[i]) {
				return g.data[i][j], true
			}
		}
	}
	return value, false
}

func (g Grid) String(rowHdrStr, colHdrStr string) (value string, ok bool) {
	if tmp, ok := g.Value(rowHdrStr, colHdrStr); ok {
		if value, ok := tmp.(string); ok {
			return value, true
		}
	}
	return "", false
}

func (g Grid) Float64(rowHdrStr, colHdrStr string) (value float64, ok bool) {
	if tmp, ok := g.Value(rowHdrStr, colHdrStr); ok {
		if value, ok := tmp.(float64); ok {
			return value, true
		}
	}
	return float64(0), true
}

func (g Grid) Serial(rowHdrStr, colHdrStr string, tz *time.Location) (value Serial, ok bool) {
	s, ok := g.Float64(rowHdrStr,colHdrStr)
	return NewSerial(s, tz), ok
}
