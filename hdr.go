package gsheets

import (
	"fmt"
	"time"
	"sync"
)

type Hdr struct {
	r []string
	m map[string]int
	mtx sync.Mutex
}

func NewHdr(row []interface{}) (h *Hdr) {
	h = new(Hdr)
	h.m = make(map[string]int)
	for col, rawTag := range row {
		strTag := fmt.Sprint(rawTag)
		h.r = append(h.r, strTag)
		h.m[strTag] = col
	}
	return h
}


// Returns a slice containing the header keys, in order
func (h Hdr) Values() ([]string) {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	return h.r
}


// The following all attempt to convert to a particular type. In case
// of error, a default value is returned and the boolean flag is set
// to false.
func (h Hdr) ChkNum(row []interface{}, key string) (float64, bool) {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	if col, found := h.m[key]; found {
		return ChkNum(row,col)
	}
	return 0, false
}

func (h Hdr) ChkStr(row []interface{}, key string) (string, bool) {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	if col, found := h.m[key]; found {
		return ChkStr(row,col)
	}
	return "", false
}

func (h Hdr) ChkDateTime(row []interface{}, key string) (time.Time, bool) {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	if col, found := h.m[key]; found {
		if serial, ok := ChkNum(row, col); ok {
			return DateTimeFromSerial(serial)
		}
	}
	return SerialTimeZero(), false
}

func (h Hdr) ChkTime(row []interface{}, key string) (time.Duration, bool) {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	if col, found := h.m[key]; found {
		if serial, ok := ChkNum(row, col); ok {
			return TimeFromSerial(serial)
		}
	}
	var dur time.Duration
	return dur, false
}


// The following all return a default value in case of error
func (h Hdr) Num(row []interface{}, key string) (float64) {
	f64, _ := h.ChkNum(row, key)
	return f64
}

func (h Hdr) Str(row []interface{}, key string) (string) {
	str, _ := h.ChkStr(row, key)
	return str
}

func (h Hdr) DateTime(row []interface{}, key string) (time.Time) {
	dt, _ := h.ChkDateTime(row, key)
	return dt
}

func (h Hdr) Time(row []interface{}, key string) (time.Duration) {
	t, _ := h.ChkTime(row, key)
	return t
}
