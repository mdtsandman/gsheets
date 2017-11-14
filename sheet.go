package gsheets

import (
	"fmt"
	"time"
	"sort"
	"errors"
)

import (
	"google.golang.org/api/sheets/v4"
)


// Sheet interface

type Sheet interface {
	Name() string
	Stale() bool
	SetStale(state bool)
	Header() *Hdr
	Refresh() error
}


// BaseSheet struct

type BaseSheet struct {
	srv *sheets.Service
	stale bool
	ssid, sheet, rng string
	cols *Hdr
}

func (s BaseSheet) Name() string {
	return s.sheet
}

func (s BaseSheet) Stale() bool {
	return s.stale
}

func (s BaseSheet) SetStale(state bool) {
	s.stale = state;
}

func (s BaseSheet) Header() (*Hdr) {
	return s.cols
}


// StrTagSheet struct

type StrTagSheet struct {
	BaseSheet
	rows map[string][][]interface{}
}

func NewStrTagSheet(srv *sheets.Service, ssid, sheet, rng string) (*StrTagSheet) {
	base := BaseSheet{srv, true, ssid, sheet, rng, nil}
	return &StrTagSheet{base, nil}
}

func (s StrTagSheet) AllRows() (rows [][]interface{}) {
	for _, value := range s.rows {
		rows = append(rows, value...)
	}
	return rows
}

func (s StrTagSheet) FindRows(tag string) (rows [][]interface{}, ok bool) {
	rows, found := s.rows[tag]
	return rows, found
}

func (s StrTagSheet) Refresh() (e error) {
	rngStr := s.sheet + "!" + s.rng
	fmt.Printf("Updating %s in cache\n", rngStr)
	req := s.srv.Spreadsheets.Values.Get(s.ssid, s.sheet + "!" + s.rng)
	resp, err := req.ValueRenderOption("UNFORMATTED_VALUE").DateTimeRenderOption("SERIAL_NUMBER").Do()
	if err != nil {
		return err
	}
	data := resp.Values
	s.cols = NewHdr(data[0])
	s.rows = make(map[string][][]interface{})
	if len(data) > 1 {
		for i, row := range data[1:] {
			rawTag := row[0]
			str, ok := rawTag.(string)
			if !ok {
				return errors.New( fmt.Sprintf(
					"Unable to convert value in row 0 of column %d of sheet %s to a string: %s",
					i,
					s.sheet,
					rawTag,
				) )
			}
			s.rows[str] = append(s.rows[str], row)
		}
	}
	s.stale = false
	return nil
}


// DateTimeTagSheet struct

type DateTimeTagSheet struct {
	BaseSheet
	rows map[float64][][]interface{}
}

func NewDateTimeTagSheet(srv *sheets.Service, ssid, sheet, rng string) (*DateTimeTagSheet) {
	base := BaseSheet{srv, true, ssid, sheet, rng, nil}
	return &DateTimeTagSheet{base, nil}
}

/*
type Int64s []int64
func (i Int64s) Len() int { return len(i) }
func (i Int64s) Swap(x,y int) { i[x], i[y] = i[y], i[x]}
func (i Int64s) Less(x,y int) bool { return i[x] < i[y] }
*/

func (s DateTimeTagSheet) FindRows(start, end time.Time) (rows [][]interface{}, ok bool) {
	sStart := float64(start.Sub(SerialTimeZero()).Nanoseconds()/1000)
	sEnd := float64(end.Sub(SerialTimeZero()).Nanoseconds()/1000)
	var keys sort.Float64Slice
	for key, _ := range s.rows {
		keys = append(keys, key)
	}
	sort.Sort(keys)
	for _, key := range keys {
		if key >= sStart && key <= sEnd {
			rows = append(rows, s.rows[key]...)
		}
	}
	return rows, len(rows) > 0
}

func (s DateTimeTagSheet) Refresh() (e error) {
	rngStr := s.sheet + "!" + s.rng
	fmt.Printf("Updating %s in cache\n", rngStr)
	req := s.srv.Spreadsheets.Values.Get(s.ssid, s.sheet + "!" + s.rng)
	resp, err := req.ValueRenderOption("UNFORMATTED_VALUE").DateTimeRenderOption("SERIAL_NUMBER").Do()
	if err != nil {
		return err
	}
	data := resp.Values
	s.cols = NewHdr(data[0])
	s.rows = make(map[float64][][]interface{})
	if len(data) > 1 {
		for i, row := range data[1:] {
			rawTag := row[0]
			millisecs, ok := rawTag.(float64)
			if !ok {
				return errors.New( fmt.Sprintf(
					"Unable to convert value in row 0 of column %d of sheet %s to a float64 value: %s",
					i,
					s.sheet,
					millisecs,
				) )
			}
			s.rows[millisecs] = append(s.rows[millisecs], row)
		}
	}
	s.stale = false
	return nil
}


// YearSheet

type YearSheet struct {
	BaseSheet
	rows [][]interface{}
}

func NewYearSheet(srv *sheets.Service, ssid, sheet, rng string) (*YearSheet) {
	base := BaseSheet{srv, true, ssid, sheet, rng, nil}
	return &YearSheet{base, nil}
}

func (s YearSheet) Rows(start, end int) (rows [][]interface{}, e error) {
	if start > end || start >= len(rows) || start < 0 {
		return rows, errors.New( fmt.Sprintf(
			"YearSheet.GetRows(start,end): Invalid parameter(s): start=%d, end=%d",
			start,
			end,
		) )
	}
	return rows[start:end], nil
}

func (s YearSheet) Refresh() (e error) {
	rngStr := s.sheet + "!" + s.rng
	fmt.Printf("Updating %s in cache\n", rngStr)
	req := s.srv.Spreadsheets.Values.Get(s.ssid, s.sheet + "!" + s.rng)
	resp, err := req.ValueRenderOption("UNFORMATTED_VALUE").DateTimeRenderOption("SERIAL_NUMBER").Do()
	if err != nil {
		return err
	}
	s.rows = resp.Values
	s.cols = NewHdr(s.rows[0])
	s.stale = false
	return nil
}
