//go:build sequential

package main

func getTopStories(numStories int) ([]item, error) {
	var client hn.Client
	ids, err := client.TopItems()
	var stories []item
	for _, id := range ids {
		hnItem, err := client.GetItem(id)
		if err != nil {
			continue
		}
		item := parseHNItem(hnItem)
		if isStoryLink(item) {
			stories = append(stories, item)
			if len(stories) >= numStories {
				break
			}
		}
	}
	return stories, err
}
