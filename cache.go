package main

import (
	"errors"
	"time"
)

type proCache struct {
	items map[string]proCacheItem
}

type proCacheItem struct {
	ulid         string
	minTime      time.Time
	maxTime      time.Time
	totalCalled  int
	periodCalled int
	size         int
}

func NewCache() *proCache {
	c := &proCache{
		items: make(map[string]proCacheItem),
	}
	return c
}

func newCacheItem(minTime time.Time, maxTime time.Time, size int) (*proCacheItem, error) {
	c := &proCacheItem{
		// ulid:         ulid,
		minTime:      minTime,
		maxTime:      maxTime,
		totalCalled:  0,
		periodCalled: -5,
		size:         size,
	}
	return c, nil
}

func (c *proCache) ItemExistedByName(ulid string) bool {
	_, ok := c.items[ulid]
	if ok {
		return true
	}
	return false
}

func (c *proCache) ItemExistedByTime(t time.Time) string {
	for ulid, item := range c.items {
		if t.After(item.minTime) && t.Before(item.maxTime) {
			return ulid
		}
	}
	return ""
}

func (c *proCache) AddCacheItem(ulid string, minTime time.Time, maxTime time.Time, size int) error {

	_, ok := c.items[ulid]
	if ok {
		return errors.New("ulid is already existed")
	}

	newItem, error := newCacheItem(minTime, maxTime, size)
	if error == nil {
		c.items[ulid] = *newItem
		return nil
	}

	return errors.New("Failed to generate the cache item")

}

func (c *proCache) DeleteCacheItem(ulid string) error {
	_, ok := c.items[ulid]
	if ok {
		delete(c.items, ulid)
		return nil
	}
	return errors.New("ulid not existed")
}

func (c *proCache) FilterCache() error {
	for ulid, item := range c.items {
		// item := c.items[ulid]
		if item.periodCalled < 0 { // new block
			item2 := item
			item2.periodCalled += 1
			c.items[ulid] = item2
		} else {
			if item.periodCalled < 5 {
				c.DeleteCacheItem(ulid)
			} else {
				item2 := item
				item2.periodCalled = 0
				c.items[ulid] = item2
			}
		}
	}
	return nil
}

func (c *proCache) CallItem(ulid string) error {

	item, ok := c.items[ulid]
	if ok {
		item.totalCalled += 1
		if item.periodCalled >= 0 {
			item.periodCalled += 1
		}
		c.items[ulid] = item
		return nil
	}
	return errors.New("ulid is in valid")

}

func main() {
	// var temp = NewCache()

	// temp.AddCacheItem("ABCD", time.Now(), time.Now(), 1024)
	// // var len int = temp.items.Len()
	// // fmt.Print(len)
	// for i := 0; i < 6; i++ {
	// 	temp.FilterCache()
	// 	fmt.Println(temp.items["ABCD"].periodCalled)
	// }

	// error := temp.CallItem("ABCD")
	// if error != nil {
	// 	fmt.Println(error)
	// }
	// fmt.Println(temp.items["ABCD"].periodCalled)

	// // temp.deleteCacheItem("ABCD")

	// slice1 := make([]proCacheItem, 10)

}
