package gsheets

import (
	"fmt"
	"sync"
)

import (
	"google.golang.org/api/sheets/v4"
)

type Meta struct {
	ssid, rngStr string
	ok bool
	mtx sync.Mutex
}

func NewMeta(ssid, rngStr string) (meta Meta) {
	return Meta{
		ssid: ssid,
		rngStr: rngStr,
		ok: false,
	}
}

type Cache struct{
	meta map[string]Meta
	data map[string][][]interface{}
	mtx sync.Mutex
}

func NewCache() (cache Cache) {
	cache.meta = make(map[string]Meta)
	cache.data = make(map[string][][]interface{})
	return cache
}

func (sc Cache) SetMeta(resource string, meta Meta) {
	sc.mtx.Lock()
	defer sc.mtx.Unlock()
	sc.meta[resource] = meta
}

func (sc Cache) Fetch(srv *sheets.Service, resource string) (data [][]interface{}, err error) {
	sc.mtx.Lock()
	defer sc.mtx.Unlock()
	meta := sc.meta[resource]
	if meta.ok {
		fmt.Printf("Using cached data for %s\n", meta.rngStr)
		return sc.data[resource], nil
	}
	fmt.Printf("No cached data found for %s\nFetching from spreadsheet via REST API\n", meta.rngStr)
	req := srv.Spreadsheets.Values.Get(meta.ssid, meta.rngStr)
	resp, err := req.ValueRenderOption("UNFORMATTED_VALUE").DateTimeRenderOption("SERIAL_NUMBER").Do()
	if err != nil {
		return data, err
	}
	sc.data[resource] = resp.Values
	meta.ok = true
	sc.meta[resource] = meta
	return sc.data[resource], nil
}
