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

func NewSheet(srv *sheets.Service, ssid, sheet, rng string) (Sheet) {
	return Sheet{	srv, true, ssid, sheet, rng, nil, nil }
}

func (s Sheet) Name() string {
	return s.sheet
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
			x, _ := tag.(time.Time)
			return func(i int) bool {
				serial, _ := data[i][0].(float64)
				y, _ := DateTimeFromSerial(serial)
				return !(y.Before(x))
			}
		case float64:
			x, _ := tag.(float64)
			return func(i int) bool {
				y, _ := data[i][0].(float64)
				return y >= x
			}
		default:
			x := fmt.Sprintf("%s",tag)
			return func(i int) bool {
				y := fmt.Sprintf("%s",data[i][0])
				return y >= x
			}
		}
	}

	next := func(tag interface{}) interface{} {
		switch tag.(type) {
		case time.Time:
			x, _ := tag.(time.Time)
			return x.Add(time.Nanosecond)
		case float64:
			x, _ := tag.(float64)
			return x + 1
		default:
			x := fmt.Sprintf("%s",tag)
			return x + "a"
		}
	}

	first := sort.Search(len(s.rows), searchFxn(s.rows, startTag))
	last := sort.Search(len(s.rows), searchFxn(s.rows, next(endTag)))

	result := s.rows[first:last]

	fmt.Printf("StartTag: %s | EndTag: %s | First: %d | Last: %d\n", startTag, next(endTag), first, last)

	return result, len(result) > 0

}

func (s *Sheet) SetStale(state bool) {
	s.stale = state;
}

func (s *Sheet) Refresh() (e error) {
	rngStr := s.sheet + "!" + s.rng
	fmt.Printf("Updating %s in cache\n", rngStr)
	req := s.srv.Spreadsheets.Values.Get(s.ssid, rngStr)
	resp, err := req.ValueRenderOption("UNFORMATTED_VALUE").DateTimeRenderOption("SERIAL_NUMBER").Do()
	if err != nil {
		return err
	}
	if len(resp.Values) == 0 {
		return errors.New(fmt.Sprintf("Unable to refrsh sheet [%s]: No header row found", s.sheet))
	}
	s.rows = resp.Values
	s.cols = NewHdr(s.rows[0])
	sort.Sort(s.rows)
	s.stale = false
	return nil
}
