package gsheets

import (
	"time"
)

// Google Sheets uses a fixed timezone (in Vancouver this is set to PST, not PDT).

// Value returned by Google Sheets API when dateTimeRenderOption = "SERIAL_NUMBER"
// is a float64 representation of the # of days since Dec 31 1899 @ 00:00.

var DSZone, _ = time.LoadLocation("America/Vancouver")

func midnight(t time.Time) (time.Time) {
	return time.Date(t.Year(),t.Month(),t.Day(),0,0,0,0,t.Location())
}

func SerialTimeZero(tz *time.Location) (time.Time) {
	return time.Date(1899,12,30,0,0,0,0,tz)
}

func FromSerialDate(ds interface{}, sheetZone *time.Location) (t time.Time, ok bool) {
	if f64, ok := ds.(float64); ok && f64 >= 0 {
		serial := SerialTimeZero(sheetZone).Add(time.Duration(f64*60*60*24)*time.Second).Round(time.Second)
		return DSAdj(serial, DSZone), true
	}
	return t, false
}

func SetSerialTime(d time.Time, ts interface{}) (t time.Time, ok bool) {
	if f64, ok := ts.(float64); ok && f64 >= 0 && f64 < 1 {
		return midnight(d).Add(time.Duration(f64*24*60*60)*time.Second).Round(time.Second), true
	}
	return t, false
}

func DSAdj(timeInFixedZone time.Time, variableZone *time.Location) (time.Time) {
	h1 := timeInFixedZone.Hour()
	h2 := timeInFixedZone.In(variableZone).Hour()
	adj := time.Duration(h1-h2) * time.Hour
	return timeInFixedZone.Add(adj).In(variableZone)
}

/*
type Serial struct {
	 float64
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
*/
