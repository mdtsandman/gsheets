package gsheets

import (
	"log"
	"time"
)

// Google Sheets uses a fixed timezone (in Vancouver this is set to PST, not PDT).

// Value returned by Google Sheets API when dateTimeRenderOption = "SERIAL_NUMBER"
// is a float64 representation of the # of days since Dec 31 1899 @ 00:00.

const dateFormat = "2006-01-02T15:04:05"

func Zones() (local, utc, fixedZone, dsZone *time.Location) {
	var err error
	local, err = time.LoadLocation("Local")
	if err != nil {
		log.Fatal("Unable to load TZ = Local")
	}
	utc, err = time.LoadLocation("UTC")
	if err != nil {
		log.Fatal("Unable to load TZ = UTC")
	}
	dsZone, err = time.LoadLocation("America/Vancouver")
	if err != nil {
		log.Fatal("Unable to load TZ = America/Vancouver")
	}
	fixedZone = time.FixedZone("PST", -8*60*60)
	return local, utc, fixedZone, dsZone
}

var Local, UTC, FixedZone, DSZone = Zones()

func Midnight(t time.Time) (time.Time) {
	return time.Date(t.Year(),t.Month(),t.Day(),0,0,0,0,t.Location())
}

func Monday(t time.Time) (time.Time) {
	offset := int(t.Weekday()) - 1 // week begins on Monday
	if (offset < 0) {
		offset = 6
	}
	return t.AddDate(0,0,0-offset)
}

func SerialTimeZero() (time.Time) {
	return time.Date(1899,12,30,0,0,0,0,FixedZone)
}

func DateTimeFromSerial(serialDateTime interface{}) (result time.Time, ok bool) {
	if f64, ok := serialDateTime.(float64); ok && f64 >= 0 {
		result := SerialTimeZero().Add(time.Duration(f64*24)*time.Hour)
		return DSAdj(result.Round(time.Minute), DSZone), true
	}
	return result, false
}

func TimeFromSerial(serialTime interface{}) (result time.Duration, ok bool) {
	if f64, ok := serialTime.(float64); ok && f64 >= 0 && f64 < 1 {
		var t time.Time
		result = time.Duration(f64*24)*time.Hour
		result = Midnight(t).Add(result).Round(time.Minute).Sub(Midnight(t)) // time.Duration should really have a Round() method!
		return result, true
	}
	return result, false
}

func DSAdj(timeInFixedZone time.Time, variableZone *time.Location) (time.Time) {
	h1 := timeInFixedZone.Hour()
	h2 := timeInFixedZone.In(variableZone).Hour()
	adj := time.Duration(h1-h2) * time.Hour
	return timeInFixedZone.Add(adj).In(variableZone)
}

func DateTimeFromStr(str string) (result time.Time, ok bool) {
	if len(str) < 19 {
		return result, false
	}
	// Google Sheets formats dates and datetimes in the timezone specified in the spreadsheet locale
	d, err := time.ParseInLocation(dateFormat,str[0:19], FixedZone)
	if err != nil {
		return result, false
	}
	return DSAdj(d.Round(time.Minute), DSZone), true
}

func TimeFromStr(str string) (result time.Duration, ok bool) {
	if len(str) < 19 {
		return result, false
	}
	// Google Sheets formats **standalone times** (not dates or datetimes) as relative to 
	// 1899-12-30T00:00:00Z00:00 (ie. in UTC) regardless of the spreadsheet locale.
	t, err := time.ParseInLocation(dateFormat,str[0:19], UTC)
	if err != nil {
	}
	t = t.In(FixedZone)
	mn := Midnight(t)
	return t.Sub(mn), true
}
