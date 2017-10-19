package gsheets

import (
	"fmt"
	"sync"
	"errors"
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

func (c Cache) Fetch(srv *sheets.Service, resource string, force bool) (data [][]interface{}, err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	meta, exists := c.meta[resource]
	if (!exists) {
		return data, errors.New(fmt.Sprintf("Attempt to fetch resource [%s], which is not registered in the cache.",resource))
	}
	if !force && meta.ok {
		fmt.Printf("Using cached data for resource [%s]\n", meta.rngStr)
		return c.data[resource], nil
	}
	fmt.Printf("Fetching resource [%s] from spreadsheet via REST API\n", meta.rngStr)
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
