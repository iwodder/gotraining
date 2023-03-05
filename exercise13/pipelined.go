//go:build pipeline

package main

import (
	"context"
	"fmt"
	"gotraining/exercise13/hn"
	"sort"
	"sync"
)

type work struct {
	idx    int
	itemId int
	item
}

func getTopStories(numStories int) ([]item, error) {
	fmt.Println("getTopStories() -> PIPELINED")
	var client hn.Client
	ids, err := client.TopItems()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return collect(cancel, numStories, retrieveStories(ctx, sendWork(ctx, ids))), nil
}

func sendWork(ctx context.Context, ids []int) <-chan work {
	workStream := make(chan work)
	go func() {
		defer close(workStream)
		for idx, v := range ids {
			select {
			case workStream <- work{idx: idx, itemId: v}:
			case <-ctx.Done():
				return
			}
		}
	}()
	return workStream
}

func retrieveStories(ctx context.Context, input <-chan work) <-chan work {
	output := make(chan work)
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			defer wg.Done()
			var client hn.Client
			for work := range input {
				i, err := client.GetItem(work.itemId)
				if err != nil {
					continue
				}
				work.item = parseHNItem(i)
				select {
				case output <- work:
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(output)
	}()
	return output
}

func collect(c context.CancelFunc, limit int, input <-chan work) []item {
	var results []work
	for w := range input {
		if isStoryLink(w.item) {
			results = append(results, w)
			if len(results) == limit {
				c()
			}
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].idx < results[j].idx
	})
	var ret []item
	for i := range results {
		ret = append(ret, results[i].item)
	}
	return ret[:limit]
}
