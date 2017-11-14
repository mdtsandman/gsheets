package gsheets

// Parameters: row = an array representing one row in a spreadsheet; i = an index in that array.
// Returns: a float64 representing the value of row at index i, or a zeroed float64 if that value can't be represented as a float64 or if i is out of bounds..
func ChkNum(row []interface{}, i int) (float64, bool) {
	if i < 0 || i >= len(row) {
		return 0, false // array index out of bounds
	}
	if f64, ok := row[i].(float64); ok {
		return f64, true
	}
	return 0, false // data is not a float64 value
}

// Parameters: row = an array representing one row in a spreadsheet; i = an index in that array.
// Returns: a string representing the value of row at index i, or an empty string if i is out of bounds.
func ChkStr(row []interface{}, i int) (string, bool) {
	if i < 0 || i >= len(row) {
		return "", false // aray index out of bounds
	}
	if str, ok := row[i].(string); ok {
		return str, true
	}
	return "", false // data is not a string
}
