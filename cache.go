package gsheets

import (
	"fmt"
	"sync"
	"errors"
)

type Cache struct {
	sheets map[string]interface{}
	mtx sync.Mutex
}

func NewCache() (c Cache) {
	c.sheets = make(map[string]interface{})
	return c
}

func (c Cache) Register(obj interface{}) (ok bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if s, ok := obj.(Sheet); ok {
		c.sheets[s.Name()] = obj
		return true
	}
	return false // can only register objects that implement the Sheet interface
}

func (c Cache) Get(name string) (obj interface{}, err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	obj, registered := c.sheets[name]
	if !registered {
		return obj, errors.New(fmt.Sprintf("Cache: Attempt to access unregistered table [%s]", name))
	}
	s, _ := obj.(Sheet)
	if s.Stale() {
		err = s.Refresh()
		return obj, err
	}
	return obj, nil
}

func (c Cache) SetStaleFlag(name string) (err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	obj, registered := c.sheets[name]
	if !registered {
		return errors.New(fmt.Sprintf("Cache: Attempt to set stale flag on unregistered resource [%s]", name))
	}
	s, _ := obj.(Sheet)
	s.SetStale(true)
	return nil
}

func (c Cache) Refresh(name string) (err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	obj, registered := c.sheets[name]
	if !registered {
		return errors.New(fmt.Sprintf("Cache: Attempt to update unregistered resource [%s]", name))
	}
	s, _ := obj.(Sheet)
	return s.Refresh()
}

func (c Cache) RefreshAll() (err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for _, obj := range c.sheets {
		s, _ := obj.(Sheet)
		if err := s.Refresh(); err != nil {
			return err
		}
	}
	return nil
}
