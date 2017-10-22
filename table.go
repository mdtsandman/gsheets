package gsheets

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

func (tbl Table) NumRows() (int) {
	return len(tbl.data)
}

func (tbl Table) ColTags() ([]string) {
	var tags []string
	for _, tag := range tbl.data[0] {
		if str, ok := tag.(string); ok {
			tags = append(tags,str)
		}
	}
	return tags
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

func (tbl Table) AddRows(rows [][]interface{}) (Table) {
	tbl.data = append(tbl.data, rows...)
	return tbl
}
