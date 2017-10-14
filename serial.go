package gsheets

import (
	"time"
)

// Google Sheets uses a fixed timezone (in Vancouver this is set to PST, not PDT).

// Value returned by Google Sheets API when dateTimeRenderOption = "SERIAL_NUMBER"
// is a float64 representation of the # of days since Dec 31 1899 @ 00:00.
type Serial struct {
	value float64
	tz *time.Location
}

func NewSerialFromFloat64(value float64, tz *time.Location) (s Serial) {
	s.value = value
	s.tz = tz
	return s
}

func NewSerial(value interface{}, tz *time.Location) (s Serial, ok bool) {
	if f64, ok := value.(float64); ok {
		return NewSerialFromFloat64(f64, tz), ok
	}
	return s, false
}

func TimeZero(tz *time.Location) (time.Time) {
	return time.Date(1899,12,30,0,0,0,0,tz)
}

func (s Serial) Time() (time.Time) {
	return TimeZero(s.tz).Add(time.Duration(s.value * 24 * 60 * 60) * time.Second)
}

func (s Serial) Add(other Serial) (Serial) {
	tmp, _ := NewSerial(s.value + other.value, s.tz)
	return tmp
}

func (s Serial) IsValidTime() (bool) {
	return s.value >= 0 && s.value < 1
}

// Given a time in a FIXED zone (eg. PST) as the first param and a zone that has a 
// daylight savings adjustment (eg. PDT) as the second, this function returns the
// equivalent time in the zone with the DS adjustment.
func DaylightSavAdj(fixed time.Time, variableTZ *time.Location) (time.Time) {
	h1 := fixed.Hour()
	h2 := fixed.In(variableTZ).Hour()
	adj := time.Duration(h1-h2) * time.Hour
	return fixed.Add(adj).In(variableTZ)
}
