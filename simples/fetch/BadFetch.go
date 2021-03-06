package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Item struct {
	Title, Channel, GUID string
}

type Fetcher interface {
	Fetch() (items []Item, next time.Time, err error)
}

type Subscription interface {
	Updates() <-chan Item
	Close() error
}

type FakeFetcher struct {
	channel string
	items   []Item
}

func (f *FakeFetcher) Fetch() (items []Item, next time.Time, err error) {
	now := time.Now()
	next = now.Add(time.Duration(rand.Intn(5)) * 500 * time.Millisecond)
	item := Item{
		Title:   fmt.Sprintf("Item %d", len(f.items)),
		Channel: f.channel,
	}
	item.GUID = item.Channel + "/" + item.Title
	f.items = append(f.items, item)
	items = []Item{item}
	return
}

type NaiveSub struct {
	closed  bool
	err     error
	updates chan Item
	fetcher Fetcher
}

func (s *NaiveSub) Close() error {
	s.closed = true
	return s.err
}

func (s *NaiveSub) Updates() <-chan Item {
	return s.updates
}

func (s *NaiveSub) Loop() {
	for {
		if s.closed {
			close(s.updates)
			return
		}

		items, next, err := s.fetcher.Fetch()

		if err != nil {
			s.err = err
			time.Sleep(10 * time.Second)
			continue
		}

		for _, item := range items {
			s.updates <- item
		}

		now := time.Now()

		if next.After(now) {
			time.Sleep(next.Sub(now))
		}
	}
}

type NaiveMerge struct {
	subs    []Subscription
	updates chan Item
}

func (m *NaiveMerge) Close() (err error) {
	for _, sub := range m.subs {
		if e := sub.Close(); err == nil && e != nil {
			err = e
		}
	}
	close(m.updates)
	return
}

func (m *NaiveMerge) Updates() <-chan Item {
	return m.updates
}

func Merge(subs ...Subscription) Subscription {
	m := &NaiveMerge{
		subs:    subs,
		updates: make(chan Item),
	}

	for _, sub := range subs {
		go func(s Subscription) {
			for it := range s.Updates() {
				m.updates <- it
			}
		}(sub)
	}

	return m
}

func Subscribe(fetcher Fetcher) Subscription {
	s := &NaiveSub{
		fetcher: fetcher,
		updates: make(chan Item),
	}

	go s.Loop()
	return s
}

func main() {
	fetcher1 := &FakeFetcher{channel: "baidu.com"}
	fetcher2 := &FakeFetcher{channel: "google.com"}
	fetcher3 := &FakeFetcher{channel: "facebook.com"}
	fetcher4 := &FakeFetcher{channel: "twitter.com"}

	m := Merge(
		Subscribe(fetcher1),
		Subscribe(fetcher2),
		Subscribe(fetcher3),
		Subscribe(fetcher4))

	time.AfterFunc(5*time.Second, func() {
		fmt.Println("closed:", m.Close())
	})

	for item := range m.Updates() {
		fmt.Println(item.Channel, item.Title)
	}

}
