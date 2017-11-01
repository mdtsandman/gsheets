package gsheets

import (
	"fmt"
	"sync"
	"errors"
)

import (
	"google.golang.org/api/sheets/v4"
)

type Resource struct {
	data [][]interface{}
	stale bool
	ssid, sheet, rng string
}

func NewResource(ssid, rng string) (r *Resource) {
	r = new(Resource)
	r.ssid = ssid
	r.rng = rng
	r.stale = true
	return r
}

type Cache struct {
	resources map[string](*Resource)
	srv *sheets.Service
	mtx sync.Mutex
}

func NewCache(srv *sheets.Service) (c Cache) {
	c.resources = make(map[string](*Resource))
	c.srv = srv
	return c
}

func (c Cache) Register(resource, ssid, rng string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	_, registered := c.resources[resource]
	if !registered {
		c.resources[resource] = NewResource(ssid, rng)
	}
	return
}

func (c Cache) fetch(r *Resource) (err error) {
	req := c.srv.Spreadsheets.Values.Get(r.ssid, r.rng)
	resp, err := req.ValueRenderOption("UNFORMATTED_VALUE").DateTimeRenderOption("SERIAL_NUMBER").Do()
	if err != nil {
		return err
	}
	r.data = resp.Values
	r.stale = false
	return nil
}

func (c Cache) Data(resource string) (data [][]interface{}, err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	r, registered := c.resources[resource]
	if !registered {
		return data, errors.New(fmt.Sprintf("Cache: Attempt to access unregistered resource [%s]", resource))
	}
	if !r.stale {
		return r.data, nil
	}
	if err := c.fetch(r); err != nil  {
		return data, err
	}
	return r.data, nil
}

func (c Cache) SetStaleFlag(resource string) (err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	r, registered := c.resources[resource]
	if !registered {
		return errors.New(fmt.Sprintf("Cache: Attempt to set stale flag on unregistered resource [%s]", resource))
	}
	r.stale = true
	return
}

func (c Cache) Update(resource string) (err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	r, registered := c.resources[resource]
	if !registered {
		return errors.New(fmt.Sprintf("Cache: Attempt to update unregistered resource [%s]", resource))
	}
	if err := c.fetch(r); err != nil {
		return err
	}
	return nil
}

func (c Cache) UpdateAll() (err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for _, r := range c.resources {
		if err := c.fetch(r); err != nil {
			return err
		}
	}
	return nil
}
