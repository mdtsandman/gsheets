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

type lines [][]interface{}
func (list lines) Len() int { return len(list) }
func (list lines) Swap(i,j int) { list[i], list[j] = list[j], list[i] }
func (list lines) Less(i,j int) bool {
	switch list[i][0].(type) {
	case time.Time:
		t1, _ := list[i][0].(time.Time)
		t2, _ := list[j][0].(time.Time)
		return t1.Before(t2)
	case float64:
		f1, _ := list[i][0].(float64)
		f2, _ := list[j][0].(float64)
		return f1 < f2
	default:
		s1 := fmt.Sprintf("%s",list[i][0])
		s2 := fmt.Sprintf("%s",list[j][0])
		return s1 < s2
	}
}

type Sheet struct {
	srv *sheets.Service
	stale bool
	ssid, sheet, rng string
	cols *Hdr
	rows lines
}

func NewSheet(srv *sheets.Service, ssid, sheet, rng string) (*Sheet) {
	return &Sheet{	srv, true, ssid, sheet, rng, nil, nil }
}

func (s Sheet) Name() string {
	return s.ssid + "-" + s.sheet
}

func (s Sheet) Stale() bool {
	return s.stale
}

func (s Sheet) Header() (*Hdr) {
	return s.cols
}

func (s Sheet) AllRows() (rows [][]interface{}) {
	return s.rows
}

func (s Sheet) Rows(startTag, endTag interface{}) (rows [][]interface{}, found bool) {

	searchFxn := func(data lines, tag interface{}) (func(int) bool) {
		switch tag.(type) {
		case time.Time:
			target, _ := tag.(time.Time)
			return func(i int) bool {
				if len(data[i]) == 0 {
					return false
				}
				serial, _ := data[i][0].(float64)
				element, _ := DateTimeFromSerial(serial)
				return !element.Before(target)
			}
		case float64:
			target, _ := tag.(float64)
			return func(i int) bool {
				if len(data[i]) == 0 {
					return false
				}
				element, _ := data[i][0].(float64)
				return element >= target
			}
		default:
			target := fmt.Sprintf("%s",tag)
			return func(i int) bool {
				if len(data[i]) == 0 {
					return false
				}
				element := fmt.Sprintf("%s",data[i][0])
				return element >= target
			}
		}
	}

	present := func(data lines, i int, tag interface{}) (bool) {
		switch tag.(type) {
		case time.Time:
			target, _ := tag.(time.Time)
			if len(data[i]) == 0 {
				return false
			}
			serial, _ := data[i][0].(float64)
			element, _ := DateTimeFromSerial(serial)
			return element.Equal(target)
		case float64:
			target, _ := tag.(float64)
			if len(data[i]) == 0 {
				return false
			}
			element, _ := data[i][0].(float64)
			return element == target
		default:
			target := fmt.Sprintf("%s",tag)
			if len(data[i]) == 0 {
				return false
			}
			element := fmt.Sprintf("%s",data[i][0])
			return element == target
		}
	}

	equal := func(s,e interface{}) bool {
		switch s.(type) {
		case time.Time:
			a, _ := s.(time.Time)
			b, ok := e.(time.Time)
			if !ok {
				serial, ok := e.(float64)
				if !ok {
					return false
				}
				b, _ = DateTimeFromSerial(serial)
			}
			return a.Equal(b)
		case float64:
			a, _ := s.(float64)
			if b, ok := e.(float64); ok {
				return a == b
			}
			return false
		default:
			a := fmt.Sprintf("%s",s)
			b := fmt.Sprintf("%s",e)
			return a == b
		}
	}

	end := len(s.rows)
	first := sort.Search(end, searchFxn(s.rows, startTag))
	last := sort.Search(end, searchFxn(s.rows, endTag))

	switch {
	case first == end || (equal(startTag,endTag) && !present(s.rows,first,startTag)):
		return rows, false
	case present(s.rows,last,endTag):
		return s.rows[first:last+1], true
	default:
		return s.rows[first:last], true
	}

}

func (s *Sheet) SetStale(stale bool) {
	s.stale = stale;
}

func (s *Sheet) Refresh() (e error) {
	rngStr := s.sheet + "!" + s.rng
	fmt.Printf("Updating %s-%s in cache\n", s.ssid, rngStr)
	req := s.srv.Spreadsheets.Values.Get(s.ssid, rngStr)
	resp, err := req.ValueRenderOption("UNFORMATTED_VALUE").DateTimeRenderOption("SERIAL_NUMBER").Do()
	if err != nil {
		return err
	}
	if len(resp.Values) == 0 {
		return errors.New(fmt.Sprintf("Unable to refresh sheet [%s]: No header row found", s.sheet))
	}
	s.cols = NewHdr(resp.Values[0])
	if len(resp.Values) > 1 {
		s.rows = lines(resp.Values[1:])
	} else {
		s.rows = nil
	}
	sort.Sort(s.rows)
	s.stale = false
	return nil
}
