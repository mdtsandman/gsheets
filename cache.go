package gsheets

import (
	"fmt"
	"sync"
	"errors"
)

type Cache struct {
	sheets map[string]*Sheet
	mtx sync.Mutex
}

func NewCache() (c Cache) {
	c.sheets = make(map[string]*Sheet)
	return c
}

func (c Cache) Register(s *Sheet) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.sheets[s.Name()] = s
	return
}

func (c Cache) Get(name string) (s *Sheet, err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	s, registered := c.sheets[name]
	if !registered {
		return s, errors.New(fmt.Sprintf("Cache: Attempt to access unregistered table [%s]", name))
	}
	if s.Stale() {
		err = s.Refresh()
		return s, err
	}
	return s, nil
}

func (c Cache) SetStaleFlag(name string) (err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	s, registered := c.sheets[name]
	if !registered {
		return errors.New(fmt.Sprintf("Cache: Attempt to set stale flag on unregistered resource [%s]", name))
	}
	s.SetStale(true)
	return nil
}

func (c Cache) Refresh(name string) (err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	s, registered := c.sheets[name]
	if !registered {
		return errors.New(fmt.Sprintf("Cache: Attempt to update unregistered resource [%s]", name))
	}
	return s.Refresh()
}

func (c Cache) RefreshAll() (error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for _, s := range c.sheets {
		if e := s.Refresh(); e != nil {
			return e
		}
	}
	return nil
}
