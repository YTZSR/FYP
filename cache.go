package main

import (
	"fmt"
	"time"
)

const CACHE_SIZE = 10

type proCache struct {
	items        []proCacheItem
	periodCalled []int
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
		items: make([]proCacheItem, CACHE_SIZE),
	}
	for i := 0; i < CACHE_SIZE; i++ {
		item := proCacheItem{
			ulid:         "",
			minTime:      time.Now(),
			maxTime:      time.Now(),
			totalCalled:  0,
			periodCalled: 0,
			size:         0,
		}
		c.items[i] = item
	}
	return c
}

func (c *proCache) AddCacheItem(ulid string, minTime time.Time, maxTime time.Time, size int) error {
	i := 0
	for {
		if i < CACHE_SIZE && c.items[i].ulid != "" {
			i++
		} else {
			break
		}
	}
	if i == CACHE_SIZE {
		minperiodCalled := int(^uint(0) >> 1)
		minminTime := time.Now()
		index := 0
		for j := 0; j < CACHE_SIZE; j++ {
			if c.items[j].periodCalled < minperiodCalled {
				minperiodCalled = c.items[j].periodCalled
				minminTime = c.items[j].minTime
				index = j
			} else if c.items[j].periodCalled == minperiodCalled {
				if minminTime.After(c.items[j].minTime) {
					minminTime = c.items[j].minTime
					index = j
				}
			}
		}
		i = index
		//TODO: Remove previous folder

	}
	c.items[i].ulid = ulid
	c.items[i].minTime = minTime
	c.items[i].maxTime = maxTime
	c.items[i].totalCalled = 0
	c.items[i].periodCalled = 0
	c.items[i].size = size
	//TODO: Input current folder

	return nil
}

func (c *proCache) ItemExistedByName(ulid string) (int, error) {
	index := 0
	for {
		if index < CACHE_SIZE && c.items[index].ulid != ulid {
			index++
		} else {
			break
		}
	}
	if index == CACHE_SIZE {
		return 0, fmt.Errorf("ulid is not existed")
	}
	return index, nil
}

func (c *proCache) ItemExistedByTime(t time.Time) (int, error) {
	for index, item := range c.items {
		if t.After(item.minTime) && t.Before(item.maxTime) {
			return index, nil
		}
	}
	return 0, fmt.Errorf("time is not in range")
}

// func (c *proCache) FilterCache() error {
// 	for index, item := range c.items {
// 		// item := c.items[ulid]
// 		if item.periodCalled < 0 { // new block
// 			item2 := item
// 			item2.periodCalled += 1
// 			c.items[ulid] = item2
// 		} else {
// 			if item.periodCalled < 5 {
// 				c.DeleteCacheItem(ulid)
// 			} else {
// 				item2 := item
// 				item2.periodCalled = 0
// 				c.items[ulid] = item2
// 			}
// 		}
// 	}
// 	return nil
// }

// func (c *proCache) CallItem(ulid string) error {

// 	item, ok := c.items[ulid]
// 	if ok {
// 		item.totalCalled += 1
// 		if item.periodCalled >= 0 {
// 			item.periodCalled += 1
// 		}
// 		c.items[ulid] = item
// 		return nil
// 	}
// 	return errors.New("ulid is in valid")

// }

func main() {
	var temp = NewCache()

	// fmt.Print(len(temp.items))
	for i := 0; i < 10; i++ {
		temp.AddCacheItem("ABCD", time.Now(), time.Now(), 1024)
	}
	temp.AddCacheItem("ABCD", time.Now(), time.Now(), 100)
	fmt.Print(temp.items[0].size)
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
