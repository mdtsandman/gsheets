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
	mtx sync.Mutex
}

func NewMeta(ssid, rngStr string) (meta Meta) {
	return Meta{
		ssid: ssid,
		rngStr: rngStr,
	}
}

type Cache struct{
	srv *sheets.Service
	meta map[string]Meta
	data map[string][][]interface{}
	mtx sync.Mutex
}

func NewCache(srv *sheets.Service) (c Cache) {
	c.srv = srv
	c.meta = make(map[string]Meta)
	c.data = make(map[string][][]interface{})
	return c
}

func (c Cache) Register(resource, ssid, rngStr string) (cache Cache, err error) {
	c.mtx.Lock()
	c.meta[resource] = NewMeta(ssid, rngStr)
	c.mtx.Unlock()
	_, err = c.Fetch(resource)
	return c, err
}

func (c Cache) Data(resource string) (data [][]interface{}, err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	_, registered := c.meta[resource]
	if (!registered) {
		return data, errors.New(fmt.Sprintf("Attempt to read unregistered resource [%s].", resource))
	}
	return c.data[resource], nil
}

func (c Cache) Fetch(resource string) (data [][]interface{}, err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	meta, registered := c.meta[resource]
	if (!registered) {
		return data, errors.New(fmt.Sprintf("Attempt to fetch unregistered resource [%s].", resource))
	}
	fmt.Printf("Fetching %s from spreadsheet %s via REST API\n", meta.rngStr, meta.ssid)
	req := c.srv.Spreadsheets.Values.Get(meta.ssid, meta.rngStr)
	resp, err := req.ValueRenderOption("UNFORMATTED_VALUE").DateTimeRenderOption("SERIAL_NUMBER").Do()
	if err != nil {
		return data, err
	}
	c.data[resource] = resp.Values
	return c.data[resource], nil
}

func (c Cache) Update(resource string, data [][]interface{}) (cache Cache, err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	_, registered := c.meta[resource]
	if (!registered) {
		return c, errors.New(fmt.Sprintf("Attempt to update unregistered resource [%s].", resource))
	}
	c.data[resource] = data
	return c, nil
}
