iackage gsheets

import (
	"fmt"
	"sync"
	"errors"
)

import (
	"google.golang.org/api/sheets/v4"
)

type Sheet interface {
	func Stale() bool
	func 
}

type Table struct {
	stale bool
	ssid, sheet, rng string
	data [][]interface{}
	cols map[string]int
}

type Grid struct {
	Table
	rows map[string]*[]interface{}
}

func (g *Grid) buildMap() {
	g.rows := make(map[string]*[]interface)
	for _, row := range g.data {
		g.rows[row[0]] = &row
	}
}

func (g *Grid) Find(tag string) (row []interface, found bool) {
	row, found = rows[tag]
	return row, found
}

type Cache struct {
	tables map[string](*Table)
	grids map[string](*Grid)
	srv *sheets.Service
	mtx sync.Mutex
}

func NewCache(srv *sheets.Service) (c Cache) {
	c.tables = make(map[string]*Table)
	c.grids = make(map[string]*Grid)
	c.srv = srv
	return c
}

func (c Cache) Register(t *Table) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	_, registered := c.tables[t.sheet]
	if !registered {
		t.tables[t] = NewResource(ssid, rng)
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
