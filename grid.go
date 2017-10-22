package gsheets

import (
	"sync"
)

type Grid struct {
	colHdr, rowHdr map[string]int
	data [][]interface{}
	mtx sync.Mutex
}

func NewGrid(data [][]interface{}) (g Grid) {
	g.data = data
	g.colHdr = make(map[string]int)
	g.rowHdr = make(map[string]int)
	for i, line := range data {
		if key, ok := line[0].(string); ok && i > 0 {
			g.rowHdr[key] = i
		}
	}
	for j, hdr := range data[0] {
	 if key, ok := hdr.(string); ok && j > 0 {
			g.colHdr[key] = j
		}
	}
	return g
}

func (g Grid) value(rowHdrStr, colHdrStr string) (value interface{}, found bool) {
	if i, present := g.rowHdr[rowHdrStr]; present {
		if j, present := g.colHdr[colHdrStr]; present {
			if j < len(g.data[i]) {
				return g.data[i][j], true
			}
		}
	}
	return value, false
}

func (g Grid) Interface(rowHdrStr, colHdrStr string) (value interface{}, found bool) {
	return g.value(rowHdrStr, colHdrStr)
}

func (g Grid) String(rowHdrStr, colHdrStr string) (value string, found bool) {
	if tmp, ok := g.value(rowHdrStr, colHdrStr); ok {
		if value, ok := tmp.(string); ok {
			return value, true
		}
	}
	return "", false
}

func (g Grid) Float64(rowHdrStr, colHdrStr string) (value float64, found bool) {
	if tmp, found := g.value(rowHdrStr, colHdrStr); found {
		if value, found := tmp.(float64); found {
			return value, true
		}
	}
	return float64(0), false
}

func (g Grid) RowTags() ([]string) {
	var tags []string
	for _, row := range g.data {
		if str, ok := row[0].(string); ok {
			tags = append(tags,str)
		}
	}
	return tags
}

func (g Grid) ColTags() ([]string) {
	var tags []string
	for _, row := range g.data[0] {
		if str, ok := row.(string); ok {
			tags = append(tags,str)
		}
	}
	return tags
}
