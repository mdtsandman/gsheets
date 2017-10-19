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

func (c Cache) SetChangedFlag(resource string) (cache Cache) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if meta, exists := c.meta[resource]; exists {
		meta.ok = false
		c.meta[resource] = meta
	}
	return cache
}

func (c Cache) SetMeta(resource string, meta Meta) (cache Cache) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.meta[resource] = meta
	return cache
}

func (c Cache) Fetch(srv *sheets.Service, resource string) (data [][]interface{}, err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	meta := c.meta[resource]
	if meta.ok {
		fmt.Printf("Using cached data for %s\n", meta.rngStr)
		return c.data[resource], nil
	}
	fmt.Printf("No cached data found for %s\nFetching from spreadsheet via REST API\n", meta.rngStr)
	req := srv.Spreadsheets.Values.Get(meta.ssid, meta.rngStr)
	resp, err := req.ValueRenderOption("UNFORMATTED_VALUE").DateTimeRenderOption("SERIAL_NUMBER").Do()
	if err != nil {
		return data, err
	}
	c.data[resource] = resp.Values
	meta.ok = true
	c.meta[resource] = meta
	return c.data[resource], nil
}
