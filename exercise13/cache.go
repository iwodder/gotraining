package main

import (
	"fmt"
	"sync"
	"time"
)

type cache struct {
	mu     sync.Mutex
	limit  int
	active []item
	next   []item
}

func newCache(limit int) *cache {
	c := &cache{limit: limit}
	go func() {
		for fetchNext := time.NewTimer(10 * time.Second); ; fetchNext.Reset(10 * time.Second) {
			<-fetchNext.C
			fmt.Println("Updating next cache entry...")
			stories, err := getTopStories(limit)
			if err != nil {
				continue
			}
			c.mu.Lock()
			c.next = stories
			c.mu.Unlock()
			fmt.Println("Updated next cache entry")
		}
	}()
	go func() {
		for swap := time.NewTimer(15 * time.Second); ; swap.Reset(15 * time.Second) {
			<-swap.C
			fmt.Println("Swapping active to next")
			c.mu.Lock()
			c.active = c.next
			c.mu.Unlock()
			fmt.Println("Swapped active to next")
		}
	}()
	return c
}

func (c *cache) getStories() ([]item, error) {
	c.mu.Lock()
	stories := c.active
	c.mu.Unlock()
	if stories == nil {
		stories, err := getTopStories(c.limit)
		if err != nil {
			return nil, err
		}
		c.mu.Lock()
		if c.active == nil {
			c.active = stories
		}
		c.mu.Unlock()
	}
	return stories, nil
}
