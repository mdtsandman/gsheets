package gsheets

import (
	"time"
)

type Table struct {
	colHdr map[string]int
	data [][]interface{}
}

func NewTable(data [][]interface{}) (tbl Table) {
	tbl.data = data
	tbl.colHdr = make(map[string]int)
	if len(tbl.data) == 0 {
		return tbl
	}
	for col, hdr := range (tbl.data)[0] {
		if key, ok := hdr.(string); ok {
			tbl.colHdr[key] = col
		}
	}
	return tbl
}

func (tbl Table) Value(row int, colHdrStr string) (value interface{}) {
	if row < 0 || row >= len(tbl.data) {
		return nil
	}
	col, present := tbl.colHdr[colHdrStr]
	if !present || col >= len((tbl.data)[row]) {
		return nil
	}
	return (tbl.data)[row][col]
}

func (tbl Table) String(row int, colHdrStr string) (value string, ok bool) {
	tmp := tbl.Value(row, colHdrStr)
	if tmp != nil {
		if value, ok := tmp.(string); ok {
			return value, true
		}
	}
	return "", false
}

func (tbl Table) Float64(row int, colHdrStr string) (value float64, ok bool) {
	tmp := tbl.Value(row, colHdrStr)
	if tmp != nil {
		if value, ok := tmp.(float64); ok {
			return value, true
		}
	}
	return float64(0), true
}

func (tbl Table) Serial(row int, colHdrStr string, tz *time.Location) (value Serial, ok bool) {
	s, ok := tbl.Float64(row,colHdrStr)
	return NewSerial(s,tz), ok
}
