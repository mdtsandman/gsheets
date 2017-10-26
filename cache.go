package gsheets

import (
	"sync"
)

type Cache struct{
	data map[string][][]interface{}
	mtx sync.Mutex
}

func NewCache() (c Cache) {
	c.data = make(map[string][][]interface{})
	return c
}

func (c Cache) Data(resource string) (data [][]interface{}, ok bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	_, registered := c.data[resource]
	if (!registered) {
		return data, false
	}
	return c.data[resource], true
}

func (c Cache) Update(resource string, data [][]interface{}) (cache Cache) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.data[resource] = data
	return c
}
